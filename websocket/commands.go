package websocket

import (
	"chess-backend/websocket/model"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func CreateNewGame(c model.CreateNewGameParams, p *model.Player) *model.Game {
	var g model.Game

	g.Id = uuid.New().String()

	if c.Color == "" {
		rand := rand.Intn(2)
		switch rand {
		case 0:
			c.Color = "w"
		case 1:
			c.Color = "b"
		}
	}

	min, _ := strconv.ParseInt(c.TimeControl, 10, 64)
	sec, _ := strconv.ParseInt(c.Adder, 10, 64)

	switch c.Color {
	case "w":
		g.White = &model.PlayerState{Player: p, TimeLeft: time.Duration(min) * time.Minute}
		g.Black = &model.PlayerState{TimeLeft: time.Duration(min) * time.Minute}
	case "b":
		g.White = &model.PlayerState{TimeLeft: time.Duration(min) * time.Minute}
		g.Black = &model.PlayerState{Player: p, TimeLeft: time.Duration(min) * time.Minute}
	}

	g.TimeControl = time.Duration(min) * time.Minute
	g.AddTimePerMove = time.Duration(sec) * time.Second
	g.GameState = model.WaitsOpponent
	g.Observers = append(g.Observers, p)
	g.BoardState = model.BoardState{CastlingRights: []string{"K", "Q", "k", "q"}, EnPassant: "-"}

	AddGame(&g)

	return &g
}

func ConnectToGame(c model.ConnectToGameParams, p *model.Player) (*model.Game, error) {

	g, err := GetGame(c.GameId)

	if err != nil {
		return nil, err
	}

	if g.GameState == model.WaitsOpponent {
		switch {
		case g.White.Player == nil:
			g.White.Player = p
		case g.Black.Player == nil:
			g.Black.Player = p
		default:
		}
		g.GameState = model.WaitsFirstMove
	}

	g.Observers = append(g.Observers, p)

	return g, nil
}
