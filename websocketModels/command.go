package websocketModels

import (
	"chess-backend/websocketModels/pieces"
	"encoding/json"
)

type Command struct {
	CommandId string          `json:"commandId"`
	Command   string          `json:"command"`
	Params    json.RawMessage `json:"params"`
}

type CreateNewGameParams struct {
	Color       string `json:"color"`
	TimeControl int64  `json:"time"`
	Adder       int64  `json:"add"`
	Friend      string `json:"friend"`
}

type ConnectToGameParams struct {
	GameId string `json:"gameId"`
}

type AcceptGameRequestParams struct {
	GameId          string `json:"gameId"`
	SenderCommandId string `json:"senderCommandId"`
}

type PlayMoveParams struct {
	GameId  string         `json:"gameId"`
	Source  string         `json:"source"`
	Target  string         `json:"target"`
	Promote pieces.Promote `json:"promote"`
}

type CheckStatusParams struct {
	Nick string `json:"nick"`
}

type CheckAllStatusParams struct {
	Nicks []string `json:"nicks"`
}

type DrawOfferParams struct {
	GameId string `json:"gameId"`
}

type SendResponseToDrawOfferParams struct {
	GameId   string `json:"gameId"`
	Response bool   `json:"response"`
}

type ResignParams struct {
	GameId string `json:"gameId"`
}
