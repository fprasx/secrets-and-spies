package rpc

import (
	"net"
	"net/rpc"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

type PeerEnd struct {
	name string         // name of the player
	uid  uuid.UUID      // unique identifier for player
	addr net.Addr       // network address of party (§3.2)
}

// Sends an RPC request to given party, synchronously waiting for reply.
// Returns false if something goes wrong.
func (p *PeerEnd) call(rpcname string, args any, reply any) {
	fail := func(err error) {
		log.Fatal("Failed to call RPC", "rpc", rpcname, "party", p.uid, "args", args, "err", err)
	}

	c, err := rpc.Dial(p.addr.Network(), p.addr.String())

	if err != nil {
		fail(err)
	}

	defer c.Close()

	err = c.Call(rpcname, args, reply)

	if err != nil {
		fail(err)
	}
}
