package mcpclient

import (
	"encoding/json"

	"github.com/chatmcp/mcprouter/service/jsonrpc"
)

// ListTools lists the tools available in the MCP server.
func (c *StdioClient) ListTools() (*jsonrpc.ListToolsResult, error) {
	request := jsonrpc.NewRequest(jsonrpc.MethodListTools, nil, 1)
	response := c.ForwardRequest(request)

	resultB, err := json.Marshal(response.Result)
	if err != nil {
		return nil, err
	}

	result := &jsonrpc.ListToolsResult{}
	if err := json.Unmarshal(resultB, result); err != nil {
		return nil, err
	}

	return result, nil
}

// CallTool calls a tool with the given name and arguments.
func (c *StdioClient) CallTool(params *jsonrpc.CallToolParams) (*jsonrpc.CallToolResult, error) {
	request := jsonrpc.NewRequest(jsonrpc.MethodCallTool, params, 1)
	response := c.ForwardRequest(request)

	resultB, err := json.Marshal(response.Result)
	if err != nil {
		return nil, err
	}

	result := &jsonrpc.CallToolResult{}
	if err := json.Unmarshal(resultB, result); err != nil {
		return nil, err
	}

	return result, nil
}
