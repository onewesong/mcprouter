package jsonrpc

import "encoding/json"

type Request struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	ID      interface{} `json:"id,omitempty"`
}

func UnmarshalRequest(data []byte) (*Request, error) {
	var j Request

	if err := json.Unmarshal(data, &j); err != nil {
		return nil, err
	}

	return &j, nil
}
