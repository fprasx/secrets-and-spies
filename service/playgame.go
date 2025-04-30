package service

import (
	"log"
	"time"

	"github.com/fprasx/secrets-and-spies/bgw"
	"github.com/fprasx/secrets-and-spies/ff"
	"github.com/fprasx/secrets-and-spies/game"
)

type ActionArgs struct {
	A          game.Action
	TurnNumber int
	Player     int
	Shares     [][][2]ff.Num
}

type ActionReply struct {
}

func (s *Spies) RPCSendAction(args *ActionArgs, reply *ActionReply) error {
	log.Printf("asdfafas\n")

	s.Lock()

	for s.b.TurnNumber != args.TurnNumber || !s.EndRound {
		log.Printf("turn: %v %v %v\n", s.b.TurnNumber, args.TurnNumber, s.EndRound)
		s.Unlock()
		time.Sleep(10 * time.Millisecond)
		s.Lock()
	}
	defer s.Unlock()
	log.Printf("hhhh%v %v\n", args.A, s.b.TurnNumber)
	s.actionSent = true
	s.EndRound = false
	s.OppoAction = args.A
	if args.A.Type == game.Move {
		log.Printf("shares%v \n", args.Shares)
		s.NewShares = args.Shares
	}
	return nil
}

func (e *Peer) SendAction(action game.Action, turnNumber int, shares [][][2]ff.Num) error {
	var args ActionArgs
	var reply ActionReply

	args.A = action
	args.TurnNumber = turnNumber
	args.Shares = shares
	for {
		log.Printf("sending action\n")
		err := e.Call("Spies.RPCSendAction", &args, &reply)
		log.Printf("err: %v\n", err)
		if err == nil {
			break
		}
	}

	return nil
}

type LocationShareArgs struct {
	Shares     [][][2]ff.Num
	Player     int
	TurnNumber int
}

type LocationShareReply struct {
}

func (s *Spies) RPCSendLocation(args *LocationShareArgs, reply *LocationShareReply) error {

	s.Lock()
	for !s.started {
		s.Unlock()
		time.Sleep(10 * time.Millisecond)
		s.Lock()
	}

	for s.b.TurnNumber != args.TurnNumber {
		//log.Printf("turn: %v %v\n", s.b.TurnNumber, args.TurnNumber)
		s.Unlock()
		time.Sleep(10 * time.Millisecond)
		s.Lock()
	}
	defer s.Unlock()
	log.Printf("received location %v %v\n", args.Player, s.b.TurnNumber)
	s.Shares[args.Player] = args.Shares
	s.receivedNumber++
	return nil
}

func (e *Peer) SendLocation(shares [][][2]ff.Num, turnNumber int, player int) error {
	var args LocationShareArgs
	var reply LocationShareReply

	args.Shares = shares
	args.Player = player
	args.TurnNumber = turnNumber
	for {
		err := e.Call("Spies.RPCSendLocation", &args, &reply)
		log.Printf("err: %v\n", err)
		if err == nil {
			break
		}
	}

	return nil
}

type ConfirmationArgs struct {
	Player     int
	TurnNumber int
}

type ConfirmationReply struct {
}

func (s *Spies) RPCConfirmation(args *ConfirmationArgs, reply *ConfirmationReply) error {
	s.Lock()
	for s.b.TurnNumber != args.TurnNumber {
		//log.Printf("turn: %v %v\n", s.b.TurnNumber, args.TurnNumber)
		s.Unlock()
		time.Sleep(10 * time.Millisecond)
		s.Lock()
	}
	defer s.Unlock()
	log.Printf("received confirmation %v %v\n", args.Player, s.b.TurnNumber)
	s.ConfirmCount++
	return nil
}

func (e *Peer) SendConfirmation(turnNumber int, player int) error {
	var args LocationShareArgs
	var reply LocationShareReply
	args.Player = player
	args.TurnNumber = turnNumber
	log.Println("send that confirmation")
	for {
		err := e.Call("Spies.RPCConfirmation", &args, &reply)
		log.Printf("err: %v\n", err)
		if err == nil {
			break
		}
	}

	return nil
}

func WaitForConfirmation(spies *Spies) bool {
	for {
		//	log.Println("coming")
		spies.Lock()
		//	log.Println("come")
		if spies.ConfirmCount == len(spies.peers)-1 {
			log.Println("done")
			spies.Unlock()
			return true
		}
		spies.Unlock()
		time.Sleep(10 * time.Millisecond)
	}
	return true
}
func (spies *Spies) DoTurn(WaitForAction func(*Spies) game.Action) error {
	b := spies.b
	playerID := spies.me

	b.StartTurn()
	if b.Turn == playerID {
		spies.MyTurn(b, playerID, WaitForAction)
	} else {
		spies.OpponentTurn(b, b.Turn, playerID)
	}
	log.Println("asdfdafa?")
	spies.Lock()
	log.Println("asdfdafa")
	spies.ConfirmCount = 0
	spies.actionSent = false
	b.CleanupTurn()
	spies.Unlock()
	return nil
}

func RPCSendVec(vec [][2]ff.Num, party int) error {
	return nil

}
func RPCSendNum(num ff.Num, party int) error {
	return nil
}

// receives the result of striking
func RPCReceiveDotProductResult(spies *Spies) int {
	return 0
}
func RPCReceiveVec(spies *Spies) [][2]ff.Num {
	return nil
}

// receives the shares of the dot products of new and old locations from every party and reconstructs secret
func RPCReceiveAndValidate(spies *Spies, player int) bool {
	//endshares := make([][2]ff.Num, len(spies.peers))
	log.Printf("new shares %v cities: %v\n", spies.NewShares, spies.b.NoCities)
	loc1 := 0
	for i := 0; i < spies.b.NoCities; i++ {
		c, _ := bgw.ReconstructSecret(spies.Shares[spies.me][i])
		if c.Eq(ff.New(1)) {
			loc1 = i
			break
		}
	}
	for i := 0; i < spies.b.NoCities; i++ {
		c, _ := bgw.ReconstructSecret(spies.NewShares[i])
		if c.Eq(ff.New(1)) {
			log.Printf("my loc %v his loc %v\n", loc1, i)
			if loc1 == i {
				return true
			}
			return false
		}
	}
	// loc1, _ := bgw.ReconstructSecret(spies.Shares[spies.me][1])
	// log.Printf("my loc loc1?: %v\n", loc1)
	// loc2, _ := bgw.ReconstructSecret(spies.NewShares[1])
	// log.Printf("his loc loc1?: %v\n", loc2)
	// for i := 0; i < len(spies.peers); i++ {
	// 	columnold := make([][2]ff.Num, spies.b.NoCities)
	// 	for j := 0; j < spies.b.NoCities; j++ {
	// 		columnold[j] = spies.Shares[spies.me][j][i]
	// 	}
	// 	columnnew := make([][2]ff.Num, spies.b.NoCities)

	// 	for j := 0; j < spies.b.NoCities; j++ {
	// 		columnnew[j] = spies.NewShares[j][i]
	// 	}
	// 	endshares[i] = bgw.DotProductShares(columnold, columnnew, i)
	// }

	//	res, _ := bgw.ReconstructSecret(endshares)
	//	log.Printf("dot prod?: %v\n", loc2)
	return false
}

// WaitForOpponentAction waits for an opponent's action to be received.
// It locks the spies structure to check if an action has been sent,
// and once received, it unlocks and prints the action.
// The function then enters a loop, simulating the wait for further processing.
// Returns a default game.Action for demonstration purposes.

func WaitForOpponentAction(spies *Spies) game.Action {
	log.Println("Waiting for opponent action")
	spies.Lock()
	for !spies.actionSent {
		spies.Unlock()
		time.Sleep(10 * time.Millisecond)
		spies.Lock()
	}
	spies.Unlock()
	log.Printf("Received action %v\n", spies.OppoAction)

	return spies.OppoAction
}

// sends one dot product share per player
func SendDotProducts(shares [][2]ff.Num, spies *Spies) {
	return
}
func WaitForDotProducts(spies *Spies) []ff.Num {
	return make([]ff.Num, len(spies.Peers()))
}
func GetLocation(target int, spies *Spies) int {
	return 0
}

func (spies *Spies) OpponentTurn(b *game.Board, playerID int, meID int) int {

	player := &b.Players[playerID]
	if player.Dead {
		return 0
	}
	for player.Energy > 0 {
		action := WaitForOpponentAction(spies)
		switch action.Type {
		case game.Move:

			//		valshares, _ := bgw.ValidateMoveShares(b.Graph, b.NoCities, shares[playerID], newshare, len(b.Players))

			//TODO: Broadcast valshares to everyone

			//This receives all the shares from everyone
			// if RPCReceiveAndValidate(spies) == 0 {
			// 	return -1
			// }

			//TODO: Send to all parties or just one?
			spies.peers[playerID].SendConfirmation(spies.b.TurnNumber, meID)

			//	for oppo, j := range b.Players {
			//TODO: Send dot product of oppo*mover to oppo
			//	}
			spies.b.ExecuteAction(action)
			//this receives just the dot product of my location and oppo location

			if RPCReceiveAndValidate(spies, playerID) {
				b.RevealPlayer(playerID, b.Players[meID].City)
				log.Printf("He moved to me!\n")
			} else {
				log.Printf("He not moved to me!")
			}
			spies.Shares[playerID] = spies.NewShares
			spies.Lock()
			spies.EndRound = true
			spies.ConfirmCount = 0
			spies.actionSent = false

			spies.Unlock()
			// for {
			// 	time.Sleep(10 * time.Millisecond)
			// }
			// case game.Strike:
			// 	dotShares := make([][2]ff.Num, len(b.Players))
			// 	for oppo, _ := range b.Players {
			// 		dotShares[oppo] = bgw.DotProductShares(shares[oppo], shares[playerID], len(b.Players))
			// 		RPCSendNum(dotShares[oppo][1], oppo)
			// 		// if (j.HasStrikeRep){
			// 		// 	RPCSendVec(shares[playerID], oppo)
			// 		// }
			// 	}
			// 	SendDotProducts(dotShares, spies)
			// 	if RPCReceiveDotProductResult(spies) != 0 {
			// 		b.Players[meID].Dead = true
			//
			// 	}
		}

	}
	return 0
}

func (spies *Spies) MyTurn(b *game.Board, playerID int, WaitForAction func(s *Spies) game.Action) int {
	player := &b.Players[playerID]
	if player.Dead {
		return 0
	}
	for player.Energy > 0 {

		action := WaitForAction(spies)
		switch action.Type {
		case game.Move:
			log.Printf("move to %d\n", action.TargetCity)
			//TODO: Broadcast Move
			shares, _ := bgw.ShareLocation(action.TargetCity, b.NoCities, b.T, len(b.Players))
			for peer := range spies.Peers() {
				if peer == spies.me {
					spies.Shares[playerID] = shares
					continue
				}

				spies.peers[peer].SendAction(game.Action{
					Type: game.Move,
				}, spies.b.TurnNumber, shares)
			}

			//	shares, _ := bgw.ShareLocation(player.City, b.NoCities, b.T, len(b.Players))
			//TODO: communicate share to each party

			//TODO: Broadcast valshares to everyone
			if !WaitForConfirmation(spies) {
				return -1
			}
			//See if anyone was in ur location
			if player.HasRapidRec {
				res := WaitForDotProducts(spies)
				for oppo, j := range res {
					if j.Eq(ff.New(1)) {
						b.RevealPlayer(oppo, player.City)
					} else if b.Players[oppo].HasStrikeRep {
						//TODO: Tell them your share
					}

				}
			}
			spies.ConfirmCount = 0
			err := b.ExecuteAction(action)
			if err != nil {
				return -1
			}
			log.Printf("Done turn %v\n", b.TurnNumber)

		case game.Strike:
			//TODO: Broadcast strike move
			if !WaitForConfirmation(spies) {
				return -1
			}
			//TODO: Broadcast all dot products
			//ALL PARTIES NEED TO RECEIVE THEIR OWN DOT PRODUCTS
			res := WaitForDotProducts(spies)
			for oppo, j := range res {
				if j.Eq(ff.New(1)) {
					b.RevealPlayer(oppo, player.City)
				} else if b.Players[oppo].HasStrikeRep {
					//TODO: Tell them your share
				}

			}

			err := b.ExecuteAction(action)
			if err != nil {
				return -1
			}
		case game.Locate:
			//TODO: Broadcast locate move
			if !WaitForConfirmation(spies) {
				return -1
			}
			loc := GetLocation(action.Target, spies)
			b.RevealPlayer(action.Target, loc)
		default:
			//TODO: Broadcast move
			if !WaitForConfirmation(spies) {
				return -1
			}
			err := b.ExecuteAction(action)

			if err != nil {
				return -1
			}

		}

	}
	return 0
}
