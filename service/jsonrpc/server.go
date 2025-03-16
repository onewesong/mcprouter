package jsonrpc

// ServerCapabilitiesTools is the tools of the server capabilities.
type ServerCapabilitiesTools struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// ServerCapabilities is the capabilities of the server.
type ServerCapabilities struct {
	Experimental map[string]interface{} `json:"experimental,omitempty"`
	Logging      *struct{}              `json:"logging,omitempty"`
	Prompts      *struct {
		ListChanged bool `json:"listChanged,omitempty"`
	} `json:"prompts,omitempty"`
	Resources *struct {
		Subscribe bool `json:"subscribe,omitempty"`
	} `json:"resources,omitempty"`
	Tools ServerCapabilitiesTools `json:"tools,omitempty"`
}

// ServerInfo is the info of the server.
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
