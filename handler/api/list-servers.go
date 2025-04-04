package api

import (
	"github.com/chatmcp/mcprouter/service/api"
	"github.com/chatmcp/mcprouter/service/mcpserver"
	"github.com/labstack/echo/v4"
)

type ListServersRequest struct {
	Page  int `json:"page,omitempty"`
	Limit int `json:"limit,omitempty"`
}

func ListServers(c echo.Context) error {
	ctx := api.GetAPIContext(c)

	req := &ListServersRequest{}

	if err := ctx.Valid(req); err != nil {
		return ctx.RespErr(err)
	}

	servers, err := mcpserver.GetHostedServers(req.Page, req.Limit)
	if err != nil {
		return ctx.RespErr(err)
	}

	return ctx.RespData(servers)
}
