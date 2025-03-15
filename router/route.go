package router

import (
	"github.com/chatmcp/mcprouter/api"
	"github.com/labstack/echo/v4"
)

// Route will create the routes for the http server
func Route(e *echo.Echo) {
	e.GET("/ping", api.Ping)
	e.GET("/sse/:key", api.SSE)
	e.POST("/messages", api.Messages)
}
