package router

import (
	"github.com/chatmcp/mcprouter/handler/api"
	"github.com/labstack/echo/v4"
)

// APIRoute will create the routes for the http server
func APIRoute(e *echo.Echo) {
	apiv1 := e.Group("/v1")

	apiv1.POST("/list-servers", api.ListServers)
	apiv1.POST("/list-tools", api.ListTools)
	apiv1.POST("/call-tool", api.CallTool)
}
