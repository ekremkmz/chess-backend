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

type Game struct {
	Id             string
	WhoWaits       PlayerState
	WhoPlays       PlayerState
	TimeControl    time.Duration
	AddTimePerMove time.Duration
	LastPlayed     *time.Time
	Started        *time.Time
	GameState      GameState
	Observers      []*Player
	BoardState     BoardState
}

func (g *Game) Play(playerId string, params PlayMoveParams) error {
	state := g.WhoPlays
	if state.Player.Id != playerId {
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
	g.WhoWaits, g.WhoPlays = g.WhoPlays, g.WhoWaits

	// Send response to other players
	data, _ := json.Marshal(params)
	for _, value := range g.Observers {
		if value.Id != playerId {
			value.Chan <- data
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
	}

	whoPlays := g.WhoPlays.ToMap()
	if whoPlays != nil {
		whoPlays, _ := json.Marshal(whoPlays)
		mymap["whoPlays"] = json.RawMessage(whoPlays)
	}

	whoWaits := g.WhoWaits.ToMap()
	if whoWaits != nil {
		whoWaits, _ := json.Marshal(whoWaits)
		mymap["whoPlays"] = json.RawMessage(whoWaits)
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
