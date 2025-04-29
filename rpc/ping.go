package rpc

type PingArgs struct {
}

type PingReply struct {
	pong string
}

func (p *Peer) Ping(args *PingArgs, reply *PingReply) error {
	reply.pong = "Pong"
	return nil
}

func (pe *PeerEnd) Ping(from string) string {
	args := PingArgs{}
	reply := PingReply{}

	pe.call("Peer.Ping", &args, &reply)
	return reply.pong
}
