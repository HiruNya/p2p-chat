package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
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

func server(ctx context.Context, topic *pubsub.Topic, mlog *messageLog, h host.Host, cli CommandLineArguments) {
	http.HandleFunc("/connect", func(writer http.ResponseWriter, request *http.Request) {
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			logger.Errorf("Could not upgrade request to a websocket request: %v", err)
			return
		}
		sub, err := topic.Subscribe()
		if err != nil {
			logger.Errorf("Could not subscribe to topic: %v", err)
			return
		}

		logger.Infof("Connected!")

		go func() {
			// Forward data to client
			defer sub.Cancel()
			conn.WriteJSON(wsMessage{
				Type: JOIN,
				Room: sub.Topic(),
			})
			for {
				messageData, err := sub.Next(ctx)
				if err != nil {
					break
				}
				msg := message{}
				err = json.Unmarshal(messageData.Data, &msg)
				if err != nil {
					logger.Errorf("Could not deserialise json: %v", err)
					break
				}
				name := msg.Name
				if name == "" {
					name = msg.ID[len(msg.ID)-6 : len(msg.ID)]
				}
				err = conn.WriteJSON(wsMessage{
					Type: MESSAGE,
					Text: msg.Text,
					User: name,
					Date: fmt.Sprintf("%02d:%02d", time.Now().Hour(), time.Now().Minute()),
				})
				if err != nil {
					logger.Errorf("Could not send message to client: %v", err)
					break
				}
			}
			logger.Infof("Closed!")
		}()

		// Forward messages from client to pubsub
		wsMsg := wsMessage{}
		for {
			if err = conn.ReadJSON(&wsMsg); err != nil {
				logger.Errorf("Could not read message: %v", err)
				break
			}
			m := message{
				Clock: mlog.clock,
				ID:    peer.Encode(h.ID()),
				Name:  wsMsg.User,
				Text:  wsMsg.Text,
			}
			b, err := json.Marshal(m)
			if err != nil {
				logger.Errorf("Could not marshal message: %v", err)
				continue
			}
			err = topic.Publish(ctx, b)
			if err != nil {
				logger.Errorf("Could not publish message: %v", err)
				continue
			}
		}

	})

	ip := cli.Ip
	if ip == "" {
		ip = "localhost"
	}
	logger.Infof("Serving websocket connection on ws://%s:%d/connect", ip, cli.WebsocketPort)
	logger.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", cli.Ip, cli.WebsocketPort), nil))
}

type wsMessage struct {
	Type string
	Text string
	User string
	Date string
	Room string
}
