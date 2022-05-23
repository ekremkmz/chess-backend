package model

import (
	"chess-backend/websocket/model/modelErrors"
	"encoding/json"
	"time"
)

type GameState int64

const (
	Undefined GameState = iota
	WaitsOpponent
	WaitsFirstMove
	Playing
	Ended
)

const (
	White Color = iota
	Black
)

type Game struct {
	Id             string
	Black          *PlayerState
	White          *PlayerState
	TimeControl    time.Duration
	AddTimePerMove time.Duration
	LastPlayed     *time.Time
	Started        *time.Time
	GameState      GameState
	Observers      []*Player
	BoardState     BoardState
	Turn           Color
	CastleSides    string
	HalfMove       int
	FullMove       int
}

func (g *Game) Play(playerNick string, params PlayMoveParams) error {

	state := g.WhoPlays()

	if state.Player.Nick != playerNick {
		return &modelErrors.IllegalTurn{}
	}

	err := g.BoardState.Play(params)

	if err != nil {
		return err
	}

	// Assign time left
	lastPlayed := time.Now()
	durationUsed := lastPlayed.Sub(*g.LastPlayed)
	state.TimeLeft -= durationUsed
	g.LastPlayed = &lastPlayed

	// Set new turns
	switch g.Turn {
	case White:
		g.Turn = Black
	case Black:
		g.Turn = White
	}

	// Send response to other players
	data, _ := json.Marshal(params)

	var res struct {
		Command string          `json:"command"`
		Data    json.RawMessage `json:"data"`
	}

	res.Command = "playMove"

	res.Data = data

	response, _ := json.Marshal(res)

	for _, obs := range g.Observers {
		if obs.Nick != playerNick {
			obs.Chan <- response
		}
	}

	return nil
}

func (g *Game) ToMap() map[string]any {
	mymap := map[string]any{
		"gameId":      g.Id,
		"timeControl": int64(g.TimeControl / time.Minute),
		"adder":       int64(g.AddTimePerMove / time.Second),
		"state":       g.GameState,
		"boardState":  g.BoardState.ToMap(),
	}

	white := g.White.ToMap()
	if white != nil {
		white, _ := json.Marshal(white)
		mymap["white"] = json.RawMessage(white)
	}

	black := g.Black.ToMap()
	if black != nil {
		black, _ := json.Marshal(black)
		mymap["black"] = json.RawMessage(black)
	}

	if g.LastPlayed != nil {
		mymap["lastPlayed"] = *g.LastPlayed
	}

	if g.Started != nil {
		mymap["started"] = *g.Started
	}

	return mymap
}

func (g *Game) CheckPlayable() error {
	for _, v := range []GameState{WaitsFirstMove, Playing} {
		if v == g.GameState {
			return nil
		}
	}
	return &modelErrors.NotPlayableError{}
}

func (g *Game) WhoPlays() *PlayerState {
	c := g.Turn

	var state *PlayerState

	switch c {
	case White:
		state = g.White
	case Black:
		state = g.Black
	}
	return state
}

func (g *Game) WhoWaits() *PlayerState {
	c := g.Turn

	var state *PlayerState

	switch c {
	case White:
		state = g.Black
	case Black:
		state = g.White
	}
	return state
}
