package proxy

import (
	"fmt"

	"github.com/chatmcp/mcprouter/service/jsonrpc"
)

func HandleListTools(request *jsonrpc.Request) *jsonrpc.Response {
	fmt.Printf("proxy handle list tools: %+v\n", request)

	tools := []map[string]interface{}{
		{
			"name":        "ping",
			"description": "Ping the server",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"message": map[string]interface{}{
						"type":        "string",
						"description": "Message to ping the server",
					},
				},
				"required": []string{"message"},
			},
		},
	}

	result := map[string]interface{}{
		"tools": tools,
	}

	return jsonrpc.NewResultResponse(result, request.ID)
}

func HandleCallTool(request *jsonrpc.Request) *jsonrpc.Response {
	fmt.Printf("proxy handle call tool: %+v,\n", request.Params)

	result := map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": "pong",
			},
		},
	}

	return jsonrpc.NewResultResponse(result, request.ID)
}
