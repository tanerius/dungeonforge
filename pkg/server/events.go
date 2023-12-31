package server

const (
	EventClientDisconnect int = iota
	EventClientConnect
	EventClientRegistered
	EventMsgReceived
)

type ClientEvent struct {
	clientId string
	client   *client
}

func (c *ClientEvent) ClientId() string {
	return c.clientId
}

func NewClientEvent(_id string, _client *client) *ClientEvent {
	return &ClientEvent{
		clientId: _id,
		client:   _client,
	}
}

type MessageEvent struct {
	clientId string
	data     []byte
}

func NewMessageEvent(_id string, _data []byte) *MessageEvent {
	return &MessageEvent{
		clientId: _id,
		data:     _data,
	}
}

func (e *MessageEvent) Data() []byte {
	return e.data
}

func (e *MessageEvent) ClientId() string {
	return e.clientId
}
