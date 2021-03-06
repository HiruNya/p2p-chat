package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p"
	relay "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/routing"
	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-pubsub"
)

var logger = log.Logger("group-chat")

func main() {
	var _ = log.SetLogLevel("group-chat", "info")
	cliArgs := ParseCliArgs()
	logger.Infof("Nickname: \"%s\"", cliArgs.Name)

	// Start host
	ctx := context.Background()
	enableRelay := libp2p.EnableRelay()
	if cliArgs.Relay {
		enableRelay = libp2p.EnableRelay(relay.OptHop)
	}
	h, err := libp2p.New(ctx,
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			return dht.New(ctx, h)
		}),
		libp2p.EnableNATService(),
		enableRelay,
		libp2p.EnableAutoRelay(),
	)
	if err != nil {
		logger.Fatalf("Could not start lbp2p host: %v", err)
	}

	// DHT for Rooms
	roomDht, err := dht.New(ctx, h,
		dht.BootstrapPeers(cliArgs.Bootstrap.toPeerAddrArray()...),
		dht.ProtocolPrefix("/hiruchat"))
	if err != nil {
		logger.Fatalf("Could not create DHT for rooms: %v", roomDht)
	}
	if err := roomDht.Bootstrap(ctx); err != nil {
		logger.Fatalf("Failed to bootstrap the room table: %v", err)
	}

	// Subscribe to messages
	ps, err := pubsub.NewFloodSub(ctx, h)
	if err != nil {
		logger.Fatalf("Could not start pubsub: %v", err)
	}

	currentRoom := newRoom(ctx, h, roomDht, ps, cliArgs.Room)

	// Connect to bootstrap peers
	for _, b := range cliArgs.Bootstrap {
		if err := h.Connect(ctx, *b); err != nil {
			logger.Errorf("Could not connect to bootstrap node: %v", err)
		}
	}

	// We need to wait 15 minutes for libp2p to advertise this host
	if cliArgs.Relay {
		fmt.Println("Relay peers must wait 15 minutes before use.")
		time.Sleep(15 * time.Minute)
		fmt.Println("Ready to go!")
		for _, a := range h.Addrs() {
			if strings.Contains(a.String(), "127.0.0.1") {
				fmt.Printf("%v/p2p/%v\n", strings.Replace(a.String(), "127.0.0.1", cliArgs.Ip, 1), h.ID())
			}
		}
	}

	logger.Infof("PeerID: %s", h.ID().String())

	go server(ctx, h, cliArgs, ps, currentRoom)

	<-ctx.Done() // Wait for the program to close
}

func (mlog *messageLog) Append(msg message) {
	mlog.mux.Lock()
	defer mlog.mux.Unlock()
	if _, ok := mlog.data[msg]; ok {
		// Message already exists
		return
	}
	name := msg.Name
	if name == "" {
		// Use the last 6 characters of the peer's address if no nickname is provided
		name = msg.ID[len(msg.ID)-6 : len(msg.ID)]
	}
	logger.Infof("%s:\t%s", name, msg.Text)
	if msg.Clock >= mlog.clock {
		mlog.clock = msg.Clock + 1
	}
}

type messageLog struct {
	mux   sync.Mutex
	data  map[message]struct{}
	clock uint
}

type message struct {
	Clock uint
	ID    string // The peer ID
	Name  string
	Text  string
}
