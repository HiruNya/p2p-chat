package main

import (
	"github.com/ipfs/go-log"
	"sync"
)

var logger = log.Logger("group-chat")

const topic = "sylo-group-chat-demo"

func main() {
	var _ = log.SetLogLevel("group-chat", "info")
	cliArgs := ParseCliArgs()
	logger.Infof("Nickname: \"%s\"", cliArgs.Name)
	mlog := messageLog{}
	mlog.data = make(map[message]struct{})
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
