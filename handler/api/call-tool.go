package api

import (
	"github.com/chatmcp/mcprouter/service/api"
	"github.com/chatmcp/mcprouter/service/jsonrpc"
	"github.com/labstack/echo/v4"
)

type CallToolRequest struct {
	Name      string                 `json:"name" validate:"required"`
	Arguments map[string]interface{} `json:"arguments" validate:"required"`
}

func CallTool(c echo.Context) error {
	ctx := api.GetAPIContext(c)

	req := &CallToolRequest{}

	if err := ctx.Valid(req); err != nil {
		return ctx.RespErr(err)
	}

	client, err := ctx.Connect()
	if err != nil {
		return ctx.RespErr(err)
	}
	defer client.Close()

	callToolResult, err := client.CallTool(&jsonrpc.CallToolParams{
		Name:      req.Name,
		Arguments: req.Arguments,
	})
	if err != nil {
		return ctx.RespErr(err)
	}

	return ctx.RespData(callToolResult)
}
