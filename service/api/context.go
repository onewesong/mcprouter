package api

import (
	"fmt"
	"net/http"
	"strings"

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
}

func createAPIMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := &APIContext{
				Context: c,
			}

			authorization := c.Request().Header.Get("Authorization")
			if authorization == "" {
				return ctx.RespNoAuthMsg("no authorization header")
			}

			apikey := strings.TrimPrefix(authorization, "Bearer ")
			if apikey == "" {
				return ctx.RespNoAuthMsg("no apikey")
			}

			serverConfig := mcpserver.GetServerConfig(apikey)
			if serverConfig == nil {
				return ctx.RespNoAuthMsg("invalid apikey")
			}

			ctx.serverConfig = serverConfig

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
	command := c.ServerCommand()
	if command == "" {
		return nil, fmt.Errorf("invalid command")
	}

	client, err := mcpclient.NewStdioClient(command)
	if err != nil {
		return nil, fmt.Errorf("connect to mcp server failed")
	}

	if _, err := client.Initialize(&jsonrpc.InitializeParams{
		ProtocolVersion: jsonrpc.JSONRPC_VERSION,
		Capabilities:    jsonrpc.ClientCapabilities{},
		ClientInfo: jsonrpc.ClientInfo{
			Name:    proxy.ProxyClientName,
			Version: proxy.ProxyClientVersion,
		},
	}); err != nil {
		client.Close()
		return nil, fmt.Errorf("connection initialize failed")
	}

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
