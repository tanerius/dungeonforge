package server

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/tanerius/dungeonforge/pkg/messages"
)

// Server side representation of the connected client
type Client struct {
	clientId        uuid.UUID
	cn              *websocket.Conn
	gameCoordinator *Coordinator
	toSend          chan *messages.Response // responses sent to users
	closeRequested  bool
	mu              sync.Mutex
	lastSeq         int64
	disconnectChan  chan string
}

type clients map[uuid.UUID]*Client

func newClient(_c *websocket.Conn) *Client {
	return &Client{
		clientId:       uuid.New(),
		cn:             _c,
		closeRequested: false,
		lastSeq:        0,
		toSend:         make(chan *messages.Response),
		disconnectChan: make(chan string),
	}
}

// Starts the client read/write pump enabling communication ability
func (c *Client) activateClientOnGameserver(_gameCoordinator *Coordinator) {
	log.Printf("%s activating...\n", c.clientId.String())
	c.gameCoordinator = _gameCoordinator
	go c.writePump()
	go c.readPump()
}

// readPump pumps messages from the websocket connection to the game coordinator.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.gameCoordinator.Unregister <- c
	}()
	//c.cn.SetReadLimit(maxMessageSize)
	//c.cn.SetReadDeadline(time.Now().Add(pongWait))
	//c.cn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	log.Printf("%s read pump starting...\n", c.clientId.String())
	for {
		var message *messages.Payload = &messages.Payload{}
		err := c.cn.ReadJSON(message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				log.Errorf("%s: %v", c.clientId.String(), err)
			}
			break
		}

		// Append the UUID
		message.ClientId = c.clientId

		// check if out of sequence
		if c.lastSeq+1 != message.Seq {
			// out of sequence
			log.Errorf("%s out of sequence: %v expected %d\n", c.clientId.String(), message, c.lastSeq+1)
			break
		} else {
			// SEND THE MESSAGE TO game
			c.lastSeq++
			c.gameCoordinator.PlayerMessagesChan <- message
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	//ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Printf("%s closing connection...\n", c.clientId.String())
		c.mu.Lock()
		if !c.closeRequested {
			c.closeRequested = true
			c.gameCoordinator.Unregister <- c
		}
		c.mu.Unlock()
		log.Printf("%s closed\n", c.clientId.String())
	}()

	log.Printf("%s write pump starting...\n", c.clientId.String())

	for {
		select {
		case message, ok := <-c.toSend:
			if !ok {
				log.Printf("%s sending channel closed\n", c.clientId.String())
				// The hub closed the channel.
				cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye")
				c.cn.WriteMessage(websocket.CloseMessage, cm)
				return
			}

			if err := c.cn.WriteJSON(message); err != nil {
				log.Errorf("%s writing response * %v", c.clientId.String(), err)
				return
			}
		case message, _ := <-c.disconnectChan:
			// The close connection
			cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, message)
			c.cn.WriteMessage(websocket.CloseMessage, cm)
			return
		}
	}
}
