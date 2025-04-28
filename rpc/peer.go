package rpc

import (
	"sync"

	"crypto/rsa"
)

// current state the peer is in
type state string

const (
	stateHosting state = "HOSTING"
	stateJoining state = "JOINING"
)

type Peer struct {
	lock  sync.Mutex     // lock to protect shared access to peer state
	peers []peerEnd      // list of other peer endpoints
	me    int            // ID for this party, assigned by host
	host  bool           // true if this peer is the host of the game
	state state          // current state of the game
	sk    rsa.PrivateKey // encryption private key
	sigk  rsa.PrivateKey // signature private key
}
