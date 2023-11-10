package main

import (
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

	msg := Client{
		Name: "Taner",
		Id:   12,
	}

	if err := conn.WriteJSON(msg); err != nil {
		log.Fatal(err)
	}

	log.Println("Client succesfully connected!")
}
