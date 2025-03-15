package proxy

import (
	"fmt"

	"github.com/chatmcp/mcprouter/service/jsonrpc"
)

func HandleInitializedNotification(request *jsonrpc.Request) *jsonrpc.Response {
	fmt.Printf("proxy handle initialized notification: %+v\n", request)

	return nil
}

func HandleCancelNotification(request *jsonrpc.Request) *jsonrpc.Response {
	fmt.Printf("proxy handle cancel notification: %+v\n", request)

	return jsonrpc.NewResultResponse(nil, request.ID)
}
