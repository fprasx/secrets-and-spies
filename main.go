package main

import (
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"

	"github.com/fprasx/secrets-and-spies/rpc"

	"github.com/charmbracelet/log"
)

func resolve(address string) (net.Addr, error) {
	u, err := url.Parse(address)
	if err != nil {
		return nil, fmt.Errorf("Invalid address format: %w", err)
	}

	network := u.Scheme
	var addrStr string

	switch network {
	case "tcp", "tcp4", "tcp6", "udp", "udp4", "udp6":
		addrStr = u.Host
	case "ip", "ip4", "ip6":
		addrStr = u.Hostname()
	case "unix", "unixgram", "unixpacket":
		addrStr = u.Path
	default:
		return nil, fmt.Errorf("Unsupported network: %s", network)
	}

	switch network {
	case "tcp", "tcp4", "tcp6":
		return net.ResolveTCPAddr(network, addrStr)
	case "udp", "udp4", "udp6":
		return net.ResolveUDPAddr(network, addrStr)
	case "ip", "ip4", "ip6":
		return net.ResolveIPAddr(network, addrStr)
	case "unix", "unixgram", "unixpacket":
		return net.ResolveUnixAddr(network, addrStr)
	default:
		return nil, fmt.Errorf("Unsupported network: %s", network)
	}
}

func main() {
	log.SetReportCaller(false)
	log.SetReportTimestamp(false)

	hostFlag := flag.Bool("host", false, "Run as host of game instead of player")
	addrFlag := flag.String("addr", "", "Network address with prefix, e.g. unix:///tmp/socket/")

	flag.Parse()

	if *addrFlag == "" {
		log.Errorf("-addr flag is required.\n")
		flag.Usage()
		os.Exit(1)
	}

	host := *hostFlag
	addr, err := resolve(*addrFlag)

	if err != nil {
		log.Errorf("-addr flag is invalid: %v.\n", addr)
		flag.Usage()
		os.Exit(1)
	}

	peer := rpc.NewPeer(addr, host)

	for {
		select {
		case reply := <-peer.Replies:
			log.Info("Recieved rpc reply, ", reply)
		}
	}
}
