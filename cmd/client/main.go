package main

import (
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const wsServerEndpoont = "ws://localhost:40000/ws"

type Client struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

func main() {
	dialer := websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, _, err := dialer.Dial(wsServerEndpoont, nil)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	/*
		msg := &messages.Person{
			Name: *proto.String("Tanerius"),
			Age:  *proto.Int32(45),
		}

		log.Println(msg.String())

		data, _ := proto.Marshal(msg)

		log.Printf("Marshalled data: %s ", data)

		if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
			log.Fatal(err)
		}
	*/

	duration := time.Duration(3) * time.Second
	time.Sleep(duration)

	log.Println("Client succesfully connected!")
}
