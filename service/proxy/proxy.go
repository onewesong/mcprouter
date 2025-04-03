package proxy

import "time"

type ProxyInfo struct {
	JSONRPCVersion     string      `json:"jsonrpc_version"`
	ProtocolVersion    string      `json:"protocol_version"`
	ConnectionTime     time.Time   `json:"connection_time"`
	ClientName         string      `json:"client_name"`
	ClientVersion      string      `json:"client_version"`
	RequestMethod      string      `json:"request_method"`
	RequestParams      interface{} `json:"request_params"`
	RequestID          interface{} `json:"request_id"`
	RequestTime        time.Time   `json:"request_time"`
	RequestFrom        string      `json:"request_from"`
	SessionID          string      `json:"session_id"`
	ServerUUID         string      `json:"server_uuid"`
	ServerKey          string      `json:"server_key"`
	ServerConfigName   string      `json:"server_config_name"`
	ServerShareProcess bool        `json:"server_share_process"`
	ServerType         string      `json:"server_type"`
	ServerURL          string      `json:"server_url"`
	ServerCommand      string      `json:"server_command"`
	ServerCommandHash  string      `json:"server_command_hash"`
	ServerName         string      `json:"server_name"`
	ServerVersion      string      `json:"server_version"`
	ResponseTime       time.Time   `json:"response_time"`
	ResponseResult     interface{} `json:"response_result"`
	ResponseError      string      `json:"response_error"`
	CostTime           int64       `json:"cost_time"`
}
