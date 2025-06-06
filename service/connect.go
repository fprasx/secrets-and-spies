package service

import (
	"fmt"
	"log"
	"slices"
)

type ConnectArgs struct {
	Peer Peer
}

type ConnectReply struct {
	You int
}

func (s *Spies) ConnectRPC(args *ConnectArgs, reply *ConnectReply) error {
	s.Lock()
	defer s.Unlock()

	if s.me != 0 {
		return fmt.Errorf("Cannot connect to non-host")
	}

	if s.started {
		return fmt.Errorf("Cannot join non-lobby game")
	}

	reply.You = s.next
	s.peers = append(s.peers, args.Peer)
	s.next++

	peers := slices.Clone(s.peers)
	go s.Broadcast(func(p *Peer) { p.Lobby(peers) })

	return nil
}

func (e *Peer) Connect(peer Peer) (int, error) {
	var args ConnectArgs
	var reply ConnectReply

	args.Peer = peer

	for {
		err := e.Call("Spies.ConnectRPC", &args, &reply)
		log.Printf("err: %v\n", err)
		if err == nil {
			break
		}
	}

	return reply.You, nil
}
