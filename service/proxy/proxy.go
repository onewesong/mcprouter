package proxy

import (
	"fmt"

	"github.com/chatmcp/mcprouter/service/jsonrpc"
)

func HandleRequest(request *jsonrpc.Request) *jsonrpc.Response {
	fmt.Printf("proxy handle request: %+v\n", request)

	switch request.Method {
	case "initialize":
		return HandleInitialize(request)
	case "notifications/initialized":
		return HandleInitializedNotification(request)
	case "tools/list":
		return HandleListTools(request)
	case "tools/call":
		return HandleCallTool(request)
	case "notifications/cancelled":
		return HandleCancelNotification(request)
	default:
		return jsonrpc.NewErrorResponse(jsonrpc.ErrorMethodNotFound, request.ID)
	}
}
