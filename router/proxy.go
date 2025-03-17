package router

import (
	"github.com/chatmcp/mcprouter/handler/proxy"
	"github.com/labstack/echo/v4"
)

// ProxyRoute will create the routes for the http server
func ProxyRoute(e *echo.Echo) {
	e.GET("/sse/:key", proxy.SSE)
	e.POST("/messages", proxy.Messages)
}
