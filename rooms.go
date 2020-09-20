package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/libp2p/go-libp2p-core/host"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"time"
)

type room struct {
	host             host.Host
	discovery        *discovery.RoutingDiscovery
	cancelAdvertiser context.CancelFunc
	mainCtx          context.Context
	topic            *pubsub.Topic
	sub              *pubsub.Subscription
	cancellableCtx   context.Context
	ps               *pubsub.PubSub
	subscribed       chan *websocket.Conn
	mlog             messageLog
}

func newRoom(mainCtx context.Context, host host.Host, table *dht.IpfsDHT, ps *pubsub.PubSub, name string) *room {
	// Advertise the room you wish to join
	topic, err := ps.Join(name)
	if err != nil {
		logger.Fatalf("Could not join pubsub topic: %v", err)
	}
	sub, err := topic.Subscribe()
	if err != nil {
		logger.Fatalf("Could not subscribe to pubsub topic: %v", err)
	}
	roomDiscovery := discovery.NewRoutingDiscovery(table)
	cancellableCtx, cancelAdvertiser := context.WithCancel(mainCtx)
	discovery.Advertise(cancellableCtx, roomDiscovery, name)
	r := room{
		discovery:        roomDiscovery,
		host:             host,
		cancelAdvertiser: cancelAdvertiser,
		cancellableCtx:   cancellableCtx,
		mainCtx:          mainCtx,
		topic:            topic,
		sub:              sub,
		ps:               ps,
		mlog:             messageLog{data: make(map[message]struct{})},
		subscribed:       make(chan *websocket.Conn, 1024),
	}
	go findMembersInRoom(cancellableCtx, host, roomDiscovery, name)
	go r.recv()
	return &r
}

// Join a new room
// All previous goroutines required for advertising the current room will be stopped.
func (r *room) joinRoom(roomName string) {
	if r.topic.String() == roomName {
		return
	}

	// Stops advertising the previous room and starts advertising the new room
	r.sub.Cancel()
	if err := r.topic.Close(); err != nil {
		logger.Errorf("Could not close topic: %v", err)
	}
	r.cancelAdvertiser()
	logger.Infof("Exited the room!")

	t, err := r.ps.Join(roomName)
	r.topic = t
	if err != nil {
		logger.Fatalf("Could not join new topic: %v", err)
	}
	r.sub, err = r.topic.Subscribe()
	if err != nil {
		logger.Fatalf("Could not subscribe to new room: %v", err)
	}
	cancellableCtx, cancelAdvertiser := context.WithCancel(r.mainCtx)
	discovery.Advertise(cancellableCtx, r.discovery, roomName)
	go findMembersInRoom(cancellableCtx, r.host, r.discovery, roomName)
	r.cancelAdvertiser = cancelAdvertiser
	r.cancellableCtx = cancellableCtx

	for i := 0; i < len(r.subscribed); i++ {
		conn := <-r.subscribed
		if err = conn.WriteJSON(wsMessage{
			Type: JOIN,
			Room: r.topic.String(),
		}); err != nil {
			logger.Errorf("Could not send join message to client: %v", err)
			continue
		}
		r.subscribed <- conn
	}
}

func findMembersInRoom(ctx context.Context, host host.Host, roomDiscovery *discovery.RoutingDiscovery, roomName string) {
	peerChannel, err := roomDiscovery.FindPeers(ctx, roomName)
	if err != nil {
		logger.Errorf("Could not find more peers in the room: %v", err)
		return
	}
	for peer := range peerChannel {
		if peer.ID == host.ID() {
			continue
		}
		logger.Infof("Found member in room!")
		if err = host.Connect(ctx, peer); err != nil {
			logger.Errorf("Could not connect to peer found in room: %v", err)
			continue
		}
		logger.Infof("Connected to member in room! %v", peer.Addrs)
	}
}

func (r *room) subscribe(conn *websocket.Conn) {
	if err := conn.WriteJSON(wsMessage{
		Type: JOIN,
		Room: r.topic.String(),
	}); err != nil {
		logger.Errorf("Client could not subscribe: %v", err)
		return
	}
	r.subscribed <- conn
}

func (r *room) recv() {
	for {
		msgData, err := r.sub.Next(r.mainCtx)
		if err != nil {
			logger.Errorf("Could not get next value: %v", err)
			continue
		}
		msg := message{}
		if err = json.Unmarshal(msgData.Data, &msg); err != nil {
			logger.Errorf("Could not deserialise json: %v", err)
			break
		}
		name := msg.Name
		if name == "" {
			name = msg.ID[len(msg.ID)-6 : len(msg.ID)]
		}
		r.mlog.Append(msg)
		subscribed_len := len(r.subscribed)
		for i := 0; i < subscribed_len; i++ {
			conn := <-r.subscribed
			err := conn.WriteJSON(wsMessage{
				Type: MESSAGE,
				Text: msg.Text,
				User: name,
				Date: fmt.Sprintf("%02d:%02d", time.Now().Hour(), time.Now().Minute()),
			})
			if err != nil {
				logger.Errorf("Could not send message to client: %v", err)
			} else {
				r.subscribed <- conn
			}
		}
	}
}

func (r *room) send(msg message) error {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	if err = r.topic.Publish(r.cancellableCtx, bytes); err != nil {
		return err
	}
	return nil
}
