package jsonrpc

import "encoding/json"

// Response is a JSON-RPC response.
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

// String returns the string representation of the response.
func (r *Response) String() string {
	b, _ := json.Marshal(r)

	return string(b)
}

// NewResultResponse creates a new result response.
func NewResultResponse(result interface{}, id interface{}) *Response {
	return &Response{
		JSONRPC: JSONRPC_VERSION,
		Result:  result,
		ID:      id,
	}
}

// NewErrorResponse creates a new error response.
func NewErrorResponse(err *Error, id interface{}) *Response {
	return &Response{
		JSONRPC: JSONRPC_VERSION,
		Error:   err,
		ID:      id,
	}
}
