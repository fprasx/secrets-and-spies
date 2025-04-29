package service

import (
	"fmt"
)

type ConnectArgs struct {
	End ClientEnd
}

type ConnectReply struct {
	You int
}

func (s *Spies) Connect(args *ConnectArgs, reply *ConnectReply) error {
	s.Lock()
	defer s.Unlock()

	if !s.isHost() {
		return fmt.Errorf("Cannot connect to non-host")
	}

	if s.state != stateLobby {
		return fmt.Errorf("Cannot join non-lobby game")
	}

	reply.You = s.next
	s.peers = append(s.peers, args.End)
	s.next++

	fmt.Printf("%v\n", s.peers)

	return nil
}

func (e *ClientEnd) Connect(end ClientEnd) (int, error) {
	var args ConnectArgs
	var reply ConnectReply

	args.End = end
	err := e.Call("Spies.Connect", &args, &reply)

	return reply.You, err
}
