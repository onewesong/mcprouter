package jsonrpc

// Result is a JSON-RPC result.
type Result struct {
}

// InitializeResult is the result of the initialize method.
type InitializeResult struct {
	Result
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      ServerInfo         `json:"serverInfo"`
	Instructions    string             `json:"instructions,omitempty"`
}
