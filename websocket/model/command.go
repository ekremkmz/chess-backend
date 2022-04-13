package model

import "encoding/json"

type Command struct {
	CommandId string          `json:"commandId"`
	Command   string          `json:"command"`
	Params    json.RawMessage `json:"params"`
}

type CreateNewGameParams struct {
	Color       string `json:"color"`
	TimeControl string `json:"time"`
	Adder       string `json:"add"`
}

type ConnectToGameParams struct {
	GameId string `json:"gameId"`
}

type GetGameParams struct {
	GameId string `json:"gameId"`
}

type PlayMoveParams struct {
	GameId string `json:"gameId"`
	Source string `json:"source"`
	Target string `json:"target"`
}
