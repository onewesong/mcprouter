package router

import (
	"github.com/chatmcp/mcprouter/handler/proxy"
	"github.com/labstack/echo/v4"
)

// ProxyRoute will create the routes for the http server
func ProxyRoute(e *echo.Echo) {
	// sse proxy
	e.GET("/sse/:key", proxy.SSE)
	e.POST("/messages", proxy.Messages)
	// streamable http proxy
	e.Any("/mcp/:key", proxy.MCP)
	// management endpoints
	e.POST("/admin/stop-server", proxy.StopServer)
	e.GET("/admin/list-clients", proxy.ListClients)
}
