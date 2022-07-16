package websocketModels

import (
	"chess-backend/syncslice"
	"time"

	ws "github.com/gorilla/websocket"
)

type Player struct {
	Nick          string
	Conn          *ws.Conn
	Latency       time.Duration
	ObserverOf    []*Game
	Chan          chan []byte
	Online        bool
	lastPing      time.Time
	OnlineFriends *syncslice.SyncSlice[*Player]
}

func (p *Player) SetPingTime() {
	p.lastPing = time.Now()
}

func (p *Player) SetLatency() {
	p.Latency = time.Now().Sub(p.lastPing) / 2
}

func (p *Player) GetStatus() map[string]any {
	json := map[string]any{
		"nick": p.Nick,
	}
	if p.Online {
		json["status"] = "online"
	} else {
		json["status"] = "offline"
	}
	return json
}
