package model

import (
	"time"
)

type Color int64

type PlayerState struct {
	Player   *Player
	TimeLeft time.Duration
}

func (p *PlayerState) ToMap() map[string]any {
	if p.Player == nil {
		return nil
	}

	return map[string]any{
		"nick":     p.Player.Nick,
		"timeleft": int64(p.TimeLeft / time.Millisecond),
	}
}
