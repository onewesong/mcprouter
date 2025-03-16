package jsonrpc

import "encoding/json"

// Request is a JSON-RPC request.
type Request struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	ID      interface{} `json:"id,omitempty"`
}

// UnmarshalRequest unmarshals a JSON-RPC request.
func UnmarshalRequest(data []byte) (*Request, error) {
	var j Request

	if err := json.Unmarshal(data, &j); err != nil {
		return nil, err
	}

	return &j, nil
}
