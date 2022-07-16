package websocketModels

import (
	"chess-backend/websocketModels/pieces"
	"strconv"
	"strings"
)

type BoardState struct {
	Board          [][]pieces.ChessPiece
	ActiveColor    pieces.Color
	CastlingRights []pieces.CastleSide
	EnPassant      *pieces.ChessCoord
	HalfMove       int32
	FullMove       int32
}

func BoardToString(board [][]pieces.ChessPiece) string {
	str := []string{}

	for _, v := range board {
		substr := ""
		counter := 0
		for _, piece := range v {
			if piece == nil {
				counter++
			} else {
				if counter != 0 {
					substr += strconv.Itoa(counter)
					counter = 0
				}
				substr += piece.ToString()
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

func StringToBoard(str string) [][]pieces.ChessPiece {
	board := make([][]pieces.ChessPiece, 8)
	for i := range board {
		board[i] = make([]pieces.ChessPiece, 8)
	}
	rows := strings.Split(str, "/")
	column := 0
	for index1, row := range rows {
		chars := strings.Split(row, "")
		column = 0
		for _, char := range chars {
			num, err := strconv.ParseInt(char, 10, 64)
			if err != nil {
				coord := &pieces.ChessCoord{Row: int64(index1), Column: int64(column)}
				board[index1][column] = charToPiece(char, coord)
				column++
			} else {
				column += int(num)
			}
		}
	}
	return board
}

func charToPiece(char string, coord *pieces.ChessCoord) pieces.ChessPiece {
	isBlack := char != strings.ToUpper(char)
	pieceColor := pieces.White
	if isBlack {
		pieceColor = pieces.Black
	}

	pieceAttr := &pieces.PieceAttr{PieceColor: pieceColor, Coordinates: coord}

	switch strings.ToUpper(char) {
	case "K":
		return &pieces.KingPiece{PieceAttr: pieceAttr}
	case "Q":
		return &pieces.QueenPiece{PieceAttr: pieceAttr}
	case "R":
		return &pieces.RookPiece{PieceAttr: pieceAttr}
	case "B":
		return &pieces.BishopPiece{PieceAttr: pieceAttr}
	case "N":
		return &pieces.KnightPiece{PieceAttr: pieceAttr}
	case "P":
		return &pieces.PawnPiece{PieceAttr: pieceAttr}
	}
	panic("Unknown piece")
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
