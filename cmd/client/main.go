package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	game "github.com/tanerius/dungeonforge/pkg/game/messaes"
	"github.com/tanerius/dungeonforge/pkg/messages"
)

const wsServerEndpoont = "ws://localhost:40000/ws"

func main() {
	sendChan := make(chan []byte)
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

				if err := conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
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

				conn.SetPingHandler(func(a string) error {
					log.Printf("received PING from server %s \n\n", wsServerEndpoont)
					conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
					if err := conn.WriteMessage(websocket.PongMessage, nil); err != nil {
						return err
					}
					return nil
				})

				log.Printf("Client * Client connected to server %s \n\n", wsServerEndpoont)

				// starting reader too
				go func() {
					for {
						var message *messages.Response = &messages.Response{}
						_, data, err := conn.ReadMessage()
						if err != nil {
							if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
								log.Errorf("Client error * : %v", err)
							} else {
								log.Infof("Client * : %v", err)
							}

							conn = nil
							seq = 1
							break
						}

						// unmarshal
						if err := json.Unmarshal(data, message); err != nil {
							log.Errorf("Client reader cannot unmarshal * %v", err)
						} else {
							log.Printf("Client * Received: %v \n", message)
						}

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
		fmt.Print("2. Send basic messages.Request message \n")
		fmt.Print("3. Close connection \n")
		fmt.Print("4. Request disconnect from server \n")
		fmt.Print("5. Send Login Request \n")
		fmt.Print("6. Test marshaling \n")
		fmt.Scanln(&i)

		if i == 0 {
			quit = true
			quitChan <- true
		} else if i == 1 {
			connectChan <- true
		} else if i == 2 {
			var YData *messages.Request = &messages.Request{
				CmdType:  messages.CmdExec,
				Seq:      1,
				DataType: 0,
			}
			data, err := json.Marshal(YData)

			if err != nil {
				log.Errorf("Client cannot marshal : %v", err)
			} else {
				sendChan <- data
			}
		} else if i == 3 {
			disconnectChan <- true
		} else if i == 4 {
			var YData *messages.Request = &messages.Request{
				CmdType:  messages.CmdDisconnect,
				Seq:      1,
				DataType: game.TypeNothing,
			}

			data, err := json.Marshal(YData)

			if err != nil {
				log.Errorf("Client cannot marshal : %v", err)
			} else {
				log.Infoln("Client sending a disconnect request to server")
				sendChan <- data
			}

		} else if i == 5 {
			var YData = &game.RequestLogin{
				Request: messages.Request{
					CmdType:  messages.CmdExec,
					Seq:      1,
					DataType: int(game.TypeLogin),
				},
				PlayerId: "tanerius@gmail.com",
				Password: "123123123",
			}

			data, err := json.Marshal(YData)

			if err != nil {
				log.Errorf("Client cannot marshal : %v", err)
			} else {
				log.Infof("Client sending %v", data)
				sendChan <- data
			}
		} else if i == 6 {
			var YData = &game.RequestLogin{
				Request: messages.Request{
					CmdType:  messages.CmdExec,
					Seq:      1,
					DataType: int(game.TypeLogin),
				},
				PlayerId: "tanerius@gmail.com",
				Password: "123123123",
			}

			var result *messages.Request = &messages.Request{}

			data, err := json.Marshal(YData)

			if err != nil {
				log.Errorf("Client cannot marshal : %v", err)
			} else {
				if err := json.Unmarshal(data, result); err != nil {
					log.Errorf("Client cannot unmarshal : %v", err)
				} else {
					log.Infof("Client marshalling worked : %v", result)
				}
			}
		}

	}
}
