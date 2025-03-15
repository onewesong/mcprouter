package api

import (
	"fmt"
	"net/http"

	"github.com/chatmcp/mcprouter/service/sse"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// SSE is a handler for the sse endpoint
func SSE(c echo.Context) error {
	ctx := sse.GetSSEContext(c)
	if ctx == nil {
		return c.String(http.StatusInternalServerError, "Failed to get SSE context")
	}

	key := c.Param("key")
	if key == "" {
		return c.String(http.StatusBadRequest, "Key is required")
	}

	mcpServerCommand := viper.GetString(fmt.Sprintf("mcp_server_commands.%s", key))
	if mcpServerCommand == "" {
		return c.String(http.StatusBadRequest, "MCP server not found")
	}

	fmt.Printf("mcp server command: %s\n", mcpServerCommand)

	writer, err := sse.NewSSEWriter(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// store session
	sessionID := uuid.New().String()
	session := sse.NewSSESession(writer, mcpServerCommand)
	ctx.StoreSession(sessionID, session)
	defer ctx.DeleteSession(sessionID)

	// response to client with endpoint url
	messagesUrl := fmt.Sprintf("/messages?sessionid=%s", sessionID)
	writer.SendEventData("endpoint", messagesUrl)

	// listen to messages
	for {
		select {
		case message := <-session.Messages():
			fmt.Printf("sse send message: %s\n", message)
			writer.SendMessage(message)
		}
	}
}
