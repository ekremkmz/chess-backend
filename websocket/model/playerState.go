package model

import (
	"time"
)

type Color int64

const (
	White Color = iota
	Black
)

type PlayerState struct {
	Player   *Player
	Color    Color
	TimeLeft time.Duration
}

func (p *PlayerState) ToMap() map[string]any {
	if p.Player == nil {
		return nil
	}
	return map[string]any{
		"playerNick": p.Player.Nick,
		"color":      p.Color,
		"timeleft":   p.TimeLeft,
	}
}
