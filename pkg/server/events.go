package server

import (
	"github.com/tanerius/dungeonforge/pkg/events"
)

type ClientConnectedEvent struct {
	clientId *Client
}

func NewClientConnectedEvent(_client *Client) ClientConnectedEvent {
	return ClientConnectedEvent{
		clientId: _client,
	}
}

func (e ClientConnectedEvent) EventId() events.EventIdType {
	return events.EventClientConnected
}

type ClientDisonnectedEvent struct {
	clientId string
}

func NewClientDisonnectedEvent(_id string) ClientDisonnectedEvent {
	return ClientDisonnectedEvent{
		clientId: _id,
	}
}

func (e ClientDisonnectedEvent) EventId() events.EventIdType {
	return events.EventClientDisconnected
}
