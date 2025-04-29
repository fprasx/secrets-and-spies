package service

import (
	"log"
	"net"
	"net/rpc"
	"os"
	"sync"

	"github.com/fprasx/secrets-and-spies/utils"
)

type state int

const (
	stateLobby state = iota
)

type Spies struct {
	state state
	lock  sync.Mutex  // lock to protect shared access to peer state
	end   ClientEnd   // my network address
	me    int         // unique id, assigned by host
	next  int         // next unique id to be assigned by host
	peers []ClientEnd // list of other peer endpoints
}

func (s *Spies) isHost() bool {
	return s.me == 0
}

func (s *Spies) Lock() {
	s.lock.Lock()
}

func (s *Spies) Unlock() {
	s.lock.Unlock()
}

func (s *Spies) serve() {
	log.Printf("Starting RPC server %v", s.end)

	rpc.Register(s)

	if utils.IsSocket(s.end.Addr.String()) {
		os.Remove(s.end.Addr.String())
	}

	l, err := net.Listen(s.end.Addr.Network(), s.end.Addr.String())

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

	s.end = ClientEnd{Addr: addr}
	s.state = stateLobby
	s.peers = []ClientEnd{}
	s.me = -1
	s.next = -1

	s.serve()

	return s
}

func (s *Spies) WithName(name string) *Spies {
	s.Lock()
	defer s.Unlock()

	s.end.Name = name

	return s
}

func (s *Spies) WithHost(host bool) *Spies {
	s.Lock()
	defer s.Unlock()

	if host {
		s.me = 0
		s.next = 1
		s.peers = append(s.peers, s.end)
	}

	return s
}

func (s *Spies) Join(hostname string) *Spies {
	s.Lock()
	defer s.Unlock()

	if s.isHost() {
		return s
	}

	addr, err := utils.ResolveAddr(hostname)

	if err != nil {
		log.Fatal(err)
	}

	end := ClientEnd{Addr: addr}
	me, err := end.Connect(end)

	if err != nil {
		log.Fatal(err)
	}

	s.me = me

	return s
}
