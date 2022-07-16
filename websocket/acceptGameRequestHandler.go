package websocket

import (
	"chess-backend/websocketModels"
	"encoding/json"
)

func acceptGameRequestHandler(msg websocketModels.Command, p *websocketModels.Player) {
	var params websocketModels.AcceptGameRequestParams

	if err := json.Unmarshal(msg.Params, &params); err != nil {
		sendErr(p, msg.CommandId, err)
		return
	}

	g, _, err := connectToGame(websocketModels.ConnectToGameParams{GameId: params.GameId}, p)

	if err != nil {
		sendErr(p, msg.CommandId, err)
		return
	}

	data, _ := json.Marshal(g.ToMap())

	var otherPlayer *websocketModels.Player
	switch {
	case g.White.Player.Nick == p.Nick:
		otherPlayer = g.Black.Player
	case g.Black.Player.Nick == p.Nick:
		otherPlayer = g.White.Player
	default:
	}

	sendSuccess(p, msg.CommandId, data)
	sendSuccess(otherPlayer, params.SenderCommandId, data)
}
