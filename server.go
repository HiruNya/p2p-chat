package main

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const MESSAGE = "MESSAGE"
const JOIN = "JOIN"

func server(ctx context.Context, h host.Host, cli CommandLineArguments, ps *pubsub.PubSub, currentRoom *room) {
	http.HandleFunc("/connect", func(writer http.ResponseWriter, request *http.Request) {
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			logger.Errorf("Could not upgrade request to a websocket request: %v", err)
			return
		}
		logger.Infof("Connected!")

		go handleWebsocketConnection(conn, h, currentRoom)
	})

	ip := cli.Ip
	if ip == "" {
		ip = "localhost"
	}
	logger.Infof("Serving websocket connection on ws://%s:%d/connect", ip, cli.WebsocketPort)
	if cli.Https {
		go func() {
			logger.Fatal(http.ListenAndServeTLS(":443", "domain.cert.pem", "private.key.pem", nil))
		}()
	}
	logger.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", cli.Ip, cli.WebsocketPort), nil))
}

func handleWebsocketConnection(conn *websocket.Conn, h host.Host, currentRoom *room) {
	currentRoom.joinRoom(currentRoom.topic.String())
	currentRoom.subscribe(conn)

	// Forward messages from client to pubsub
	wsMsg := wsMessage{}
	for {
		if err := conn.ReadJSON(&wsMsg); err != nil {
			logger.Errorf("Could not read message: %v", err)
			break
		}
		if wsMsg.Type == JOIN {
			currentRoom.joinRoom(wsMsg.Room)
			continue
		}
		m := message{
			Clock: currentRoom.mlog.clock,
			ID:    peer.Encode(h.ID()),
			Name:  wsMsg.User,
			Text:  wsMsg.Text,
		}
		if err := currentRoom.send(m); err != nil {
			logger.Errorf("Could not send message: %v", err)
		}
	}
}

type wsMessage struct {
	Type string
	Text string
	User string
	Date string
	Room string
}
