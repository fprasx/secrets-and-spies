package rpc

import (
	"net"
	"net/rpc"
	"os"
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

	Replies <-chan any // used to receive RPCs
}

func (p *Peer) end() PeerEnd {
	return PeerEnd{uid: p.uid, addr: p.addr}
}

func (p *Peer) serve() {
	log.Infof("Starting RPC server at %v", p.addr)

	rpc.Register(p)

	log.Infof("%v %v", p.addr.Network(), p.addr.String())
	l, err := net.Listen(p.addr.Network(), p.addr.String())

	if err != nil {
		log.Fatal("Failed to listen on address", "err", err)
	}

	close := func() {
		l.Close()
		switch p.addr.Network() {
		case "unix", "unixgram", "unixpacket":
			os.Remove(p.addr.String())
		}
	}

	go func() {
		defer close()
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

func NewPeer(addr net.Addr, host bool) *Peer {
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
	pr.Replies = make(<-chan any)

	if host {
		pr.peers = []PeerEnd{pr.end()}
	} else {
		pr.peers = []PeerEnd{}
	}

	// start servicing RPC requests
	pr.serve()

	return pr
}
