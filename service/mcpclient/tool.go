package mcpclient

import (
	"github.com/chatmcp/mcprouter/service/jsonrpc"
)

// ListTools lists the tools available in the MCP server.
func (c *StdioClient) ListTools() (*jsonrpc.ListToolsResult, error) {
	request := jsonrpc.NewRequest(jsonrpc.MethodListTools, nil, 1)

	response, err := c.ForwardMessage(request)
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, response.Error
	}

	result := &jsonrpc.ListToolsResult{}
	if err := response.UnmarshalResult(result); err != nil {
		return nil, err
	}

	return result, nil
}

// CallTool calls a tool with the given name and arguments.
func (c *StdioClient) CallTool(params *jsonrpc.CallToolParams) (*jsonrpc.CallToolResult, error) {
	request := jsonrpc.NewRequest(jsonrpc.MethodCallTool, params, 1)

	response, err := c.ForwardMessage(request)
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, response.Error
	}

	result := &jsonrpc.CallToolResult{}
	if err := response.UnmarshalResult(result); err != nil {
		return nil, err
	}

	return result, nil
}
