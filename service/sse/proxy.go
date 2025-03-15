package sse

import (
	"fmt"

	"github.com/chatmcp/mcprouter/service/jsonrpc"
	"github.com/chatmcp/mcprouter/service/mcpclient"
)

func ForwardRequest(client *mcpclient.StdioClient, request *jsonrpc.Request) *jsonrpc.Response {
	fmt.Printf("forward request: %+v\n", request)

	response, err := client.SendRequest(request)
	if err != nil {
		fmt.Printf("failed to forward request: %v\n", err)
		return jsonrpc.NewErrorResponse(jsonrpc.ErrorProxyError, request.ID)
	}

	fmt.Printf("forward response: %+v\n", response)

	return response

	// switch request.Method {
	// case "initialize":
	// 	return HandleInitialize(request)
	// case "notifications/initialized":
	// 	return HandleInitializedNotification(request)
	// case "tools/list":
	// 	return HandleListTools(request)
	// case "tools/call":
	// 	return HandleCallTool(request)
	// case "notifications/cancelled":
	// 	return HandleCancelNotification(request)
	// default:
	// 	return jsonrpc.NewErrorResponse(jsonrpc.ErrorMethodNotFound, request.ID)
	// }
}
