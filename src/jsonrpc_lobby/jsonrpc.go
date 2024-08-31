package jsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

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

// Good to use "/rpc" as endpoint
func (s *JsonRpcService) Run(endpoint, port string) {
	mux := http.NewServeMux()
	mux.HandleFunc(endpoint, s.handleRPCRequest)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Channel to listen for interrupt or termination signals.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Run the server in a goroutine so that it doesn't block.
	go func() {
		s.logger.LogInfo("JSON-RPC server listening on port  " + port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.LogError(err, "Could not listen on port  "+port)
		}
	}()

	// Wait for an interrupt signal.
	<-stop

	// Gracefully shutdown the server with a timeout of 5 seconds.
	s.logger.LogInfo("Shutting down the server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		s.logger.LogError(err, "Server foced to shutdown")
	}

	s.logger.LogInfo("Server exiting")
}
