package main

import (
	"github.com/ipfs/go-log"
)

var logger = log.Logger("group-chat")

const topic = "sylo-group-chat-demo"

func main() {
	var _ = log.SetLogLevel("group-chat", "info")
	cliArgs := ParseCliArgs()
	logger.Infof("Nickname: \"%s\"", cliArgs.Name)
}
