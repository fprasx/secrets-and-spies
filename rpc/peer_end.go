package rpc

import (
	"crypto/rsa"
	"net"
	"net/rpc"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

type peerEnd struct {
	name string        // name of the player
	uid  uuid.UUID     // unique identifier for player
	addr net.Addr      // network address of party (§3.2)
	pk   rsa.PublicKey // encryption public key (§3.2)
	vk   rsa.PublicKey // signature verification key (§3.2)
}

// Sends an RPC request to given party, synchronously waiting for reply.
// Returns false if something goes wrong.
func (p *peerEnd) call(rpcname string, args any, reply any) {
	_, err := rpc.Dial(p.addr.Network(), p.addr.String())

	if err != nil {
		log.Fatalf("Failed to call RPC", "rpc", rpcname, "party", p.uid)
	}
}
