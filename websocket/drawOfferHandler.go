package websocket

import (
	"chess-backend/data"
	"chess-backend/websocketModels"
	"chess-backend/websocketModels/gameCommands"
	"chess-backend/websocketModels/modelErrors"
	"encoding/json"
)

func drawOfferHandler(msg websocketModels.Command, p *websocketModels.Player) {
	var params websocketModels.DrawOfferParams

	if err := json.Unmarshal(msg.Params, &params); err != nil {
		sendErr(p, msg.CommandId, err)
		return
	}

	g, err := data.GetGame(params.GameId)

	if err != nil {
		sendErr(p, msg.CommandId, err)
		return
	}

	if g.GameState != websocketModels.Playing {
		sendErr(p, msg.CommandId, &modelErrors.NotPlayableError{})
		return
	}

	g.DrawOfferedFrom = p.Nick

	data, _ := json.Marshal(map[string]any{
		"gameId": g.Id,
		"from":   p.Nick,
	})

	sendSuccess(p, msg.CommandId, data)
	g.CommandChan <- gameCommands.BroadcastToObserversCommand{
		Id:          msg.CommandId,
		CommandName: msg.Command,
		Exception:   []string{p.Nick},
		Message:     data,
	}
}
