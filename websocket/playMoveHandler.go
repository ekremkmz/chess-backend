package websocket

import (
	"chess-backend/data"
	"chess-backend/websocket/moveHelper"
	"chess-backend/websocketModels"
	"chess-backend/websocketModels/gameCommands"
	"chess-backend/websocketModels/modelErrors"
	"encoding/json"
	"time"
)

func playMoveHandler(msg websocketModels.Command, p *websocketModels.Player) {
	var params websocketModels.PlayMoveParams

	now := time.Now().UTC()

	if err := json.Unmarshal(msg.Params, &params); err != nil {
		sendErr(p, msg.CommandId, err)
		return
	}

	g, err := data.GetGame(params.GameId)

	if err != nil {
		sendErr(p, msg.CommandId, err)
		return
	}

	if err = checkGamePlayable(g, &params, &now, p); err != nil {
		switch err.(type) {
		// This cases doesn't contain locked timer
		case *modelErrors.IllegalTurn:
		case *modelErrors.GameEnderLockTriggeredError:
		default:
			g.CountDownTimer.Unlock(false)
		}
		sendErr(p, msg.CommandId, err)
		return
	}

	// We are killing timer because move is acceptable
	g.CountDownTimer.Unlock(true)

	playMove(g, &params, &now)

	isChecked, hasValid := moveHelper.CheckThereIsAnyValidMove(g)

	isStealMate := !isChecked && !hasValid
	isCheckMate := isChecked && !hasValid

	response := map[string]any{
		"lastplayed": g.LastPlayed.UnixMilli(),
		"move": map[string]any{
			"target": params.Target,
			"source": params.Source,
		},
		"gameId": params.GameId,
		"white":  g.White.ToMap(),
		"black":  g.Black.ToMap(),
	}

	switch {
	case isCheckMate:
		response["special"] = "checkmate"
		g.GameState = websocketModels.Ended

	case isChecked:
		response["special"] = "check"

	case isStealMate:
		response["special"] = "stealmate"
		g.GameState = websocketModels.Ended

	}

	data, _ := json.Marshal(response)

	sendSuccess(p, msg.CommandId, data)

	// We send response to the player first because we will gonna start timer
	// Observers can wait :P
	command := websocketModels.Command{CommandId: msg.CommandId, Command: msg.Command, Params: data}
	otherPlayersState := g.WhoPlays()
	sendCommand(otherPlayersState.Player, command)

	switch g.GameState {
	case websocketModels.Playing:
		// We are giving latency balance to this player because
		// this message gonna travel too
		latency := otherPlayersState.Player.Latency
		g.CountDownTimer.Start(otherPlayersState.TimeLeft+latency, g)
	case websocketModels.Ended:
		//TODO:
	case websocketModels.WaitsFirstMove:
		// Means black should play its first move
		g.CountDownTimer.Start(30*time.Second, g)
	}

	g.CommandChan <- gameCommands.BroadcastToObserversCommand{
		Id:          msg.CommandId,
		CommandName: msg.Command,
		Exception:   []string{p.Nick, otherPlayersState.Player.Nick},
		Message:     data,
	}
}
