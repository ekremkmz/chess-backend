package data

import (
	"chess-backend/websocketModels"
	"chess-backend/websocketModels/modelErrors"
	"sync"
)

var activeGames = sync.Map{}

func AddGame(g *websocketModels.Game) {
	activeGames.Store(g.Id, g)
}

func GetGame(id string) (*websocketModels.Game, error) {
	if val, ok := activeGames.Load(id); ok {
		return val.(*websocketModels.Game), nil
	}
	return nil, &modelErrors.GameNotFoundError{}
}
