package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/chatmcp/mcprouter/service/api"
	"github.com/chatmcp/mcprouter/service/jsonrpc"
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

	proxyInfo := ctx.ProxyInfo()
	proxyInfo.RequestMethod = jsonrpc.MethodListTools

	tools, err := client.ListTools()
	if err != nil {
		return ctx.RespErr(err)
	}

	proxyInfo.ResponseResult = tools

	proxyInfo.ResponseTime = time.Now()
	proxyInfo.CostTime = proxyInfo.ResponseTime.Sub(proxyInfo.RequestTime).Milliseconds()

	proxyInfoB, _ := json.Marshal(proxyInfo)
	fmt.Printf("proxyInfo: %s\n", string(proxyInfoB))

	return ctx.RespData(tools)
}
