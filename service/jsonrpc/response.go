package jsonrpc

import "encoding/json"

type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

func (r *Response) String() string {
	b, _ := json.Marshal(r)

	return string(b)
}

func NewResultResponse(result interface{}, id interface{}) *Response {
	return &Response{
		JSONRPC: JSONRPC_VERSION,
		Result:  result,
		ID:      id,
	}
}

func NewErrorResponse(err *Error, id interface{}) *Response {
	return &Response{
		JSONRPC: JSONRPC_VERSION,
		Error:   err,
		ID:      id,
	}
}
