package proxy

import "time"

type ProxyInfo struct {
	JSONRPCVersion     string        `json:"jsonrpc_version"`
	ProtocolVersion    string        `json:"protocol_version"`
	ClientName         string        `json:"client_name"`
	ClientVersion      string        `json:"client_version"`
	RequestMethod      string        `json:"request_method"`
	RequestParams      interface{}   `json:"request_params"`
	RequestID          interface{}   `json:"request_id"`
	RequestTime        time.Time     `json:"request_time"`
	SSERequestTime     time.Time     `json:"sse_request_time"`
	SessionID          string        `json:"session_id"`
	ServerUUID         string        `json:"server_uuid"`
	ServerKey          string        `json:"server_key"`
	ServerConfigName   string        `json:"server_config_name"`
	ServerShareProcess bool          `json:"server_share_process"`
	ServerCommand      string        `json:"server_command"`
	ServerCommandHash  string        `json:"server_command_hash"`
	ServerName         string        `json:"server_name"`
	ServerVersion      string        `json:"server_version"`
	ResponseTime       time.Time     `json:"response_time"`
	ResponseDuration   time.Duration `json:"response_duration"`
	ResponseResult     string        `json:"response_result"`
	ResponseError      string        `json:"response_error"`
}
