package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/chatmcp/mcprouter/service/jsonrpc"
	"github.com/chatmcp/mcprouter/service/mcpclient"
	"github.com/chatmcp/mcprouter/service/mcpserver"
	"github.com/chatmcp/mcprouter/service/proxy"
	"github.com/labstack/echo/v4"
)

type APIErrorCode int

const (
	APIErrorFail    APIErrorCode = -1
	APIErrorSuccess APIErrorCode = 0
	APIErrorNoAuth  APIErrorCode = -2
)

type APIResponse struct {
	Code    APIErrorCode `json:"code"`
	Message string       `json:"message"`
	Data    interface{}  `json:"data,omitempty"`
}

// APIContext is the context for the API request
type APIContext struct {
	echo.Context
	serverConfig *mcpserver.ServerConfig
	clientInfo   *jsonrpc.ClientInfo
	proxyInfo    *proxy.ProxyInfo
}

func createAPIMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := &APIContext{
				Context: c,
			}

			header := c.Request().Header

			authorization := header.Get("Authorization")
			if authorization == "" {
				return ctx.RespNoAuthMsg("no authorization header")
			}

			apikey := strings.TrimSpace(strings.ReplaceAll(authorization, "Bearer", ""))
			if apikey == "" {
				return ctx.RespNoAuthMsg("no authorization key")
			}

			serverConfig := mcpserver.GetServerConfig(apikey)
			if serverConfig == nil || serverConfig.Command == "" {
				return ctx.RespNoAuthMsg("invalid authorization key")
			}

			ctx.serverConfig = serverConfig

			clientInfo := header.Get("X-Client-Info")
			if clientInfo == "" {
				return ctx.RespNoAuthMsg("no client info")
			}

			ctx.clientInfo = &jsonrpc.ClientInfo{}
			if err := json.Unmarshal([]byte(clientInfo), ctx.clientInfo); err != nil {
				return ctx.RespNoAuthMsg("invalid client info")
			}

			return next(ctx)
		}
	}
}

// GetAPIContext returns the APIContext from the echo.Context
func GetAPIContext(c echo.Context) *APIContext {
	if ctx, ok := c.(*APIContext); ok {
		return ctx
	}

	return nil
}

// Valid: valid request params
func (c *APIContext) Valid(req interface{}) error {
	if err := c.Bind(req); err != nil {
		if v, ok := err.(*echo.HTTPError); ok {
			return fmt.Errorf("%s", v.Message)
		}

		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	return nil
}

// ProxyInfo returns the proxy info
func (c *APIContext) ProxyInfo() *proxy.ProxyInfo {
	return c.proxyInfo
}

func (c *APIContext) SetProxyInfo(proxyInfo *proxy.ProxyInfo) {
	c.proxyInfo = proxyInfo
}

// ClientInfo returns the client info
func (c *APIContext) ClientInfo() *jsonrpc.ClientInfo {
	return c.clientInfo
}

// ServerConfig returns the server config
func (c *APIContext) ServerConfig() *mcpserver.ServerConfig {
	return c.serverConfig
}

// ServerCommand returns the server command
func (c *APIContext) ServerCommand() string {
	return c.ServerConfig().Command
}

// Connect connects to the mcp server
func (c *APIContext) Connect() (*mcpclient.StdioClient, error) {
	serverConfig := c.ServerConfig()
	clientInfo := c.ClientInfo()

	header := c.Request().Header

	proxyInfo := &proxy.ProxyInfo{
		ClientName:         clientInfo.Name,
		ClientVersion:      clientInfo.Version,
		ServerUUID:         serverConfig.ServerUUID,
		ServerKey:          serverConfig.ServerKey,
		ServerConfigName:   serverConfig.ServerName,
		ServerShareProcess: serverConfig.ShareProcess,
		ServerCommand:      serverConfig.Command,
		ServerCommandHash:  serverConfig.CommandHash,
		ConnectionTime:     time.Now(),
		RequestTime:        time.Now(),
		RequestID:          header.Get("X-Request-ID"),
		RequestFrom:        header.Get("X-Request-From"),
	}

	command := c.ServerCommand()
	if command == "" {
		return nil, fmt.Errorf("invalid command")
	}

	client, err := mcpclient.NewStdioClient(command)
	if err != nil {
		return nil, fmt.Errorf("connect to mcp server failed")
	}

	// initialize get server info
	result, err := client.Initialize(&jsonrpc.InitializeParams{
		ProtocolVersion: jsonrpc.JSONRPC_VERSION,
		Capabilities:    jsonrpc.ClientCapabilities{},
		ClientInfo: jsonrpc.ClientInfo{
			Name:    proxy.ProxyClientName,
			Version: proxy.ProxyClientVersion,
		},
	})

	if err != nil {
		client.Close()
		return nil, fmt.Errorf("connection initialize failed")
	}

	proxyInfo.ServerName = result.ServerInfo.Name
	proxyInfo.ServerVersion = result.ServerInfo.Version
	proxyInfo.JSONRPCVersion = jsonrpc.JSONRPC_VERSION
	proxyInfo.ProtocolVersion = result.ProtocolVersion

	c.SetProxyInfo(proxyInfo)

	if err := client.NotificationsInitialized(); err != nil {
		client.Close()
		return nil, fmt.Errorf("connection notifications initialized failed")
	}

	return client, nil
}

func (c *APIContext) RespErr(err error) error {
	return c.JSON(http.StatusOK, APIResponse{
		Code:    APIErrorFail,
		Message: err.Error(),
	})
}

func (c *APIContext) RespErrMsg(message string) error {
	return c.JSON(http.StatusOK, APIResponse{
		Code:    APIErrorFail,
		Message: message,
	})
}

func (c *APIContext) RespOK() error {
	return c.JSON(http.StatusOK, APIResponse{
		Code:    APIErrorSuccess,
		Message: "ok",
	})
}

func (c *APIContext) RespOKMsg(message string) error {
	return c.JSON(http.StatusOK, APIResponse{
		Code:    APIErrorSuccess,
		Message: message,
	})
}

func (c *APIContext) RespNoAuth() error {
	return c.JSON(http.StatusOK, APIResponse{
		Code:    APIErrorNoAuth,
		Message: "no auth",
	})
}

func (c *APIContext) RespNoAuthMsg(message string) error {
	return c.JSON(http.StatusOK, APIResponse{
		Code:    APIErrorNoAuth,
		Message: message,
	})
}

func (c *APIContext) RespData(data interface{}) error {
	return c.JSON(http.StatusOK, APIResponse{
		Code:    APIErrorSuccess,
		Message: "ok",
		Data:    data,
	})
}

func (c *APIContext) RespJSON(code int, message string, data interface{}) error {
	return c.JSON(code, APIResponse{
		Code:    APIErrorCode(code),
		Message: message,
		Data:    data,
	})
}
