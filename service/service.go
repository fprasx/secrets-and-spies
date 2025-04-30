package service

import (
	"log"
	"net"
	"net/rpc"
	"os"
	// "strconv"
	"sync"
	"time"

	"slices"

	"github.com/fprasx/secrets-and-spies/bgw"
	"github.com/fprasx/secrets-and-spies/ff"
	"github.com/fprasx/secrets-and-spies/game"
	"github.com/fprasx/secrets-and-spies/utils"
)

type state int

type Spies struct {
	b              *game.Board
	started        bool       // true if the game has started
	lock           sync.Mutex // lock to protect shared access to peer state
	me             int        // unique id, assigned by host
	next           int        // next unique id to be assigned by host
	peer           Peer       // my network address
	peers          []Peer     // list of other peer endpoints
	OppoAction     game.Action
	actionSent     bool
	receivedNumber int
	Shares         [][][][2]ff.Num
	NewShares      [][][2]ff.Num
	ConfirmCount   int
	EndRound       bool
}

func (s *Spies) Started() bool {
	s.Lock()
	defer s.Unlock()

	return s.started
}

func (s *Spies) Peers() []Peer {
	s.Lock()
	defer s.Unlock()

	return slices.Clone(s.peers)
}

func (s *Spies) Players() int {
	s.Lock()
	defer s.Unlock()

	return len(s.peers)
}

func (s *Spies) IsHost() bool {
	s.Lock()
	defer s.Unlock()

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

// alice := service.New("unix:/tmp/spies/a").WithName("Alice").WithHost(true)
//bob := service.New("unix:/tmp/spies/b").WithName("Bob").WithHost(false)
// bob.Join("unix:/tmp/spies/a")

func New(hostname string) *Spies {
	s := new(Spies)

	addr, err := utils.ResolveAddr(hostname)

	if err != nil {
		log.Fatal(err)
	}

	s.peer = Peer{Addr: addr}
	s.started = false
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
	if s.IsHost() {
		return s
	}

	addr, err := utils.ResolveAddr(hostname)

	if err != nil {
		log.Fatal(err)
	}

	end := Peer{Addr: addr}

	s.Lock()
	peer := s.peer
	s.Unlock()
	log.Printf("me %v \n", s.me)
	me, err := end.Connect(peer)

	if err != nil {
		log.Fatal(err)
	}

	s.Lock()
	s.me = me
	s.Unlock()

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

func (s *Spies) PlayGame() {
	t := 2
	s.b = game.NewBoard(4, 2,
		[][]int{{1, 1, 1, 1},
			{1, 1, 0, 1},
			{1, 0, 1, 1},
			{1, 1, 1, 1}}, []int{0, 1}, t, ff.New(1))
	s.actionSent = false
	s.Shares = make([][][][2]ff.Num, len(s.b.Players))
	s.NewShares = make([][][2]ff.Num, len(s.b.Players))
	s.EndRound = true
	s.Lock()
	for s.started == false {
		s.Unlock()
		time.Sleep(10 * time.Millisecond)
		s.Lock()
	}
	s.Unlock()
	log.Printf("STARTED %v\n", s.me)
	shares, _ := bgw.ShareLocation(s.b.Players[s.me].City, s.b.NoCities, t, len(s.b.Players))
	for i := range s.b.Players {
		column := make([][2]ff.Num, s.b.NoCities)
		for j := 0; j < s.b.NoCities; j++ {
			column[j] = shares[j][i]
		}
		if i != s.me {
			go s.peers[i].SendLocation(shares, 0, s.me)
		} else {
			s.Shares[s.me] = shares
		}
	}
	// for i := range s.b.Players {
	// 	column := make([][2]ff.Num, s.b.NoCities)
	// 	for j := 0; j < s.b.NoCities; j++ {
	// 		column[j] = shares[j][i]
	// 	}
	// 	if i != s.me {
	// 		go s.peers[i].SendLocation(column, 0, s.me)
	// 	} else {
	// 		s.Shares[s.me] = column
	// 	}
	// }
	s.Lock()
	for s.receivedNumber != len(s.b.Players)-1 {
		s.Unlock()
		time.Sleep(10 * time.Millisecond)
		s.Lock()
	}
	s.Unlock()
	// log.Println("all locations received")
	// index := 0
	// var actions []string
	// if s.me == 0 {
	//
	// 	actions = []string{"m1", "m1", "m0", "m0"}
	// } else {
	// 	actions = []string{"m0", "m0", "m3", "m3"}
	// }
	//
	// for {
	//
	// 	s.DoTurn(s.b, s.me, nil, func(spies *Spies) game.Action {
	// 		if index >= len(actions) {
	// 			log.Println("Ran out of actions")
	// 			for {
	// 				time.Sleep(10 * time.Millisecond)
	// 			}
	// 		}
	// 		input := actions[index]
	// 		index++
	// 		if input[0:1] == "m" {
	// 			city, _ := strconv.Atoi(input[1:])
	// 			return game.Action{
	// 				Type:       game.Move,
	// 				TargetCity: city,
	// 			}
	// 		}
	// 		panic("invalid input")
	// 		return game.Action{}
	// 		// return game.Action{
	// 		// 	Type:       game.Move,
	// 		// 	TargetCity: 1,
	// 		// }
	// 	})
	// }
}
func (s *Spies) HostStart() {
	s.Lock()
	if s.started {
		return
	}

	peers := slices.Clone(s.peers)
	utils.Assert(s.me == 0, "Expected to be host")
	s.Unlock()

	s.Broadcast(func(e *Peer) { e.Start(peers) })

	s.Lock()
	s.started = true
	s.Unlock()
}
