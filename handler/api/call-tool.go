package api

import (
	"encoding/json"
	"log"
	"time"

	"github.com/chatmcp/mcprouter/model"
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

	proxyInfo := ctx.ProxyInfo()
	proxyInfo.RequestMethod = jsonrpc.MethodCallTool

	requestParams := &jsonrpc.CallToolParams{
		Name:      req.Name,
		Arguments: req.Arguments,
	}

	proxyInfo.RequestParams = requestParams

	callToolResult, err := client.CallTool(requestParams)
	if err != nil {
		return ctx.RespErr(err)
	}

	proxyInfo.ResponseResult = callToolResult

	proxyInfo.ResponseTime = time.Now()
	proxyInfo.CostTime = proxyInfo.ResponseTime.Sub(proxyInfo.RequestTime).Milliseconds()

	proxyInfoB, _ := json.Marshal(proxyInfo)

	if proxyInfo.RequestMethod == "tools/call" {
		if err := model.CreateServerLog(proxyInfo.ToServerLog()); err != nil {
			log.Printf("save server log failed: %v\n", err)
		} else {
			log.Printf("save server log ok: %s\n", proxyInfo.RequestID)
		}
	}
	log.Printf("proxyInfo: %s\n", string(proxyInfoB))

	return ctx.RespData(callToolResult)
}
