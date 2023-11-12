package server

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type client struct {
	entityId uuid.UUID
	cn       *websocket.Conn
}

type clients map[uuid.UUID]*client

func newConnection(_c *websocket.Conn) *client {
	return &client{
		entityId: uuid.New(),
		cn:       _c,
	}
}
