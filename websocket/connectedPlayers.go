package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"chess-backend/websocket/model"
	"chess-backend/websocket/model/modelErrors"

	ws "github.com/gorilla/websocket"
)

var ConnectedPlayers = sync.Map{}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

func AddPlayer(p *model.Player) {
	ConnectedPlayers.Store(p.Nick, p)
	go playerReadHandler(p)
	go playerWriteHandler(p)
}

func GetPlayer(nick string) (*model.Player, error) {
	if val, ok := ConnectedPlayers.Load(nick); ok {
		return val.(*model.Player), nil
	}
	return nil, &modelErrors.PlayerNotFoundError{Nick: nick}
}

func DeletePlayer(nick string) {
	ConnectedPlayers.Delete(nick)
}

func sendAck(p *model.Player, commandId string) {
	js := map[string]any{
		"commandId": commandId,
		"command":   "ack",
	}
	msg, _ := json.Marshal(js)
	p.Chan <- msg
}

func sendErr(p *model.Player, commandId string, err error) {
	js := map[string]any{
		"commandId": commandId,
		"command":   "error",
		"data":      err.Error(),
	}
	msg, _ := json.Marshal(js)
	p.Chan <- msg
}

func sendSuccess(p *model.Player, commandId string, data json.RawMessage) {
	js := map[string]any{
		"commandId": commandId,
		"command":   "success",
		"data":      data,
	}
	msg, _ := json.Marshal(js)
	p.Chan <- msg
}

func playerWriteHandler(p *model.Player) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		p.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-p.Chan:
			if !ok {
				p.Conn.WriteMessage(ws.CloseMessage, []byte{})
				return
			}
			if err := p.Conn.WriteMessage(ws.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			p.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.Conn.WriteMessage(ws.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func playerReadHandler(p *model.Player) {
	defer func() {
		close(p.Chan)
		DeletePlayer(p.Nick)
	}()

	p.Conn.SetPongHandler(func(string) error {
		return p.Conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	var msg model.Command
	for {
		// Read messages as Command structs from the connection as long as the connection
		if err := p.Conn.ReadJSON(&msg); err != nil {
			switch err.(type) {
			case *json.SyntaxError:
				log.Printf(err.Error())
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
			var params model.CreateNewGameParams

			if err := json.Unmarshal(msg.Params, &params); err != nil {
				sendErr(p, msg.CommandId, err)
				break
			}

			g := CreateNewGame(params, p)

			data, _ := json.Marshal(g.ToMap())

			sendSuccess(p, msg.CommandId, data)

		case "connectToGame":
			var params model.ConnectToGameParams

			if err := json.Unmarshal([]byte(msg.Params), &params); err != nil {
				sendErr(p, msg.CommandId, err)
				break
			}

			g, err := ConnectToGame(params, p)

			if err != nil {
				sendErr(p, msg.CommandId, err)
				break
			}

			data, _ := json.Marshal(g.ToMap())

			sendSuccess(p, msg.CommandId, data)

		case "playMove":
			var params model.PlayMoveParams

			if err := json.Unmarshal([]byte(msg.Params), &params); err != nil {
				sendErr(p, msg.CommandId, err)
				break
			}

			g, err := GetGame(params.GameId)

			if err != nil {
				sendErr(p, msg.CommandId, err)
				break
			}

			if err = g.CheckPlayable(); err != nil {
				sendErr(p, msg.CommandId, err)
				break
			}

			if err := g.Play(p.Nick, params); err != nil {
				sendErr(p, msg.CommandId, err)
				break
			}

			data, _ := json.Marshal(g.WhoWaits().ToMap())

			sendSuccess(p, msg.CommandId, data)

		default:
		}
	}
}
