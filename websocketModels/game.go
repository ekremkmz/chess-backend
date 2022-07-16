package websocketModels

import (
	"chess-backend/cancellableTimer"
	"chess-backend/syncslice"
	"chess-backend/websocketModels/gameCommands"
	"chess-backend/websocketModels/pieces"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type GameState int64

const (
	Undefined GameState = iota
	WaitsOpponent
	WaitsFirstMove
	Playing
	Ended
)

type Game struct {
	Id              string
	Black           *PlayerState
	White           *PlayerState
	TimeControl     time.Duration
	AddTimePerMove  time.Duration
	LastPlayed      *time.Time
	CreatedAt       *time.Time
	GameState       GameState
	Observers       *syncslice.SyncSlice[*Player]
	BoardState      BoardState
	CountDownTimer  *cancellableTimer.CancellableTimer
	CommandChan     chan gameCommands.GameCommand
	DrawOfferedFrom string
	Winner          int64 //0:White,1:Black,2:Draw
}

func (g *Game) WhenTimeout() {
	whoPlays := g.WhoPlays()
	result := "w"
	g.Winner = 0
	if g.White.Player.Nick == whoPlays.Player.Nick {
		result = "b"
		g.Winner = 1
	}
	reason := "timeout"

	if g.GameState == WaitsFirstMove {
		reason = "firstmove"
	}
	g.CommandChan <- gameCommands.EndGameCommand{Id: uuid.New().String(), GameId: g.Id, Reason: reason, Result: result}
}

func (g *Game) WhoPlays() *PlayerState {
	c := g.BoardState.ActiveColor

	var state *PlayerState

	switch c {
	case pieces.White:
		state = g.White
	case pieces.Black:
		state = g.Black
	}
	return state
}

func (g *Game) WhoWaits() *PlayerState {
	c := g.BoardState.ActiveColor

	var state *PlayerState

	switch c {
	case pieces.White:
		state = g.Black
	case pieces.Black:
		state = g.White
	}
	return state
}

func (g *Game) ToMap() map[string]any {
	mymap := map[string]any{
		"gameId":      g.Id,
		"timeControl": int64(g.TimeControl / time.Minute),
		"adder":       int64(g.AddTimePerMove / time.Second),
		"state":       g.GameState,
		"boardState":  g.BoardState.ToMap(),
	}

	white := g.White
	if white.Player != nil {
		white, _ := json.Marshal(white.ToMap())
		mymap["white"] = json.RawMessage(white)
	}

	black := g.Black
	if black.Player != nil {
		black, _ := json.Marshal(black.ToMap())
		mymap["black"] = json.RawMessage(black)
	}

	if g.LastPlayed != nil {
		mymap["lastPlayed"] = g.LastPlayed.UnixMilli()
	}

	if g.CreatedAt != nil {
		mymap["createdAt"] = g.CreatedAt.UnixMilli()
	}

	return mymap
}
