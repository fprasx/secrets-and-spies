package service

import (
	"slices"

	"github.com/fprasx/secrets-and-spies/utils"
)

type LobbyArgs struct {
	Peers []Peer
}

type LobbyReply struct{}

func (s *Spies) Lobby(args *LobbyArgs, reply *LobbyReply) error {
	s.Lock()
	defer s.Unlock()

	utils.Assert(!s.IsHost(), "Host should not be receiving LobbyUpdate RPC")

	args.Peers = slices.Clone(args.Peers)

	return nil
}

func (p *Peer) Lobby(peers []Peer) {
	var args LobbyArgs
	var reply LobbyReply

	args.Peers = slices.Clone(peers)

	for range 5 {
		err := p.Call("Spies.Lobby", &args, &reply)
		if err == nil {
			break
		}
	}
}
