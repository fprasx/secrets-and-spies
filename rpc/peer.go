package rpc

import (
	"crypto/rand"
	"crypto/rsa"
	"net"
	"net/rpc"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

// current state the peer is in
type state string

const (
	stateHosting state = "HOSTING"
	stateJoining state = "JOINING"
)

type Peer struct {
	lock  sync.Mutex // lock to protect shared access to peer state
	peers []PeerEnd  // list of other peer endpoints
	host  bool       // true if this peer is the host of the game
	state state      // current state of the game

	name string          // name of this player
	uid  uuid.UUID       // ID for this party
	addr net.Addr        // my network address
	pk   rsa.PublicKey   // encryption public key
	vk   rsa.PublicKey   // signature verification key
	sk   *rsa.PrivateKey // encryption private key
	sigk *rsa.PrivateKey // signature private key
}

func (p *Peer) end() PeerEnd {
	return PeerEnd{name: p.name, uid: p.uid, addr: p.addr, pk: p.pk, vk: p.vk}
}

func (p *Peer) serve() {
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

func Make(name string, addr net.Addr, host bool) *Peer {
	fail := func(err error) {
		log.Fatal("Failed to initialize peer", "err", err)
	}

	uid, err := uuid.NewV7()

	if err != nil {
		fail(err)
	}

	sk, err := rsa.GenerateKey(rand.Reader, 1024)

	if err != nil {
		fail(err)
	}

	sigk, err := rsa.GenerateKey(rand.Reader, 1024)

	if err != nil {
		fail(err)
	}

	pr := new(Peer)

	pr.lock = sync.Mutex{}
	pr.host = host
	pr.state = stateHosting

	pr.name = name
	pr.uid = uid
	pr.addr = addr
	pr.sk = sk
	pr.sigk = sigk
	pr.pk = sk.PublicKey
	pr.vk = sigk.PublicKey

	if host {
		pr.peers = []PeerEnd{pr.end()}
	} else {
		pr.peers = []PeerEnd{}
	}

	// start servicing RPC requests
	pr.serve()

	return pr
}
