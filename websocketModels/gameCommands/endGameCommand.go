package gameCommands

import "encoding/json"

type EndGameCommand struct {
	Id     string
	GameId string
	Reason string
	Result string
}

func (c EndGameCommand) Command() string   { return "endGame" }
func (c EndGameCommand) CommandId() string { return c.Id }
func (c EndGameCommand) Data() json.RawMessage {
	js := map[string]any{
		"gameId": c.GameId,
		"reason": c.Reason,
		"result": c.Result,
	}
	data, _ := json.Marshal(js)
	return json.RawMessage(data)
}
func (c EndGameCommand) Except() []string { return []string{} }
