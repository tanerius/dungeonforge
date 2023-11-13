package server

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/tanerius/dungeonforge/pkg/messages"
)

// Server side representation of the connected client
type client struct {
	clientId        uuid.UUID
	cn              *websocket.Conn
	gameCoordinator *Coordinator
	toSend          chan *messages.Response // responses sent to users
	closeRequested  bool
	mu              sync.Mutex
	isValidated     bool
}

type clients map[uuid.UUID]*client

func newConnection(_c *websocket.Conn) *client {
	return &client{
		clientId:       uuid.New(),
		cn:             _c,
		closeRequested: false,
		isValidated:    false,
	}
}

// Starts the client read/write pump enabling communication ability
func (c *client) activateClientOnGameserver() {
	go c.writePump()
	go c.readPump()
}

// use to send messages to client from outside
// TODO: use context here to timeout the messasge
func (c *client) SendMessage(_m *messages.Response) {
	if !c.isValidated {
		return
	}

	c.toSend <- _m
}

// readPump pumps messages from the websocket connection to the game coordinator.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *client) readPump() {
	defer func() {
		c.mu.Lock()
		if !c.closeRequested {
			c.closeRequested = true
			c.gameCoordinator.unregister <- c
		}
		c.mu.Unlock()
	}()
	//c.cn.SetReadLimit(maxMessageSize)
	//c.cn.SetReadDeadline(time.Now().Add(pongWait))
	//c.cn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var message *messages.Payload = &messages.Payload{}
		err := c.cn.ReadJSON(message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("client error: %v", err)
			}
			break
		}

		if !c.isValidated {
			if message.Cmd != messages.CmdValidate || message.Seq != 1 {
				break // disconnect this client
			}
			// TODO: send to a validation channel on coordinator
		}

		// Append the UUID
		message.ClientId = c.clientId

		// SEND THE MESSAGE TO game
		c.gameCoordinator.playerMessages <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *client) writePump() {
	//ticker := time.NewTicker(pingPeriod)
	defer func() {
		c.mu.Lock()
		if !c.closeRequested {
			c.closeRequested = true
			c.gameCoordinator.unregister <- c
		}
		c.mu.Unlock()
	}()

	for {
		select {
		case message, ok := <-c.toSend:
			// Never sent messages to an unvalidated client
			if c.isValidated {

				if !ok {
					// The hub closed the channel.
					cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye")
					c.cn.WriteMessage(websocket.CloseMessage, cm)
					return
				}

				if err := c.cn.WriteJSON(message); err != nil {
					log.Error(err)
					return
				}
			}
		}
	}
}
