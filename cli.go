package main

import (
	"flag"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"strings"
)

type CommandLineArguments struct {
	Bootstrap bootstraps
	Name      string
	Relay     bool
	ReadOnly  bool
	Ip        string
	WebSocket bool
	WebsocketPort uint64
}

func ParseCliArgs() CommandLineArguments {
	args := CommandLineArguments{}
	flag.Var(&args.Bootstrap, "bootstrap", "Will connect to this `PEER` to bootstrap the network")
	flag.StringVar(&args.Name, "nickname", "", "This `NAME` will be attached to your messages")
	flag.BoolVar(&args.Relay, "relay", false, "Allows other peers to relay through this peer")
	flag.BoolVar(&args.ReadOnly, "ro", false, "Disable input and just observe the chat")
	flag.StringVar(&args.Ip, "ip", "", "Public `IP` address (for relay peers)")
	flag.BoolVar(&args.WebSocket, "websocket", false, "Enable websocket connections to connect to this peer")
	flag.Uint64Var(&args.WebsocketPort, "wsport", 0, "The port where a websocket connection will be opened")
	flag.Parse()
	return args
}

type bootstraps []*peer.AddrInfo

func (bs *bootstraps) String() string {
	stringList := make([]string, len(*bs))
	for i, addr := range *bs {
		stringList[i] = addr.String()
	}
	return strings.Join(stringList, ",")
}

func (bs *bootstraps) Set(str string) error {
	addr, err := multiaddr.NewMultiaddr(str)
	if err != nil {
		return err
	}
	info, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		return err
	}
	*bs = append(*bs, info)
	return nil
}
