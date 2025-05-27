package proxy

import (
	"net/http"

	"github.com/chatmcp/mcprouter/service/proxy"
	"github.com/labstack/echo/v4"
)

// ClientInfo represents information about a running MCP client
type ClientInfo struct {
	ServerKey string `json:"server_key"`
	Status    string `json:"status"`
}

// ListClientsResponse defines the response structure
type ListClientsResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Clients []ClientInfo `json:"clients,omitempty"`
}

// ListClients is a handler for listing currently running MCP clients
func ListClients(c echo.Context) error {
	ctx := proxy.GetSSEContext(c)
	if ctx == nil {
		return c.JSON(http.StatusInternalServerError, ListClientsResponse{
			Success: false,
			Message: "Failed to get proxy context",
		})
	}

	var clients []ClientInfo

	// 遍历所有存储的客户端
	ctx.RangeClients(func(key, value interface{}) bool {
		if serverKey, ok := key.(string); ok {
			clients = append(clients, ClientInfo{
				ServerKey: serverKey,
				Status:    "running",
			})
		}
		return true
	})

	return c.JSON(http.StatusOK, ListClientsResponse{
		Success: true,
		Message: "success",
		Clients: clients,
	})
}
