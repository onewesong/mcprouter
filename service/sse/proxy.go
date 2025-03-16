package sse

import (
	"fmt"

	"github.com/chatmcp/mcprouter/service/jsonrpc"
	"github.com/chatmcp/mcprouter/service/mcpclient"
)

// ForwardRequest forwards a JSON-RPC request to the MCP server and returns the response
func ForwardRequest(client *mcpclient.StdioClient, request *jsonrpc.Request) *jsonrpc.Response {
	fmt.Printf("forward request: %+v\n", request)

	response, err := client.SendRequest(request)
	if err != nil {
		fmt.Printf("failed to forward request: %v\n", err)
		return jsonrpc.NewErrorResponse(jsonrpc.ErrorProxyError, request.ID)
	}

	fmt.Printf("forward response: %+v\n", response)

	return response
}
