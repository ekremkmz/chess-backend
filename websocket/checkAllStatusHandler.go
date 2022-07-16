package websocket

import (
	"chess-backend/data"
	"chess-backend/websocketModels"
	"encoding/json"
)

func checkAllStatusHandler(msg websocketModels.Command, p *websocketModels.Player) {
	var params websocketModels.CheckAllStatusParams

	if err := json.Unmarshal(msg.Params, &params); err != nil {
		sendErr(p, msg.CommandId, err)
		return
	}

	result := []map[string]any{}
	for _, nick := range params.Nicks {
		player, err := data.GetPlayer(nick)
		if err != nil {
			result = append(result, map[string]any{
				"nick":   nick,
				"status": "offline",
			})
		} else {
			result = append(result, player.GetStatus())
			data, _ := json.Marshal(map[string]any{
				"nick":   p.Nick,
				"status": "online",
			})
			cmd := websocketModels.Command{CommandId: msg.CommandId, Command: "statusUpdate", Params: data}
			player.OnlineFriends.Add(p)
			p.OnlineFriends.Add(player)
			sendCommand(player, cmd)
		}
	}

	data, _ := json.Marshal(map[string]any{
		"statuses": result,
	})

	sendSuccess(p, msg.CommandId, data)
}
