package server

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/tanerius/dungeonforge/pkg/messages"
)

// Server side representation of the connected client
type Client struct {
	clientId         string
	cn               *websocket.Conn
	started          bool
	lastSeq          int64
	ConnectionFailed chan struct{}
	wg               sync.WaitGroup
	pingTime         time.Time
}

type clients map[string]*Client

func newClient(_c *websocket.Conn) *Client {
	return &Client{
		clientId:         uuid.NewString(),
		cn:               _c,
		started:          false,
		lastSeq:          0,
		ConnectionFailed: make(chan struct{}),
	}
}

// readPump pumps messages from the websocket connection to the game coordinator.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump(msgChan chan<- *messages.Request) {
	defer func() {
		c.wg.Done()
		log.Debugf("%s read pump stopped.\n", c.clientId)
	}()

	log.Debugf("%s read pump starting...\n", c.clientId)

	// if nothing is received in 15 seconds then kill connection
	c.cn.SetReadDeadline(time.Now().Add(15 * time.Second))

	c.cn.SetPongHandler(func(string) error {
		// Reset read timer since pong was sent
		duration := time.Since(c.pingTime)
		log.Debugf("%s PING %dms \n", c.clientId, duration.Milliseconds())
		c.cn.SetReadDeadline(time.Now().Add(15 * time.Second))
		return nil
	})

	var isLost bool = false

	for {
		var message *messages.Request = &messages.Request{}
		isLost = false

		err := c.cn.ReadJSON(message)

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				log.Errorf("%s: %v", c.clientId, err)
			}
			isLost = true
			message.CmdType = messages.CmdLost
		}

		// Append the UUID
		message.ClientId = c.clientId
		c.lastSeq++
		msgChan <- message
		if isLost {
			break
		}
		c.cn.SetReadDeadline(time.Now().Add(15 * time.Second))
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump(toSend <-chan *messages.Response) {

	ticker := time.NewTicker(10 * time.Second)

	defer func() {
		ticker.Stop()
		c.wg.Done()
		log.Debugf("%s write pump stopped.\n", c.clientId)
	}()

	log.Debugf("%s write pump starting...\n", c.clientId)

	for {
		select {
		case message := <-toSend:

			c.cn.SetWriteDeadline(time.Now().Add(12 * time.Second))
			if message.Cmd == messages.CmdDisconnect {
				cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, message.Msg)
				if err := c.cn.WriteMessage(websocket.CloseMessage, cm); err != nil {
					log.Errorf("%s writing close message * %v", c.clientId, err)
					return
				}
			} else {
				if err := c.cn.WriteJSON(message); err != nil {
					log.Errorf("%s writing response * %v", c.clientId, err)
					return
				}
			}

			ticker.Reset(10 * time.Second)

		case <-ticker.C:
			if err := c.cn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
			c.pingTime = time.Now()
			c.cn.SetWriteDeadline(time.Now().Add(12 * time.Second))
		}

	}
}

func (c *Client) DeActivateClient() {
	log.Debugf("%s deactivating... ", c.clientId)
	c.cn.Close()
	c.wg.Wait()
	log.Debugf("%s deactivated ", c.clientId)
}

func (c *Client) ActivateClient(input <-chan *messages.Response, output chan<- *messages.Request) error {
	if c.started {
		return errors.New("client already activated")
	}

	c.started = true
	c.wg.Add(2)
	go c.writePump(input)
	go c.readPump(output)

	return nil
}

func (c *Client) ID() string {
	return c.clientId
}
