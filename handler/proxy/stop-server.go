package proxy

import (
	"fmt"
	"net/http"

	"github.com/chatmcp/mcprouter/service/proxy"
	"github.com/labstack/echo/v4"
)

// StopServerRequest defines the request structure for stopping MCP server
type StopServerRequest struct {
	ServerKey string `json:"server_key"`
	Force     bool   `json:"force,omitempty"`
}

// StopServerResponse defines the response structure
type StopServerResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// StopServer is a handler for stopping a specific MCP server managed by proxy
func StopServer(c echo.Context) error {
	ctx := proxy.GetSSEContext(c)
	if ctx == nil {
		return c.JSON(http.StatusInternalServerError, StopServerResponse{
			Success: false,
			Message: "Failed to get proxy context",
		})
	}

	var req StopServerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, StopServerResponse{
			Success: false,
			Message: fmt.Sprintf("Invalid request: %v", err),
		})
	}

	if req.ServerKey == "" {
		return c.JSON(http.StatusBadRequest, StopServerResponse{
			Success: false,
			Message: "server_key is required",
		})
	}

	// 获取并停止指定的客户端连接
	client := ctx.GetClient(req.ServerKey)
	if client == nil {
		return c.JSON(http.StatusNotFound, StopServerResponse{
			Success: false,
			Message: fmt.Sprintf("MCP server '%s' not found or not running", req.ServerKey),
		})
	}

	// 删除客户端连接，这会自动关闭MCP server进程
	ctx.DeleteClient(req.ServerKey)

	return c.JSON(http.StatusOK, StopServerResponse{
		Success: true,
		Message: fmt.Sprintf("MCP server '%s' stopped successfully", req.ServerKey),
	})
}
