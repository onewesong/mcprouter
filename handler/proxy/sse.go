package proxy

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chatmcp/mcprouter/service/mcpserver"
	"github.com/chatmcp/mcprouter/service/proxy"
	"github.com/chatmcp/mcprouter/util"
	"github.com/labstack/echo/v4"
)

// SSE is a handler for the sse endpoint
func SSE(c echo.Context) error {
	ctx := proxy.GetSSEContext(c)
	if ctx == nil {
		return c.String(http.StatusInternalServerError, "Failed to get SSE context")
	}

	return c.String(http.StatusNotFound, "sse connection not supported now, please use streamable http or stdio connection")

	req := c.Request()

	key := c.Param("key")
	if key == "" {
		return c.String(http.StatusBadRequest, "Key is required")
	}

	serverConfig := mcpserver.GetServerConfig(key)
	if serverConfig == nil {
		return c.String(http.StatusBadRequest, "Invalid server config")
	}

	writer, err := proxy.NewSSEWriter(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// sessionID := uuid.New().String()
	sessionID := util.MD5(key)

	proxyInfo := &proxy.ProxyInfo{
		ServerKey:          key,
		ConnectionTime:     time.Now(),
		SessionID:          sessionID,
		ServerUUID:         serverConfig.ServerUUID,
		ServerConfigName:   serverConfig.ServerName,
		ServerShareProcess: serverConfig.ShareProcess,
		ServerType:         serverConfig.ServerType,
		ServerURL:          serverConfig.ServerURL,
		ServerCommand:      serverConfig.Command,
		ServerCommandHash:  serverConfig.CommandHash,
	}

	// store session
	session := proxy.NewSSESession(writer, serverConfig, proxyInfo)
	ctx.StoreSession(sessionID, session)
	defer ctx.DeleteSession(sessionID)

	// Setup heartbeat ticker
	// heartbeatInterval := 30 * time.Second // adjust interval as needed
	// heartbeatTicker := time.NewTicker(heartbeatInterval)
	// defer heartbeatTicker.Stop()

	// Setup idle timeout
	// idleTimeout := 5 * time.Minute // adjust timeout as needed
	// idleTimer := time.NewTimer(idleTimeout)
	// defer idleTimer.Stop()

	// // Reset idle timer when activity occurs
	// resetIdleTimer := func() {
	// 	if !idleTimer.Stop() {
	// 		<-idleTimer.C
	// 	}
	// 	idleTimer.Reset(idleTimeout)
	// }

	go func() {
		for {
			select {
			case <-session.Done():
				return
			case <-req.Context().Done():
				return
				// case <-heartbeatTicker.C:
				// 	// Send heartbeat comment
				// 	// if err := writer.SendHeartbeat(); err != nil {
				// 	// 	session.Close()
				// 	// 	return
				// 	// }
				// case <-idleTimer.C:
				// 	// Close connection due to inactivity
				// 	session.Close()
				// 	return
			}
		}
	}()

	// response to client with endpoint url
	messagesUrl := fmt.Sprintf("/messages?sessionid=%s", sessionID)
	writer.SendEventData("endpoint", messagesUrl)

	// listen to messages
	for {
		select {
		case message := <-session.Messages():
			// Reset idle timer on message activity
			// resetIdleTimer()

			if err := writer.SendMessage(message); err != nil {
				fmt.Printf("sse failed to send message to session %s: %v\n", sessionID, err)
				session.Close() // Close session on send error
				return nil      // Exit the handler
			}
		case <-session.Done():
			fmt.Printf("session %s closed \n", sessionID)
			return nil
		case <-req.Context().Done():
			fmt.Println("sse request done")
			session.Close()
			return nil
		}
	}
}
