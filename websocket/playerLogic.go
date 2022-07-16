package websocket

import (
	"chess-backend/data"
	"chess-backend/websocketModels"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 10 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 7) / 10
)

func playerConnects(p *websocketModels.Player) {
	data.AddPlayer(p)
	go playerReadHandler(p)
	go playerWriteHandler(p)
}

func playerDisconnects(p *websocketModels.Player) {
	log.Printf("%s disconnected.", p.Nick)

	_, loaded := data.LoadAndDeletePlayer(p.Nick)
	if !loaded {
		return
	}
	for _, g := range p.ObserverOf {
		g.Observers.Remove(p)
	}
	data, _ := json.Marshal(map[string]any{
		"nick":   p.Nick,
		"status": "offline",
	})
	cmd := websocketModels.Command{CommandId: uuid.New().String(), Command: "statusUpdate", Params: data}
	for pItem := range p.OnlineFriends.Iter() {
		player := pItem.Value
		player.OnlineFriends.Remove(p)
		sendCommand(player, cmd)
	}
}

func sendAck(p *websocketModels.Player, commandId string) {
	if !p.Online {
		return
	}
	js := map[string]any{
		"commandId": commandId,
		"command":   "ack",
	}
	msg, _ := json.Marshal(js)
	p.Chan <- msg
}

func sendErr(p *websocketModels.Player, commandId string, err error) {
	if !p.Online {
		return
	}
	js := map[string]any{
		"commandId": commandId,
		"command":   "error",
		"params":    err.Error(),
	}
	msg, _ := json.Marshal(js)
	p.Chan <- msg
}

func sendSuccess(p *websocketModels.Player, commandId string, data json.RawMessage) {
	if !p.Online {
		return
	}
	js := map[string]any{
		"commandId": commandId,
		"command":   "success",
		"params":    data,
	}
	msg, _ := json.Marshal(js)
	p.Chan <- msg
}

func sendCommand(p *websocketModels.Player, command websocketModels.Command) {
	if !p.Online {
		return
	}
	msg, _ := json.Marshal(command)
	p.Chan <- msg
}

func playerWriteHandler(p *websocketModels.Player) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		p.Online = false
		ticker.Stop()
		p.Conn.Close()
		if r := recover(); r != nil {
			log.Printf("Recovered in playerReadHandler %v", r)
		}
	}()

	for {
		select {
		case msg, ok := <-p.Chan:
			p.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				p.Conn.WriteMessage(ws.CloseMessage, []byte{})
				return
			}
			if err := p.Conn.WriteMessage(ws.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			p.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			p.SetPingTime()
			if err := p.Conn.WriteMessage(ws.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func playerReadHandler(p *websocketModels.Player) {
	defer func() {
		p.Online = false
		close(p.Chan)
		playerDisconnects(p)
		if r := recover(); r != nil {
			log.Printf("Recovered in playerReadHandler %v", r)
		}
	}()

	p.Conn.SetPongHandler(func(string) error {
		p.SetLatency()
		return nil //p.Conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	p.Online = true
	var msg websocketModels.Command
	for {
		// Read messages as Command structs from the connection as long as the connection is open
		if err := p.Conn.ReadJSON(&msg); err != nil {
			log.Printf("error: %v", err)
			switch err.(type) {
			case *json.SyntaxError:
				continue
			case *ws.CloseError:
				if ws.IsUnexpectedCloseError(err, ws.CloseGoingAway, ws.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
			default:
			}
			break
		}

		// Log messages
		logmsg, _ := json.Marshal(msg)
		log.Printf(string(logmsg))

		sendAck(p, msg.CommandId)

		// Handle command that are sent
		switch msg.Command {
		case "createNewGame":
			createNewGameHandler(msg, p)
		case "connectToGame":
			connectToGameHandler(msg, p)
		case "playMove":
			playMoveHandler(msg, p)
		case "checkStatus":
			checkStatusHandler(msg, p)
		case "checkAllStatus":
			checkAllStatusHandler(msg, p)
		case "acceptGameRequest":
			acceptGameRequestHandler(msg, p)
		case "drawOffer":
			drawOfferHandler(msg, p)
		case "sendResponseToDrawOffer":
			sendResponseToDrawOffer(msg, p)
		case "resign":
			resignHandler(msg, p)
		default:
		}
	}
}
