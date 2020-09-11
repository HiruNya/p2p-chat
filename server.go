package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const MESSAGE = "MESSAGE"

func server() {
	http.HandleFunc("/connect", func(writer http.ResponseWriter, request *http.Request) {
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			logger.Errorf("Could not upgrade request to a websocket request: %v", err)
		}
		logger.Infof("Connected!")

		i := 0
		for {
			err = conn.WriteJSON(WsMessage{
				Type: MESSAGE,
				Text: fmt.Sprintf("This is message %d", i),
				User: fmt.Sprintf("User%d", i),
				Date: fmt.Sprintf("%02d:%02d", i, i),
			})
			logger.Infof("Written!")
			if err != nil {
				logger.Errorf("Could not send JSON message: %v", err)
				break
			}
			time.Sleep(time.Second * 10)
			i++
		}
	})

	logger.Fatal(http.ListenAndServe(":8000", nil))
}

type WsMessage struct {
	Type string
	Text string
	User string
	Date string
}
