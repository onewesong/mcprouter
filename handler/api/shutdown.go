package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/chatmcp/mcprouter/service/api"
	"github.com/chatmcp/mcprouter/service/mcpserver"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// StopServerRequest defines the request structure for stop-server endpoint
type StopServerRequest struct {
	ServerKey string `json:"server_key" validate:"required"` // MCP server的key
	Force     bool   `json:"force,omitempty"`                // 是否强制停止
}

// ProxyStopServerResponse defines the response from proxy server
type ProxyStopServerResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Shutdown is a handler for stopping a specific MCP server
func Shutdown(c echo.Context) error {
	ctx := api.GetAPIContext(c)

	req := &StopServerRequest{}
	if err := ctx.Valid(req); err != nil {
		return ctx.RespErr(err)
	}

	// 检查server key是否有效
	serverConfig := mcpserver.GetServerConfig(req.ServerKey)
	if serverConfig == nil {
		return ctx.RespErrMsg(fmt.Sprintf("Invalid server key: %s", req.ServerKey))
	}

	// 获取proxy服务器地址
	proxyPort := viper.GetInt("proxy_server.port")
	if proxyPort == 0 {
		proxyPort = 8025
	}
	proxyURL := fmt.Sprintf("http://127.0.0.1:%d/admin/stop-server", proxyPort)

	// 构造请求
	requestData := map[string]interface{}{
		"server_key": req.ServerKey,
		"force":      req.Force,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return ctx.RespErrMsg(fmt.Sprintf("Failed to marshal request: %v", err))
	}

	// 发送请求到proxy服务器
	resp, err := http.Post(proxyURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return ctx.RespErrMsg(fmt.Sprintf("Failed to connect to proxy server: %v", err))
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ctx.RespErrMsg(fmt.Sprintf("Failed to read proxy response: %v", err))
	}

	var proxyResp ProxyStopServerResponse
	if err := json.Unmarshal(body, &proxyResp); err != nil {
		return ctx.RespErrMsg(fmt.Sprintf("Failed to parse proxy response: %v", err))
	}

	// 检查proxy服务器的响应
	if !proxyResp.Success {
		return ctx.RespErrMsg(proxyResp.Message)
	}

	return ctx.RespOKMsg(proxyResp.Message)
}
