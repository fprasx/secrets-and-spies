package service

import (
	"fmt"
	"slices"

	"github.com/fprasx/secrets-and-spies/utils"
)

type LobbyArgs struct {
	Peers []Peer // list of other peers in the game
	Start bool   // true if game is starting
}

type LobbyReply struct{}

func (s *Spies) Lobby(args *LobbyArgs, reply *LobbyReply) error {
	s.Lock()
	defer s.Unlock()

	utils.Assert(!s.IsHost(), "Host should not be receiving Lobby RPC")

	s.peers = slices.Clone(args.Peers)

	if args.Start {
		utils.Assert(
			s.state == stateInit,
			fmt.Sprintf("s.state should be 0, instead got %v", s.state),
		)

		s.state = stateSeed
	}

	return nil
}

func (p *Peer) Lobby(peers []Peer) {
	var args LobbyArgs
	var reply LobbyReply

	args.Start = false
	args.Peers = slices.Clone(peers)

	for {
		err := p.Call("Spies.Lobby", &args, &reply)
		if err == nil {
			break
		}
	}
}

func (p *Peer) Start(peers []Peer) {
	var args LobbyArgs
	var reply LobbyReply

	args.Start = true
	args.Peers = slices.Clone(peers)

	for {
		err := p.Call("Spies.Lobby", &args, &reply)
		if err == nil {
			break
		}
	}
}
