package utils

import (
	"encoding/gob"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"net/url"
	"os"
)

func Assert(cond bool, msg string) {
	if !cond {
		log.Fatalf("Assertion failed: %v", msg)
	}
}

func ResolveAddr(address string) (net.Addr, error) {
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

func ValidateAddr(address string) error {
	if _, err := ResolveAddr(address); err != nil {
		return err
	}

	return nil
}

func IsSocket(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.Mode().Type() == fs.ModeSocket
}

func RegisterRpcTypes() {
	gob.Register(&net.UnixAddr{})
	gob.Register(&net.TCPAddr{})
	gob.Register(&net.UDPAddr{})
	gob.Register(&net.IPAddr{})
}

func RegisterLogger() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(io.Discard)
}

func AddrString(addr net.Addr) string {
	return fmt.Sprintf("%s://%s", addr.Network(), addr.String())
}
