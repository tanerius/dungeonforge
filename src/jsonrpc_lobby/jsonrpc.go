package jsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
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

func NewMockedLobby(logger logging.ILogger) *JsonRpcService {
	l := &JsonRpcService{}
	l.logger = logger
	logger.LogInfo("Setting up IConfig")
	l.conf = config.NewIConfig(true)
	logger.LogInfo("Setting up ID")
	l.id = uuid.NewString()
	root_user, _ := l.conf.ReadKeyString("root_user")
	root_pass, _ := l.conf.ReadKeyString("root_password")
	hostname, _ := l.conf.ReadKeyString("host")
	hostname = fmt.Sprintf("mongodb://%s:%s@%s:27017/", root_user, root_pass, hostname)
	logger.LogInfo("Setting up DB " + hostname)
	if db, err := database.NewMongoDBWrapper(context.Background(), hostname, 100); err == nil {
		l.db = db
		logger.LogInfo("Lobby ready!")
		return l
	}

	return nil
}

func NewLobby(logger logging.ILogger) *JsonRpcService {
	l := &JsonRpcService{}
	l.conf = config.NewIConfig(false)
	l.id = uuid.NewString()
	root_user, _ := l.conf.ReadKeyString("root_user")
	root_pass, _ := l.conf.ReadKeyString("root_password")
	hostname, _ := l.conf.ReadKeyString("host")
	hostname = fmt.Sprintf("mongodb://%s:%s@%s:27017/", root_user, root_pass, hostname)
	logger.LogInfo("Setting up DB " + hostname)
	if db, err := database.NewMongoDBWrapper(context.Background(), hostname, 100); err == nil {
		l.db = db
		return l
	}

	return nil
}

func (r *JsonRpcService) GetId() string {
	return r.id
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

// HandleRPCRequest is the main handler for JSON-RPC requests.
func (s *JsonRpcService) handleRPCRequest(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var result interface{}
	var err error

	switch req.Method {
	case "login":
		result, err = s.login(context.TODO(), req.Params)
	case "register":
		result, err = s.register(context.TODO(), req.Params)
	case "roll":
		result, err = s.roll(context.Background(), req.Params)

	default:
		err = fmt.Errorf("method not found")
	}

	var response Response
	response.Jsonrpc, _ = s.conf.ReadKeyString("apiVersion")
	response.ID = req.ID

	if err != nil {
		response.Error = err.Error()
	} else {
		response.Result = result
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (s *JsonRpcService) Run(port string) {
	http.HandleFunc("/rpc", s.handleRPCRequest)
	fmt.Println("JSON-RPC server listening on port 8080...")
	s.logger.LogInfo("JSON-RPC server listening on port  " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
