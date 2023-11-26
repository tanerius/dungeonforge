package server

import (
	"github.com/tanerius/dungeonforge/pkg/events"
)

type ClientEvent struct {
	clientId string
	client   *Client
	eventId  events.EventIdType
}

func NewClientEvent(_eventId events.EventIdType, _id string, _client *Client) *ClientEvent {
	return &ClientEvent{
		clientId: _id,
		client:   _client,
		eventId:  _eventId,
	}
}

func (e ClientEvent) EventId() events.EventIdType {
	return e.eventId
}

type MessageEvent struct {
	clientId string
	data     []byte
	eventId  events.EventIdType
}

func NewMessageEvent(_eventId events.EventIdType, _id string, _data []byte) *MessageEvent {
	return &MessageEvent{
		clientId: _id,
		data:     _data,
		eventId:  _eventId,
	}
}

func (e *MessageEvent) EventId() events.EventIdType {
	return e.eventId
}
