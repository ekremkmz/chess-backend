package data

import (
	"chess-backend/websocketModels"
	"chess-backend/websocketModels/modelErrors"
	"sync"
)

var connectedPlayers = sync.Map{}

func AddPlayer(p *websocketModels.Player) {
	connectedPlayers.Store(p.Nick, p)
}

func GetPlayer(nick string) (*websocketModels.Player, error) {
	if val, ok := connectedPlayers.Load(nick); ok {
		return val.(*websocketModels.Player), nil
	}
	return nil, &modelErrors.PlayerNotFoundError{Nick: nick}
}

func LoadAndDeletePlayer(nick string) (*websocketModels.Player, bool) {
	if val, ok := connectedPlayers.Load(nick); ok {
		connectedPlayers.Delete(nick)
		return val.(*websocketModels.Player), true
	}
	return nil, false
}
