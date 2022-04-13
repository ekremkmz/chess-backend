package websocket

import (
	"chess-backend/websocket/model"
	"chess-backend/websocket/model/modelErrors"
	"sync"
)

var ActiveGames = sync.Map{}

func AddGame(g *model.Game) {
	ActiveGames.LoadOrStore(g.Id, g)
}

func GetGame(id string) (*model.Game, error) {
	if val, ok := ActiveGames.Load(id); ok {
		return val.(*model.Game), nil
	}
	return nil, &modelErrors.GameNotFoundError{}
}
