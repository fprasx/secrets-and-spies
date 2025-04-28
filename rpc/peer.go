package rpc

import (
	"net"
	"net/rpc"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

type Peer struct {
	lock  sync.Mutex // lock to protect shared access to peer state
	peers []PeerEnd  // list of other peer endpoints
	host  bool       // true if this peer is the host of the game

	uid  uuid.UUID // ID for this party
	addr net.Addr  // my network address
}

func (p *Peer) end() PeerEnd {
	return PeerEnd{uid: p.uid, addr: p.addr}
}

func (p *Peer) serve() {
	log.Info("Starting RPC server at %v", p.addr)

	err := rpc.Register(p)
	if err != nil {
		log.Fatal("Failed to register RPC server", "err", err)
	}

	l, err := net.Listen(p.addr.Network(), p.addr.String())
	if err != nil {
		log.Fatal("Failed to listen on address", "err", err)
	}

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Error("Network accept error", "err", err)
				continue
			}
			go rpc.ServeConn(conn)
		}
	}()
}

func Make(addr net.Addr, host bool) *Peer {
	fail := func(err error) {
		log.Fatal("Failed to initialize peer", "err", err)
	}

	uid, err := uuid.NewV7()

	if err != nil {
		fail(err)
	}

	pr := new(Peer)

	pr.lock = sync.Mutex{}
	pr.host = host
	pr.uid = uid
	pr.addr = addr

	if host {
		pr.peers = []PeerEnd{pr.end()}
	} else {
		pr.peers = []PeerEnd{}
	}

	// start servicing RPC requests
	pr.serve()

	return pr
}
