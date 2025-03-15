package proxy

import "github.com/chatmcp/mcprouter/service/jsonrpc"

func HandleInitialize(request *jsonrpc.Request) *jsonrpc.Response {
	// todo: request backend server

	result := jsonrpc.InitializeResult{
		ProtocolVersion: jsonrpc.LATEST_PROTOCOL_VERSION,
		Capabilities: jsonrpc.ServerCapabilities{
			Tools: jsonrpc.ServerCapabilitiesTools{
				ListChanged: true,
			},
		},
		ServerInfo: jsonrpc.ServerInfo{
			Name:    jsonrpc.PROXY_SERVER_NAME,
			Version: jsonrpc.PROXY_SERVER_VERSION,
		},
		Instructions: "",
	}

	return jsonrpc.NewResultResponse(result, request.ID)
}
