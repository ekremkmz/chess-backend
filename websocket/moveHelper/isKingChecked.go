package moveHelper

import "chess-backend/websocketModels/pieces"

func IsKingChecked(board [][]pieces.ChessPiece, color pieces.Color) bool {
	// Find king
	kingCoord := FindKingCoord(board, color)

	// Check if king is capturable
	for _, row := range board {
		for _, piece := range row {
			if piece != nil && piece.Color() != color {
				if piece.CalculateMovableCoords(board)[kingCoord.Row][kingCoord.Column] {
					return true
				}
			}
		}
	}
	return false
}

func FindKingCoord(board [][]pieces.ChessPiece, color pieces.Color) *pieces.ChessCoord {
	for _, row := range board {
		for _, piece := range row {
			if _, ok := piece.(*pieces.KingPiece); ok && piece.Color() == color {
				return piece.Coord()
			}
		}
	}
	return nil
}
