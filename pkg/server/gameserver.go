package server

import "github.com/tanerius/dungeonforge/pkg/messages"

type GameServer interface {
	ResponseChannel() <-chan *messages.Response
	MessageChannel() chan<- *messages.Payload
}
