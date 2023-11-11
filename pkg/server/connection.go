package server

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/tanerius/dungeonforge/pkg/messages"
	"google.golang.org/protobuf/proto"
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

func (c *connection) ReadMessage() (*messages.Payload, error) {
	dt, data, err := c.cn.ReadMessage()

	if err != nil {
		log.Error("pipe -> cannot read message from connection")
		return nil, err
	}

	if dt != websocket.BinaryMessage {
		log.Error("pipe -> non binary message received")
		return nil, errors.New("pipe -> received non binary message in pipe")
	}

	var msg *messages.Payload = &messages.Payload{}

	if err = proto.Unmarshal(data, msg); err != nil {
		return nil, errors.New("pipe -> unmashalling error: " + err.Error())
	}

	return msg, nil
}

func (c *connection) WriteMessage(data []byte) error {
	if err := c.cn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		log.Error("pipe -> error writing socket: " + err.Error())
		return err
	}

	return nil
}
