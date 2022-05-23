package model

import (
	"strconv"
	"strings"
)

type BoardState struct {
	Board          [8][8]string
	ActiveColor    Color
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

func (b *BoardState) ToMap() map[string]any {
	return map[string]any{
		"board":          BoardToString(b.Board),
		"activeColor":    b.ActiveColor,
		"castlingRights": b.CastlingRights,
		"enPassant":      b.EnPassant,
		"halfMove":       b.HalfMove,
		"fullMove":       b.FullMove,
	}
}

func BoardToString(board [8][8]string) string {
	str := []string{}
	counter := 0
	for _, v := range board {
		substr := ""
		for _, v2 := range v {
			if v2 == "" {
				counter++
			} else {
				if counter != 0 {
					substr += strconv.Itoa(counter)
					counter = 0
				}
				substr += v2
			}
		}
		if counter != 0 {
			substr += strconv.Itoa(counter)
			counter = 0
		}
		str = append(str, substr)
	}
	return strings.Join(str, "/")
}

func StringToBoard(str string) [8][8]string {
	//TODO:fixle
	var board [8][8]string
	rows := strings.Split(str, "/")
	column := 0
	for index1, row := range rows {
		chars := strings.Split(row, "")
		column = 0
		for _, char := range chars {
			num, err := strconv.ParseInt(char, 10, 64)
			if err == nil {
				copy(board[index1][column:], make([]string, num))
				column += int(num)
			} else {
				board[index1][column] = char
				column++
			}
		}
	}
	return board
}
