package model

import "strconv"

type BoardState struct {
	Board          [8][8]string
	ActiveColor    string
	CastlingRights []string
	EnPassant      string
	HalfMove       int32
	FullMove       int32
}

type coord struct {
	x int32
	y int32
}

func (b *BoardState) Play(params PlayMoveParams) error {
	board := b.Board
	x, err := strconv.ParseInt(string(params.Source[0]), 10, 64)
	if err != nil {
		return err
	}

	y, err := strconv.ParseInt(string(params.Source[1]), 10, 64)
	if err != nil {
		return err
	}

	source := coord{x: int32(x), y: int32(y)}

	x2, err := strconv.ParseInt(string(params.Target[0]), 10, 64)
	if err != nil {
		return err
	}

	y2, err := strconv.ParseInt(string(params.Target[1]), 10, 64)
	if err != nil {
		return err
	}

	target := coord{x: int32(x2), y: int32(y2)}

	piece := board[source.x][source.y]
	board[target.x][target.y] = piece
	board[source.x][source.y] = ""
	return nil
}
