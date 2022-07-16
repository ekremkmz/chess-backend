package websocket

import (
	"chess-backend/data"
	"chess-backend/websocketModels"
	"encoding/json"
)

func createNewGameHandler(msg websocketModels.Command, p *websocketModels.Player) {
	var params websocketModels.CreateNewGameParams

	if err := json.Unmarshal(msg.Params, &params); err != nil {
		sendErr(p, msg.CommandId, err)
		return
	}

	g := createNewGame(params, p)

	if params.Friend != "" {
		friend, err := data.GetPlayer(params.Friend)
		if err != nil {
			sendErr(p, msg.CommandId, err)
			return
		}

		data, _ := json.Marshal(map[string]any{
			"gameId":          g.Id,
			"nick":            p.Nick,
			"color":           params.Color,
			"time":            params.TimeControl,
			"add":             params.Adder,
			"senderCommandId": msg.CommandId,
		})
		cmd := websocketModels.Command{CommandId: msg.CommandId, Command: "gameRequest", Params: data}
		sendCommand(friend, cmd)
		return
	}

	data, _ := json.Marshal(g.ToMap())

	sendSuccess(p, msg.CommandId, data)
}
