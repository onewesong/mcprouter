package proxy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/chatmcp/mcprouter/model"
	"github.com/chatmcp/mcprouter/service/jsonrpc"
	"github.com/chatmcp/mcprouter/service/mcpclient"
	"github.com/chatmcp/mcprouter/service/mcpserver"
	"github.com/chatmcp/mcprouter/service/proxy"
	"github.com/chatmcp/mcprouter/util"
	"github.com/labstack/echo/v4"
)

// MCP is a handler for the mcp endpoint
func MCP(c echo.Context) error {
	ctx := proxy.GetSSEContext(c)
	if ctx == nil {
		return c.String(http.StatusInternalServerError, "Failed to get SSE context")
	}

	req := c.Request()
	method := req.Method
	header := req.Header

	accept := header.Get("Accept")
	// accept: application/json, text/event-stream

	path := req.URL.Path

	sessionID := req.Header.Get("Mcp-Session-Id")

	log.Printf("method: %s, accept: %s, sessionID: %s, path: %s\n", method, accept, sessionID, path)

	if method != http.MethodPost {
		// todo: return event-stream response when method is GET
		// todo: delete session when method is DELETE
		return c.String(http.StatusMethodNotAllowed, "Method Not Allowed")
	}

	key := c.Param("key")
	if key == "" {
		return c.String(http.StatusBadRequest, "Key is required")
	}

	serverConfig := mcpserver.GetServerConfig(key)
	if serverConfig == nil {
		return c.String(http.StatusBadRequest, "Invalid server config")
	}

	request, err := ctx.GetJSONRPCRequest()
	if err != nil {
		return ctx.JSONRPCError(jsonrpc.ErrorParseError, nil)
	}

	if request.Result != nil || request.Error != nil {
		// notification
		return ctx.JSONRPCAcceptResponse(nil)
	}

	proxyInfo := &proxy.ProxyInfo{
		ServerKey:          key,
		SessionID:          sessionID,
		ServerUUID:         serverConfig.ServerUUID,
		ServerConfigName:   serverConfig.ServerName,
		ServerShareProcess: serverConfig.ShareProcess,
		ServerType:         serverConfig.ServerType,
		ServerURL:          serverConfig.ServerURL,
		ServerCommand:      serverConfig.Command,
		ServerCommandHash:  serverConfig.CommandHash,
		RequestID:          header.Get("X-Request-ID"),
		RequestFrom:        header.Get("X-Request-From"),
	}

	// log.Printf("request: %+v\n", request)

	if request.Method == "initialize" {
		paramsB, _ := json.Marshal(request.Params)
		params := &jsonrpc.InitializeParams{}
		if err := json.Unmarshal(paramsB, params); err != nil {
			return ctx.JSONRPCError(jsonrpc.ErrorParseError, nil)
		}

		// start new session
		sessionID = util.MD5(key)

		proxyInfo.ConnectionTime = time.Now()
		proxyInfo.ClientName = params.ClientInfo.Name
		proxyInfo.ClientVersion = params.ClientInfo.Version
		proxyInfo.ProtocolVersion = params.ProtocolVersion
		proxyInfo.SessionID = sessionID

		err := proxy.StoreProxyInfo(sessionID, proxyInfo)

		log.Printf("store proxy info with client info: %s, %v, %s\n", sessionID, err, proxyInfo.ClientName)

		ctx.Response().Header().Set("Mcp-Session-Id", sessionID)
	} else {
		// not initialize request, check session
		if sessionID == "" {
			return c.String(http.StatusBadRequest, "Invalid session ID")
		}

		_proxyInfo, err := proxy.GetProxyInfo(sessionID)

		log.Printf("get proxy info from cache: %s, %v, %+v\n", sessionID, err, _proxyInfo)

		if _proxyInfo != nil && _proxyInfo.SessionID == sessionID {
			proxyInfo = _proxyInfo
		}
	}

	proxyInfo.JSONRPCVersion = request.JSONRPC
	proxyInfo.RequestMethod = request.Method
	proxyInfo.RequestTime = time.Now()
	proxyInfo.RequestParams = request.Params

	if request.ID != nil {
		proxyInfo.RequestID = request.ID
	}

	client := ctx.GetClient(key)

	if client == nil {
		_client, err := mcpclient.NewClient(serverConfig)
		if err != nil {
			fmt.Printf("connect to mcp server failed: %v\n", err)
			return ctx.JSONRPCError(jsonrpc.ErrorProxyError, request.ID)
		}

		if err := _client.Error(); err != nil {
			fmt.Printf("mcp server run failed: %v\n", err)
			return ctx.JSONRPCError(jsonrpc.ErrorProxyError, request.ID)
		}

		ctx.StoreClient(key, _client)

		client = _client

		client.OnNotification(func(message []byte) {
			fmt.Printf("received notification: %s\n", message)
		})
	}

	if client == nil {
		return ctx.JSONRPCError(jsonrpc.ErrorProxyError, request.ID)
	}

	response, err := client.ForwardMessage(request)
	if err != nil {
		fmt.Printf("forward message failed: %v\n", err)
		client.Close()
		ctx.DeleteClient(key)
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

			proxyInfo.ResponseResult = response

			err = proxy.StoreProxyInfo(sessionID, proxyInfo)

			log.Printf("store proxy info with server info: %s, %v, %s\n", sessionID, err, proxyInfo.ServerName)
		}
	}

	proxyInfo.ResponseResult = response

	proxyInfo.ResponseTime = time.Now()
	costTime := proxyInfo.ResponseTime.Sub(proxyInfo.RequestTime)
	proxyInfo.CostTime = costTime.Milliseconds()

	proxyInfoB, _ := json.Marshal(proxyInfo)

	if proxyInfo.RequestMethod == "tools/call" {
		if err := model.CreateServerLog(proxyInfo.ToServerLog()); err != nil {
			log.Printf("save server log failed: %v\n", err)
		} else {
			log.Printf("save server log ok: %s\n", proxyInfo.RequestID)
		}
	}
	log.Printf("proxyInfo: %s\n", string(proxyInfoB))

	// notification
	if response == nil {
		return ctx.JSONRPCAcceptResponse(response)
	}

	return ctx.JSONRPCResponse(response)
}
