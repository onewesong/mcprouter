package proxy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/chatmcp/mcprouter/service/jsonrpc"
	"github.com/chatmcp/mcprouter/service/mcpclient"
	"github.com/chatmcp/mcprouter/service/proxy"
	"github.com/labstack/echo/v4"
)

// Messages is a handler for the messages endpoint
func Messages(c echo.Context) error {
	ctx := proxy.GetSSEContext(c)
	if ctx == nil {
		return c.String(http.StatusInternalServerError, "Failed to get SSE context")
	}

	sessionID := ctx.QueryParam("sessionid")
	if sessionID == "" {
		return ctx.JSONRPCError(jsonrpc.ErrorInvalidParams, nil)
	}

	session := ctx.GetSession(sessionID)
	if session == nil {
		return ctx.JSONRPCError(jsonrpc.ErrorInvalidParams, nil)
	}

	request, err := ctx.GetJSONRPCRequest()
	if err != nil {
		return ctx.JSONRPCError(jsonrpc.ErrorParseError, nil)
	}

	proxyInfo := session.ProxyInfo()
	sseKey := session.Key()

	proxyInfo.JSONRPCVersion = request.JSONRPC
	proxyInfo.RequestMethod = request.Method
	proxyInfo.RequestTime = time.Now()
	proxyInfo.RequestParams = request.Params

	if request.ID != nil {
		proxyInfo.RequestID = request.ID
	}

	if request.Method == "initialize" {
		paramsB, _ := json.Marshal(request.Params)
		params := &jsonrpc.InitializeParams{}
		if err := json.Unmarshal(paramsB, params); err != nil {
			return ctx.JSONRPCError(jsonrpc.ErrorParseError, nil)
		}

		proxyInfo.ClientName = params.ClientInfo.Name
		proxyInfo.ClientVersion = params.ClientInfo.Version
		proxyInfo.ProtocolVersion = params.ProtocolVersion

		session.SetProxyInfo(proxyInfo)
		ctx.StoreSession(sessionID, session)
	}

	client := ctx.GetClient(sseKey)

	if client == nil {
		command := proxyInfo.ServerCommand
		_client, err := mcpclient.NewStdioClient(command)
		if err != nil {
			fmt.Printf("connect to mcp server failed: %v\n", err)
			return ctx.JSONRPCError(jsonrpc.ErrorProxyError, request.ID)
		}

		if err := _client.Error(); err != nil {
			fmt.Printf("mcp server run failed: %v\n", err)
			return ctx.JSONRPCError(jsonrpc.ErrorProxyError, request.ID)
		}

		ctx.StoreClient(sseKey, _client)

		client = _client

		client.OnNotification(func(message []byte) {
			fmt.Printf("received notification: %s\n", message)
			session.SendMessage(string(message))
		})
	}

	if client == nil {
		return ctx.JSONRPCError(jsonrpc.ErrorProxyError, request.ID)
	}

	response, err := client.ForwardMessage(request)
	if err != nil {
		fmt.Printf("forward message failed: %v\n", err)
		session.Close()
		ctx.DeleteClient(sseKey)
		return ctx.JSONRPCError(jsonrpc.ErrorProxyError, request.ID)
	}

	if response != nil {
		if request.Method == "initialize" && response.Result != nil {
			resultB, _ := json.Marshal(response.Result)
			result := &jsonrpc.InitializeResult{}
			if err := json.Unmarshal(resultB, result); err != nil {
				fmt.Printf("unmarshal initialize result failed: %v\n", err)
				return ctx.JSONRPCError(jsonrpc.ErrorParseError, request.ID)
			}

			proxyInfo.ServerName = result.ServerInfo.Name
			proxyInfo.ServerVersion = result.ServerInfo.Version

			session.SetProxyInfo(proxyInfo)
			ctx.StoreSession(sessionID, session)
		}

		// not notification message, send sse message
		session.SendMessage(response.String())
	}

	proxyInfo.ResponseTime = time.Now()
	proxyInfo.ResponseDuration = time.Since(proxyInfo.RequestTime)
	proxyInfoB, _ := json.Marshal(proxyInfo)

	fmt.Printf("proxyInfo: %s, cost: %f\n", string(proxyInfoB), proxyInfo.ResponseDuration.Seconds())

	return ctx.JSONRPCResponse(response)
}
