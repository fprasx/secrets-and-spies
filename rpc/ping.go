package rpc

import (
	"fmt"
)

type PingRequest struct {
	from string // who is the ping from
}

type PingReply struct {
	pong string
}

func (p *Peer) Ping(args *PingRequest, reply *PingReply) error {
	reply.pong = fmt.Sprintf("Pong %s", args.from)
	return nil
}

func (pe *PeerEnd) Ping(from string) string {
	args := PingRequest{from}
	reply := PingReply{}

	pe.call("Peer.Ping", &args, &reply)
	return reply.pong
}
