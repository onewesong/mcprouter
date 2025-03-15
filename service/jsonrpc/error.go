package jsonrpc

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewError(code int, message string, data interface{}) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

var (
	ErrorParseError     = NewError(-32700, "Parse error", nil)
	ErrorInvalidRequest = NewError(-32600, "Invalid Request", nil)
	ErrorMethodNotFound = NewError(-32601, "Method not found", nil)
	ErrorInvalidParams  = NewError(-32602, "Invalid params", nil)
	ErrorInternalError  = NewError(-32603, "Internal error", nil)
)
