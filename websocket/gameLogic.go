package websocket

import (
	"chess-backend/data"
	"chess-backend/websocketModels"
	"chess-backend/websocketModels/gameCommands"
	"encoding/json"
	"log"
	"sync"

	"golang.org/x/exp/slices"
)

func newGame(g *websocketModels.Game) {
	data.AddGame(g)
	go gameCommandListener(g)
}

func sendGameToREST(g *websocketModels.Game) {
	// TODO
}

func gameCommandListener(g *websocketModels.Game) {
	wg := sync.WaitGroup{}
	for {
		c := <-g.CommandChan
		wg.Add(1)

		go func() {
			sendCommandToObservers(g, c.Command(), c.CommandId(), c.Data(), c.Except())
			wg.Done()
		}()

		if c, ok := c.(gameCommands.EndGameCommand); ok {
			log.Printf("Game %s ended: %v", g.Id, c)
			wg.Wait()
			sendGameToREST(g)
			return
		}
	}
}

func sendCommandToObservers(g *websocketModels.Game, command string, commandId string, data json.RawMessage, except []string) {
	cmd := websocketModels.Command{CommandId: commandId, Command: command, Params: data}
	for obsItem := range g.Observers.Iter() {
		observer := obsItem.Value
		if !slices.Contains(except, observer.Nick) && observer.Online {
			sendCommand(observer, cmd)
		}
	}
}
