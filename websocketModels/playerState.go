package websocketModels

import (
	"time"
)

type PlayerState struct {
	Player   *Player
	TimeLeft time.Duration
}

func (p *PlayerState) ToMap() map[string]any {
	return map[string]any{
		"nick":     p.Player.Nick,
		"timeleft": int64(p.TimeLeft / time.Millisecond),
	}
}
