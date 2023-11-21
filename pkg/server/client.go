package server

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
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
func (c *Client) readPump(msgChan chan<- []byte) {
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

	for {
		_, message, err := c.cn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				log.Errorf("%s: %v", c.clientId, err)
			}
			// Close the channel to notify the server that there will be nothing to read anymore
			close(msgChan)
			return
		}

		c.lastSeq++
		msgChan <- message
		c.cn.SetReadDeadline(time.Now().Add(15 * time.Second))
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump(toSend <-chan []byte) {

	ticker := time.NewTicker(10 * time.Second)

	defer func() {
		ticker.Stop()
		c.wg.Done()
		log.Debugf("%s write pump stopped.\n", c.clientId)
	}()

	log.Debugf("%s write pump starting...\n", c.clientId)

	for {
		select {
		case message, ok := <-toSend:
			if !ok {
				// close when channel is closed
				log.Debugf("%s sending bye to peer...\n", c.clientId)
				c.cn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.cn.WriteMessage(websocket.BinaryMessage, message); err != nil {
				log.Errorf("%s writing response * %v", c.clientId, err)
				return
			}

			ticker.Reset(10 * time.Second)
			c.cn.SetWriteDeadline(time.Now().Add(12 * time.Second))
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

func (c *Client) ActivateClient(input <-chan []byte, output chan<- []byte) error {
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
