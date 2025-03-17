package api

import (
	"github.com/chatmcp/mcprouter/service/api"
	"github.com/labstack/echo/v4"
)

// ListTools is a handler for the list tools endpoint
func ListTools(c echo.Context) error {
	ctx := api.GetAPIContext(c)

	client, err := ctx.Connect()
	if err != nil {
		return ctx.RespErr(err)
	}
	defer client.Close()

	tools, err := client.ListTools()
	if err != nil {
		return ctx.RespErr(err)
	}

	return ctx.RespData(tools)
}
