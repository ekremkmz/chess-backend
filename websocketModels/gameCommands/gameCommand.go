package gameCommands

import "encoding/json"

type GameCommand interface {
	Command() string
	CommandId() string
	Data() json.RawMessage
	Except() []string
}
