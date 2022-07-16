package websocket

import (
	"chess-backend/websocketModels"
	"chess-backend/websocketModels/gameCommands"
	"encoding/json"
)

func connectToGameHandler(msg websocketModels.Command, p *websocketModels.Player) {
	var params websocketModels.ConnectToGameParams

	if err := json.Unmarshal(msg.Params, &params); err != nil {
		sendErr(p, msg.CommandId, err)
		return
	}

	g, ps, err := connectToGame(params, p)

	if err != nil {
		sendErr(p, msg.CommandId, err)
		return
	}

	data, _ := json.Marshal(g.ToMap())

	sendSuccess(p, msg.CommandId, data)

	if ps != nil {
		return
	}

	toObservers := map[string]any{
		"gameId":   g.Id,
		"observer": p.Nick,
	}

	data, _ = json.Marshal(toObservers)

	g.CommandChan <- gameCommands.BroadcastToObserversCommand{
		Id:          msg.CommandId,
		CommandName: msg.Command,
		Exception:   []string{p.Nick},
		Message:     data,
	}
}
