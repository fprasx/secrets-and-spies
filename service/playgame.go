package service

import (
	"github.com/fprasx/secrets-and-spies/bgw"
	"github.com/fprasx/secrets-and-spies/ff"
	"github.com/fprasx/secrets-and-spies/game"
)

func (spies *Spies) doTurn(b *game.Board, playerID int, shares [][][2]ff.Num) error {
	b.StartTurn()
	if b.Turn == playerID {
		MyTurn(b, playerID, spies)
	} else {
		OpponentTurn(b, b.Turn, spies, playerID, shares)
	}
	b.CleanupTurn()
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
func RPCReceiveAndValidate(spies *Spies) int {
	return 1
}
func WaitForAction(spies *Spies) game.Action {
	return game.Action{
		Type:       game.Move,
		TargetCity: 1,
	}
}
func WaitForOpponentAction(spies *Spies) game.Action {
	return game.Action{
		Type:       game.Move,
		TargetCity: 1,
	}
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
func WaitForConfirmation(spies *Spies) bool {
	return true
}
func SendConfirmation(spies *Spies) bool {
	return true
}

func OpponentTurn(b *game.Board, playerID int, spies *Spies, meID int, shares [][][2]ff.Num) int {

	player := &b.Players[playerID]
	if player.Dead {
		return 0
	}
	for player.Energy > 0 {
		action := WaitForOpponentAction(spies)
		switch action.Type {
		case game.Move:
			newshare := RPCReceiveVec(spies)
			//		valshares, _ := bgw.ValidateMoveShares(b.Graph, b.NoCities, shares[playerID], newshare, len(b.Players))

			//TODO: Broadcast valshares to everyone

			//This receives all the shares from everyone
			if RPCReceiveAndValidate(spies) == 0 {
				return -1
			}

			//TODO: Send to all parties or just one?
			if !SendConfirmation(spies) {
				return -1
			}

			//	for oppo, j := range b.Players {
			//TODO: Send dot product of oppo*mover to oppo
			//	}
			//this receives just the dot product of my location and oppo location
			if RPCReceiveAndValidate(spies) == 1 {
				b.RevealPlayer(playerID, b.Players[meID].City)

			}
			shares[playerID] = newshare
		case game.Strike:
			dotShares := make([][2]ff.Num, len(b.Players))
			for oppo, _ := range b.Players {
				dotShares[oppo] = bgw.DotProductShares(shares[oppo], shares[playerID], len(b.Players))
				RPCSendNum(dotShares[oppo][1], oppo)
				// if (j.HasStrikeRep){
				// 	RPCSendVec(shares[playerID], oppo)
				// }
			}
			SendDotProducts(dotShares, spies)
			if RPCReceiveDotProductResult(spies) != 0 {
				b.Players[meID].Dead = true

			}
		}

	}
	return 0
}
func MyTurn(b *game.Board, playerID int, spies *Spies) int {
	player := &b.Players[playerID]
	if player.Dead {
		return 0
	}
	for player.Energy > 0 {
		action := WaitForAction(spies)
		switch action.Type {
		case game.Move:
			//TODO: Broadcast Move
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
			err := b.ExecuteAction(action)
			if err != nil {
				return -1
			}

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
