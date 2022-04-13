package model

import (
	ws "github.com/gorilla/websocket"
)

type Player struct {
	Id   string
	Conn *ws.Conn
	//LastSeen    time.Time
	//IsOnline    bool
	ActiveGames []*Game
	Chan        chan []byte
}
