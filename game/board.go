package game

import (
	"errors"

	ff "github.com/fprasx/secrets-and-spies/ff"
)

type Board struct {
	Graph           [][]ff.Num    // adjacency matrix: Graph[i][j] true if city i connected to city j
	Players         []PlayerState // indexed by player ID
	Territories     []int         // city ownership: -1 means uncaptured, else player ID
	Turn            int           // which player's turn it is
	TurnNumber      int           // overall turn number (to handle city disappearance every 2 turns)
	seed            ff.Num
	noCities        int
	cityToBeRemoved int
}

type PlayerState struct {
	City         int  // current city ID
	Energy       int  // actions available this turn
	Intel        int  // amount of intel
	HasStrikeRep bool // bought strike reports
	HasRapidRec  bool // bought rapid recon
	Dead         bool // whether the player is dead
	Revealed     bool // whether the player is revealed
	DeepCover    bool // whether currently under Deep Cover
	nextEnergy   int
}

func NewBoard(noCities, numPlayers int, g [][]int, initialCities []int, p int, seed ff.Num) *Board {
	territories := make([]int, noCities)
	for i := range territories {
		territories[i] = -1 // all cities uncaptured initially except for initialCities
	}

	players := make([]PlayerState, numPlayers)
	for i := range players {
		players[i] = PlayerState{
			City:       initialCities[i],
			Energy:     2,
			Intel:      0,
			Revealed:   false,
			Dead:       false,
			nextEnergy: 2,
		}
		territories[initialCities[i]] = i
	}
	graph := make([][]ff.Num, noCities)
	for i := 0; i < noCities; i++ {
		graph[i] = make([]ff.Num, noCities)
		for j := 0; j < noCities; j++ {
			graph[i][j] = ff.New(int64(g[i][j]))
		}
	}
	return &Board{
		Graph:           graph,
		Players:         players,
		Territories:     territories,
		Turn:            0,
		TurnNumber:      0,
		seed:            seed,
		noCities:        noCities,
		cityToBeRemoved: -1,
	}
}

// Action represents a player action
// ActionType defines types like Move, Stay, Strike, CaptureCity, SpendIntel
// MoveTarget or SpendDetail are used depending on action type

type Action struct {
	Type       ActionType
	TargetCity int // for Move or CaptureCity
}

type ActionType int

// To kill a player, reveal their location and call a strike on their location
const (
	Move        ActionType = iota
	Strike                 //populate targetcity with the city the player is on
	CaptureCity            //populate targetcity with the city the player is on
	SpendIntel
	Locate
	DeepCoverSpend
	SecretMission
	StrikeReports
	RapidRecon
)

// returns 0 if player is dead
func (b *Board) StartTurn() int {
	playerId := b.Turn
	player := &b.Players[playerId]
	player.DeepCover = false
	// Add 2 intel per city owned

	for _, owner := range b.Territories {
		if owner == playerId {
			b.Players[owner].Intel += 2
		}
	}
	if player.Dead {
		return 0
	}
	player.Energy = player.nextEnergy
	player.nextEnergy = 2
	player.DeepCover = false
	return 1
}
func (b *Board) RevealPlayer(playerID int, city int) {
	b.Players[playerID].Revealed = true
	b.Players[playerID].City = city
}
func (b *Board) ExecuteAction(action Action) error {
	playerID := b.Turn
	player := &b.Players[playerID]
	b.Players[playerID].Energy--
	if player.Energy < 0 {
		return errors.New("not enough energy")
	}
	switch action.Type {
	case Move:
		if action.TargetCity >= 0 {
			player.City = action.TargetCity
		}
		player.Revealed = false // unless reveal happens below
	case Strike:
		for id, other := range b.Players {
			if id != playerID && !other.Dead && other.City == action.TargetCity {
				b.Players[id].Dead = true
			}
		}
	case CaptureCity:
		b.Territories[action.TargetCity] = playerID
		player.Revealed = true
	case Locate:
		if player.Intel < 5 {
			return errors.New("not enough intel for locate")
		}
		player.Intel -= 5
	case DeepCoverSpend:
		if player.Intel < 10 {
			return errors.New("not enough intel for deep cover")
		}
		player.Intel -= 10
		player.DeepCover = true

	case SecretMission:
		if player.Intel < 20 {
			return errors.New("not enough intel for secret mission")
		}
		player.Intel -= 20
		player.nextEnergy++
	case StrikeReports:
		if player.Intel < 10 {
			return errors.New("not enough intel for strike reports")
		}
		player.Intel -= 10
		player.HasStrikeRep = true
	case RapidRecon:
		if player.Intel < 20 {
			return errors.New("not enough intel for rapid recon")
		}
		player.Intel -= 20
		player.HasRapidRec = true
	default:
		return errors.New("unknown action")
	}

	return nil
}

// cleanupTurn handles end-of-turn duties: intel generation, turn advancing
func (b *Board) cleanupTurn() {

	// Advance turn

	b.Turn = (b.Turn + 1) % len(b.Players)
	if b.Turn == 0 {
		if b.cityToBeRemoved == -1 {
			//TODO: USE SEED
			b.cityToBeRemoved = 0
		} else {
			b.removeCity(b.cityToBeRemoved)
			b.cityToBeRemoved = -1
		}
	}
	b.TurnNumber++
}
func (b *Board) removeCity(city int) {
	for row := 0; row < b.noCities; row++ {
		b.Graph[row][city] = ff.New(0)
	}
}
