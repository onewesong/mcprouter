package mcpclient

import (
	"encoding/json"

	"github.com/chatmcp/mcprouter/service/jsonrpc"
)

// Initialize initializes the client.
func (c *StdioClient) Initialize(params *jsonrpc.InitializeParams) (*jsonrpc.InitializeResult, error) {
	request := jsonrpc.NewRequest(jsonrpc.MethodInitialize, params, 0)

	response := c.ForwardRequest(request)

	resultB, err := json.Marshal(response.Result)
	if err != nil {
		return nil, err
	}

	result := &jsonrpc.InitializeResult{}
	if err := json.Unmarshal(resultB, result); err != nil {
		return nil, err
	}

	return result, nil
}

// NotificationsInitialized sends the initialized notification to the server.
func (c *StdioClient) NotificationsInitialized() error {
	request := jsonrpc.NewRequest(jsonrpc.MethodInitializedNotification, nil, nil)

	c.ForwardRequest(request)

	return nil
}
