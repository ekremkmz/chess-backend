package websocket

import (
	"chess-backend/data"
	"chess-backend/websocketModels"
	"chess-backend/websocketModels/gameCommands"
	"chess-backend/websocketModels/modelErrors"
	"chess-backend/websocketModels/pieces"
	"encoding/json"
	"time"
)

func sendResponseToDrawOffer(msg websocketModels.Command, p *websocketModels.Player) {
	var params websocketModels.SendResponseToDrawOfferParams

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
		"gameId":   g.Id,
		"response": params.Response,
	}

	if params.Response {
		switch {
		case g.DrawOfferedFrom == "":
			sendErr(p, msg.CommandId, &modelErrors.NoDrawOfferError{})
			return
		case g.DrawOfferedFrom == g.White.Player.Nick && g.Black.Player.Nick == p.Nick:
			fallthrough
		case g.DrawOfferedFrom == g.Black.Player.Nick && g.White.Player.Nick == p.Nick:
			g.GameState = websocketModels.Ended
			g.Winner = 2
		default:
			sendErr(p, msg.CommandId, &modelErrors.PrivilegeError{})
			return
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
