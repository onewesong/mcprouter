package mcpserver

// ServerConfig is the config for the remote mcp server
type ServerConfig struct {
	ServerUUID   string `json:"server_uuid,omitempty"`
	ServerName   string `json:"server_name,omitempty"`
	ServerKey    string `json:"server_key,omitempty"`
	Command      string `json:"command"`
	CommandHash  string `json:"command_hash,omitempty"`
	ShareProcess bool   `json:"share_process,omitempty"`
	ServerType   string `json:"server_type,omitempty"`
	ServerURL    string `json:"server_url,omitempty"`
	ServerParams string `json:"server_params,omitempty"`
}
