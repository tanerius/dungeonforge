package server

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type connection struct {
	entityId         uuid.UUID
	cn               *websocket.Conn
	expectedInterval int32
}

type connections map[uuid.UUID]*connection

func newConnection(_c *websocket.Conn) *connection {
	return &connection{
		entityId:         uuid.New(),
		cn:               _c,
		expectedInterval: 1,
	}
}
