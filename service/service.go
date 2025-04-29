package service

import (
	"log"
	"net"
	"net/rpc"
	"os"
	"sync"

	"github.com/fprasx/secrets-and-spies/utils"
	"slices"
)

type state int

const (
	stateLobby state = iota
)

type Spies struct {
	state state
	lock  sync.Mutex // lock to protect shared access to peer state
	me    int        // unique id, assigned by host
	next  int        // next unique id to be assigned by host
	peer  Peer       // my network address
	peers []Peer     // list of other peer endpoints
}

func (s *Spies) Peers() []Peer {
	return slices.Clone(s.peers)
}

func (s *Spies) IsHost() bool {
	return s.me == 0
}

func (s *Spies) Lock() {
	s.lock.Lock()
}

func (s *Spies) Unlock() {
	s.lock.Unlock()
}

func (s *Spies) serve() {
	log.Printf("Starting RPC server %v", s.peer)

	rpc.Register(s)

	if utils.IsSocket(s.peer.Addr.String()) {
		os.Remove(s.peer.Addr.String())
	}

	l, err := net.Listen(s.peer.Addr.Network(), s.peer.Addr.String())

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Fatal(err)
				continue
			}
			go rpc.ServeConn(conn)
		}
	}()
}

func New(hostname string) *Spies {
	s := new(Spies)

	addr, err := utils.ResolveAddr(hostname)

	if err != nil {
		log.Fatal(err)
	}

	s.peer = Peer{Addr: addr}
	s.state = stateLobby
	s.peers = []Peer{}
	s.me = -1
	s.next = -1

	s.serve()

	return s
}

func (s *Spies) WithName(name string) *Spies {
	s.Lock()
	defer s.Unlock()

	s.peer.Name = name

	return s
}

func (s *Spies) WithHost(host bool) *Spies {
	s.Lock()
	defer s.Unlock()

	if host {
		s.me = 0
		s.next = 1
		s.peers = append(s.peers, s.peer)
	}

	return s
}

func (s *Spies) Join(hostname string) *Spies {
	s.Lock()
	defer s.Unlock()

	if s.IsHost() {
		return s
	}

	addr, err := utils.ResolveAddr(hostname)

	if err != nil {
		log.Fatal(err)
	}

	end := Peer{Addr: addr}
	me, err := end.Connect(s.peer)

	if err != nil {
		log.Fatal(err)
	}

	s.me = me

	return s
}

func (s *Spies) Broadcast(thunk func(e *Peer)) {
	var wg sync.WaitGroup

	s.Lock()
	peers := slices.Clone(s.peers)
	me := s.me
	s.Unlock()

	for i := range peers {
		if i == me {
			continue
		}

		wg.Add(1)

		go func() {
			defer wg.Done()
			thunk(&peers[i])
		}()
	}

	wg.Wait()
}
