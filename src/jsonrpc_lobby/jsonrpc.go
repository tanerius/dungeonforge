package jsonrpc

import (
	"encoding/json"

	"github.com/tanerius/dungeonforge/pkg/config"
	"github.com/tanerius/dungeonforge/pkg/database"
	"github.com/tanerius/dungeonforge/pkg/logging"
)

// Request defines the structure of a JSON-RPC request.
type Request struct {
	Jsonrpc string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  json.RawMessage  `json:"params"`
	ID      *json.RawMessage `json:"id"`
}

// Response defines the structure of a JSON-RPC response.
type Response struct {
	Jsonrpc string           `json:"jsonrpc"`
	Code    int              `json:"code"`
	Result  interface{}      `json:"result,omitempty"`
	Error   interface{}      `json:"error,omitempty"`
	ID      *json.RawMessage `json:"id"`
}

type JsonRpcService struct {
	conf   config.IConfig
	id     string
	db     *database.MongoDB
	logger logging.ILogger
}

func getResponder(rpcVer string, id *json.RawMessage, code int, resp interface{}) *Response {
	return &Response{
		Jsonrpc: rpcVer,
		ID:      id,
		Code:    code,
		Result:  resp,
	}
}

func getError(rpcVer string, id *json.RawMessage, code int, resp interface{}) *Response {
	return &Response{
		Jsonrpc: rpcVer,
		ID:      id,
		Code:    code,
		Error:   resp,
	}
}
