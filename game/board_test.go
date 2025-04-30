// board_test.go

package game

import (
	"testing"

	ff "github.com/fprasx/secrets-and-spies/ff"
)

func TestBoard_ExecuteAction_Move(t *testing.T) {
	// Create a new board with 2 cities and 1 player
	board := &Board{
		Graph: [][]ff.Num{
			{ff.New(1), ff.New(1)},
			{ff.New(1), ff.New(1)},
		},
		Players: []PlayerState{
			{City: 0, Energy: 2, Intel: 0, Revealed: false, Dead: false, NextEnergy: 2},
		},
		Territories:     []int{-1, -1},
		Turn:            0,
		TurnNumber:      0,
		seed:            ff.New(1),
		NoCities:        2,
		CityToBeRemoved: -1,
		T:               2,
	}

	// Create a move action
	action := Action{
		Type:       Move,
		TargetCity: 1,
	}

	// Execute the action
	err := board.ExecuteAction(action)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Check if the player moved to the new city
	if board.Players[0].City != 1 {
		t.Errorf("expected player to move to city 1, but got city %d", board.Players[0].City)
	}
}

func TestBoard_ExecuteAction_Strike(t *testing.T) {
	// Create a new board with 2 cities and 2 players
	board := &Board{
		Graph: [][]ff.Num{
			{ff.New(1), ff.New(1)},
			{ff.New(1), ff.New(1)},
		},
		Players: []PlayerState{
			{City: 0, Energy: 2, Intel: 0, Revealed: false, Dead: false, NextEnergy: 2},
			{City: 0, Energy: 2, Intel: 0, Revealed: false, Dead: false, NextEnergy: 2},
		},
		Territories:     []int{-1, -1},
		Turn:            0,
		TurnNumber:      0,
		seed:            ff.New(1),
		NoCities:        2,
		CityToBeRemoved: -1,
		T:               2,
	}

	// Create a strike action
	action := Action{
		Type:       Strike,
		TargetCity: 0,
	}

	// Execute the action
	err := board.ExecuteAction(action)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Check if the other player is dead
	if !board.Players[1].Dead {
		t.Errorf("expected player 1 to be dead, but got %v", board.Players[1].Dead)
	}
}
func TestBoard_ExecuteAction_SecretMission(t *testing.T) {
	// Create a new board with 2 cities and 1 player
	board := &Board{
		Graph: [][]ff.Num{
			{ff.New(1), ff.New(1)},
			{ff.New(1), ff.New(1)},
		},
		Players: []PlayerState{
			{City: 0, Energy: 2, Intel: 18, Revealed: false, Dead: false, NextEnergy: 2},
			{City: 1, Energy: 2, Intel: 20, Revealed: false, Dead: false, NextEnergy: 2},
		},
		Territories:     []int{0, 1},
		Turn:            0,
		TurnNumber:      0,
		seed:            ff.New(1),
		NoCities:        2,
		CityToBeRemoved: -1,
		T:               2,
	}
	// Create a secret mission action
	action := Action{
		Type: SecretMission,
	}
	ecode := board.StartTurn()
	if ecode == 0 {
		t.Errorf("expected not dead, got %v", ecode)
	}
	// Execute the action
	err := board.ExecuteAction(action)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Check if the player's intel was spent correctly
	if board.Players[0].Intel != 0 {
		t.Errorf("expected player's intel to be 0, but got %d", board.Players[0].Intel)
	}

	// Check if the player's next energy was increased
	if board.Players[0].NextEnergy != 3 {
		t.Errorf("expected player's next energy to be 3, but got %d", board.Players[0].NextEnergy)
	}

	// Create a move action
	action = Action{
		Type:       Move,
		TargetCity: 1,
	}

	// Execute the action
	err = board.ExecuteAction(action)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Check if the player moved to the new city
	if board.Players[0].City != 1 {
		t.Errorf("expected player to move to city 1, but got city %d", board.Players[0].City)
	}
	board.CleanupTurn()
	board.StartTurn()
	board.CleanupTurn()
	board.StartTurn()
	// Check if the player's energy was increased
	if board.Players[0].Energy != 3 {
		t.Errorf("expected player's energy to be 3, but got %d", board.Players[0].Energy)
	}
}
