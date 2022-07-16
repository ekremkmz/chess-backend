package websocket

import (
	"chess-backend/data"
	"chess-backend/websocketModels"
	"encoding/json"
)

func checkStatusHandler(msg websocketModels.Command, p *websocketModels.Player) {
	var params websocketModels.CheckStatusParams

	if err := json.Unmarshal(msg.Params, &params); err != nil {
		sendErr(p, msg.CommandId, err)
		return
	}

	player, err := data.GetPlayer(params.Nick)

	if err != nil {
		sendErr(p, msg.CommandId, err)
		return
	}
	player.OnlineFriends.Add(p)
	p.OnlineFriends.Add(player)

	data, _ := json.Marshal(player.GetStatus())

	sendSuccess(p, msg.CommandId, data)
}
