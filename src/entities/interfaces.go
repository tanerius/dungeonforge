package entities

import "context"

type Entity interface {
	GetId() string
}

// Wraps writing to DB functions for entities
type EntityWriter interface {
	Write(context.Context) error
}

type LobbyServiceController interface {
	GetId() string
	// What type of lobby grpc of jsonrpc
	GetType() string
	// This function is called when a lobby should be started. True if started
	Start() bool
	// This function is called when a lobby should be started.
	Stop()
	// Get status of the LobbyService
	Status() string
}
