package main

import (
	"fmt"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/tanerius/dungeonforge/pkg/messages"
)

const wsServerEndpoont = "ws://localhost:40000/ws"

func main() {
	sendChan := make(chan *messages.Payload)
	quitChan := make(chan bool)
	connectChan := make(chan bool)
	disconnectChan := make(chan bool)
	var conn *websocket.Conn = nil
	var quit bool = false
	var seq int64 = 1

	go func() {
	L:
		for {
			select {
			case msg := <-sendChan:
				if conn == nil {
					log.Errorln("Client * Cannot send to a nil conn")
					continue
				}
				if err := conn.WriteJSON(msg); err != nil {
					log.Errorf("Client * %v", err)
					continue
				}
				seq++
			case <-disconnectChan:
				if conn == nil {
					log.Errorln("Connection cannot be closed on a nil conn")
					continue
				}
				cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye")
				if err := conn.WriteMessage(websocket.CloseMessage, cm); err != nil {
					log.Errorf("Client *  %v ", err)
				}
				conn.Close()
				log.Println("Connection closed")
				conn = nil
				seq = 1
			case <-quitChan:
				break L
			case <-connectChan:
				if conn != nil {
					log.Errorln("Client * Already appears connected")
					continue
				}
				var err error
				dialer := websocket.Dialer{
					ReadBufferSize:  1024,
					WriteBufferSize: 1024,
				}
				conn, _, err = dialer.Dial(wsServerEndpoont, nil)

				if err != nil {
					log.Errorf("%v", err)
					continue
				}

				log.Printf("Client * Client connected to server %s \n\n", wsServerEndpoont)

				// starting reader too
				go func() {
					for {
						var message *messages.Response = &messages.Response{}
						err := conn.ReadJSON(message)
						if err != nil {
							if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
								log.Errorf("Client reader * %v", err)
							} else {
								log.Printf("Client reader * server said: %v", err)
							}
							conn = nil
							seq = 1
							break
						}

						log.Printf("Client * Received: %v \n", message)
					}
				}()
			}
		}
	}()

	var i int = -1

	for !quit {
		fmt.Print("\n\nWelcome to server test client\n\n")
		fmt.Print("Choices: \n")
		fmt.Print("0. Quit \n")
		fmt.Print("1. Establish connection \n")
		fmt.Print("2. Send a random message \n")
		fmt.Print("3. Close connection \n")
		fmt.Scanln(&i)

		if i == 0 {
			quit = true
			quitChan <- true
		} else if i == 1 {
			connectChan <- true
		} else if i == 2 {
			var YData *messages.Payload = &messages.Payload{
				Token: "y",
				Seq:   seq,
				Cmd:   messages.CmdValidate,
				Data: messages.PersonJson{
					Name: "Tanerius",
					Age:  45,
				},
			}
			sendChan <- YData
		} else if i == 3 {
			disconnectChan <- true
		}

	}
}
