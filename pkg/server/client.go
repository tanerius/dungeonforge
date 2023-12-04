package server

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/tanerius/EventGoRound/eventgoround"
)

// Server side representation of the connected client
type client struct {
	clientId     string
	eventManager *eventgoround.EventManager
	cn           *websocket.Conn
	started      bool
	lastSeq      int64
	wg           sync.WaitGroup
	pingTime     time.Time
	sendChannel  chan []byte
}

type clients map[string]*client

func newClient(_c *websocket.Conn, _e *eventgoround.EventManager) *client {
	return &client{
		clientId:     uuid.NewString(),
		eventManager: _e,
		cn:           _c,
		started:      false,
		lastSeq:      0,
		sendChannel:  make(chan []byte),
	}
}

// readPump pumps messages from the websocket connection to the game coordinator.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *client) readPump() {

	defer func() {
		event := eventgoround.NewEvent(EventClientDisconnect, NewClientEvent(c.clientId, c))
		c.eventManager.DispatchPriorityEvent(event)
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
			return
		}

		c.lastSeq++
		event := eventgoround.NewEvent(EventMsgReceived, NewMessageEvent(c.clientId, message))
		c.eventManager.DispatchEvent(event)

		c.cn.SetReadDeadline(time.Now().Add(15 * time.Second))
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *client) writePump() {

	ticker := time.NewTicker(10 * time.Second)

	defer func() {
		ticker.Stop()
		c.wg.Done()
		log.Debugf("%s write pump stopped.\n", c.clientId)
	}()

	log.Debugf("%s write pump starting...\n", c.clientId)

	for {
		select {
		case message, ok := <-c.sendChannel:
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

// shuts down currently connected client
func (c *client) deActivateClient() {
	log.Debugf("%s deactivating... ", c.clientId)
	c.cn.Close()
	close(c.sendChannel)
	c.wg.Wait()
	log.Debugf("%s deactivated ", c.clientId)
}

func (c *client) activateClient() error {
	if c.started {
		return errors.New("client already activated")
	}

	c.started = true
	c.wg.Add(2)
	go c.writePump()
	go c.readPump()

	return nil
}

func (c *client) ID() string {
	return c.clientId
}
