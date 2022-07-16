package websocket

import (
	"chess-backend/data"
	"chess-backend/websocketModels"
	"chess-backend/websocketModels/gameCommands"
	"chess-backend/websocketModels/pieces"
	"encoding/json"
	"time"
)

func resignHandler(msg websocketModels.Command, p *websocketModels.Player) {
	var params websocketModels.ResignParams

	if err := json.Unmarshal(msg.Params, &params); err != nil {
		sendErr(p, msg.CommandId, err)
		return
	}

	g, err := data.GetGame(params.GameId)

	if err != nil {
		sendErr(p, msg.CommandId, err)
		return
	}

	js := map[string]any{
		"gameId": g.Id,
		"who":    p.Nick,
	}

	g.CountDownTimer.Cancel()

	now := time.Now().UTC()
	delta := now.Sub(*g.LastPlayed)

	switch g.BoardState.ActiveColor {
	case pieces.White:
		g.White.TimeLeft = g.White.TimeLeft - time.Duration(delta)
		js["playerState"] = g.White.ToMap()
	case pieces.Black:
		g.Black.TimeLeft = g.Black.TimeLeft - time.Duration(delta)
		js["playerState"] = g.Black.ToMap()
	}
	g.GameState = websocketModels.Ended

	switch {
	case g.White.Player.Nick == p.Nick:
		g.Winner = 1
	case g.Black.Player.Nick == p.Nick:
		g.Winner = 0
	}

	data, _ := json.Marshal(js)

	sendSuccess(p, msg.CommandId, data)
	g.CommandChan <- gameCommands.BroadcastToObserversCommand{
		Id:          msg.CommandId,
		CommandName: msg.Command,
		Exception:   []string{p.Nick},
		Message:     data,
	}
}
