package moveHelper

import (
	"chess-backend/websocketModels"
	"chess-backend/websocketModels/pieces"
)

func CheckThereIsAnyValidMove(g *websocketModels.Game) (isChecked bool, hasValid bool) {
	strBoard := websocketModels.BoardToString(g.BoardState.Board)
	copyBoard := websocketModels.StringToBoard(strBoard)

	isChecked = IsKingChecked(copyBoard, g.BoardState.ActiveColor)
	kingCoord := FindKingCoord(copyBoard, g.BoardState.ActiveColor)

	if isChecked {
		// If king is checked, check first if it can move
		// There is no need to start from irrelevant pieces
		king := copyBoard[kingCoord.Row][kingCoord.Column]

		ml := king.CalculateMovableCoords(copyBoard)

		for i, row := range ml {
			for j, movable := range row {
				innerCopyBoard := websocketModels.StringToBoard(strBoard)
				if movable {
					mH := MoveHelper{
						Board:       innerCopyBoard,
						Source:      kingCoord,
						Target:      &pieces.ChessCoord{Row: int64(i), Column: int64(j)},
						PlayerColor: king.Color(),
					}
					mH.MakeMove(true)
					notValid := IsKingChecked(innerCopyBoard, g.BoardState.ActiveColor)
					if !notValid {
						return isChecked, true
					}
				}
			}
		}
	}

	// Check all pieces until there is a valid move
	for i, row := range copyBoard {
		for j, piece := range row {
			if piece != nil && piece.Color() == g.BoardState.ActiveColor {
				coord := &pieces.ChessCoord{Row: int64(i), Column: int64(j)}
				// We already checked king stuation
				if *coord != *kingCoord {
					ml := piece.CalculateMovableCoords(copyBoard)
					for k, row := range ml {
						for l, movable := range row {
							innerCopyBoard := websocketModels.StringToBoard(strBoard)
							if movable {
								mH := MoveHelper{
									Board:       innerCopyBoard,
									Source:      coord,
									Target:      &pieces.ChessCoord{Row: int64(k), Column: int64(l)},
									PlayerColor: piece.Color(),
									EnPassant:   g.BoardState.EnPassant,
								}
								mH.MakeMove(true)
								notValid := IsKingChecked(innerCopyBoard, g.BoardState.ActiveColor)
								if !notValid {
									return isChecked, true
								}
							}
						}
					}
				}
			}
		}
	}

	return isChecked, false
}
