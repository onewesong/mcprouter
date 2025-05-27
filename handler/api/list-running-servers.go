package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/chatmcp/mcprouter/service/api"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// RunningServerInfo represents information about a running MCP server
type RunningServerInfo struct {
	ServerKey string `json:"server_key"`
	Status    string `json:"status"`
}

// ProxyListClientsResponse defines the response from proxy server
type ProxyListClientsResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Clients []RunningServerInfo `json:"clients,omitempty"`
}

// ListRunningServers is a handler for listing currently running MCP servers
func ListRunningServers(c echo.Context) error {
	ctx := api.GetAPIContext(c)

	// 获取proxy服务器地址
	proxyPort := viper.GetInt("proxy_server.port")
	if proxyPort == 0 {
		proxyPort = 8025
	}
	proxyURL := fmt.Sprintf("http://127.0.0.1:%d/admin/list-clients", proxyPort)

	// 发送请求到proxy服务器
	resp, err := http.Get(proxyURL)
	if err != nil {
		return ctx.RespErrMsg(fmt.Sprintf("Failed to connect to proxy server: %v", err))
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ctx.RespErrMsg(fmt.Sprintf("Failed to read proxy response: %v", err))
	}

	var proxyResp ProxyListClientsResponse
	if err := json.Unmarshal(body, &proxyResp); err != nil {
		return ctx.RespErrMsg(fmt.Sprintf("Failed to parse proxy response: %v", err))
	}

	// 检查proxy服务器的响应
	if !proxyResp.Success {
		return ctx.RespErrMsg(proxyResp.Message)
	}

	return ctx.RespData(proxyResp.Clients)
}
