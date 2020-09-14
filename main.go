package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	"os"
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

const topic = "sylo-group-chat-demo"

func main() {
	var _ = log.SetLogLevel("group-chat", "info")
	cliArgs := ParseCliArgs()
	logger.Infof("Nickname: \"%s\"", cliArgs.Name)
	mlog := messageLog{}
	mlog.data = make(map[sentMessage]struct{})
	mlog.history = make([]message, 0, 10)

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

	// Subscribe to messages
	ps, err := pubsub.NewFloodSub(ctx, h)
	if err != nil {
		logger.Fatalf("Could not start pubsub: %v", err)
	}
	t, err := ps.Join(topic)
	if err != nil {
		logger.Fatalf("Could not join pubsub topic: %v", err)
	}
	sub, err := t.Subscribe()
	if err != nil {
		logger.Fatalf("Could not subscribe to topic: %v", err)
	}

	// Spawn a goroutine to handle incoming messages
	go handleMessages(ctx, h, sub, t, &mlog)

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

	if cliArgs.WebSocket {
		go server(ctx, t, &mlog, h, cliArgs)
	}

	// If this is in read-only mode, then all you have to do is wait
	if cliArgs.ReadOnly {
		<-ctx.Done() // Wait for the program to close
		return
	}

	// Send messages
	fmt.Println("Welcome to the chat!")
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		m := message{
			Clock: mlog.clock,
			ID:    peer.Encode(h.ID()),
			Name:  cliArgs.Name,
			Text:  s.Text(),
		}
		b, err := json.Marshal(m)
		if err != nil {
			logger.Errorf("Could not marshal message: %v", err)
			continue
		}
		err = t.Publish(ctx, b)
		if err != nil {
			logger.Errorf("Could not publish message: %v", err)
			continue
		}
	}
	if err = s.Err(); err != nil {
		logger.Fatalf("Input scanner error: %v", err)
	}
}

func handleMessages(ctx context.Context, h host.Host, sub *pubsub.Subscription, topic *pubsub.Topic, mlog *messageLog) {
	b, err := json.Marshal(message{
		Type:          JOIN,
		Clock:         mlog.clock,
		ID:            peer.Encode(h.ID()),
		HistoryNumber: 10,
	})
	if err != nil {
		logger.Errorf("Could not serialise join message: %v", err)
	}
	if err = topic.Publish(ctx, b); err != nil {
		logger.Errorf("Could not send join message: %v", err)
	}
	gotHistory := false
	for {
		select {
		case <-ctx.Done():
			return
		default:
			next, err := sub.Next(ctx)
			if err != nil {
				logger.Fatalf("Could not get message: %v", err)
			}
			msg := message{}
			err = json.Unmarshal(next.Data, &msg)
			logger.Infof("Message: %v", msg.Type, msg.ID)
			if msg.Type == JOIN {
				if msg.ID == peer.Encode(h.ID()) {
					continue
				}
				msgs := mlog.getHistory(msg.HistoryNumber)
				b, err := json.Marshal(message{
					Type:          JOIN_REPLY,
					Clock:         msg.Clock,
					ID:            peer.Encode(h.ID()),
					HistoryNumber: len(msgs),
					History:       msgs,
					ReplyTo:       msg.ID,
				})
				if err != nil {
					logger.Errorf("Could not reply to history request: %v", err)
				}
				if err = topic.Publish(ctx, b); err != nil {
					logger.Errorf("Could not send reply to history request: %v", err)
				}
				logger.Infof("%s joined!", msg.ID)
				continue
			} else if msg.Type == JOIN_REPLY {
				if msg.ReplyTo == peer.Encode(h.ID()) && !gotHistory {
					for _, val := range msg.History {
						mlog.Append(val)
					}
					gotHistory = true
				}
				continue
			}
			if err != nil {
				logger.Errorf("Could not decode message: %v", err)
				continue
			}
			mlog.Append(msg)
		}
	}
}

func (mlog *messageLog) Append(msg message) {
	mlog.mux.Lock()
	defer mlog.mux.Unlock()
	sm := sentMessage{
		Clock: msg.Clock,
		ID:    msg.ID,
		Name:  msg.Name,
		Text:  msg.Text,
	}
	if _, ok := mlog.data[sm]; ok {
		// Message already exists
		return
	}
	name := msg.Name
	if name == "" {
		// Use the last 6 characters of the peer's address if no nickname is provided
		name = msg.ID[len(msg.ID)-6 : len(msg.ID)]
	}
	// Note: Shouldn't the message be added to the set here? (It's not in the tutorial)
	if len(mlog.history) == 10 {
		mlog.history = append(mlog.history[1:], msg)
	} else {
		mlog.history = append(mlog.history, msg)
	}
	logger.Infof("%v", mlog.history)
	logger.Infof("%s:\t%s", name, msg.Text)
	if msg.Clock >= mlog.clock {
		mlog.clock = msg.Clock + 1
	}
}

func (mlog *messageLog) getHistory(n int) []message {
	l := len(mlog.history)
	if l < n {
		n = l
	}
	return mlog.history[:n]
}

type messageLog struct {
	mux     sync.Mutex
	data    map[sentMessage]struct{}
	clock   uint
	history []message
}

type message struct {
	Type string
	// Message
	Clock uint
	ID    string // The peer ID
	Name  string
	Text  string
	// Join
	HistoryNumber int
	// Join Reply
	History []message
	ReplyTo string
}

type sentMessage struct {
	Clock uint
	ID    string
	Name  string
	Text  string
}
