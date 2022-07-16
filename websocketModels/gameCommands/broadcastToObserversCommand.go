package gameCommands

import "encoding/json"

type BroadcastToObserversCommand struct {
	Id          string
	CommandName string
	Exception   []string
	Message     []byte
}

func (c BroadcastToObserversCommand) Command() string       { return c.CommandName }
func (c BroadcastToObserversCommand) CommandId() string     { return c.Id }
func (c BroadcastToObserversCommand) Data() json.RawMessage { return json.RawMessage(c.Message) }
func (c BroadcastToObserversCommand) Except() []string      { return c.Exception }
