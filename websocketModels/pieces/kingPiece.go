package pieces

type KingPiece struct {
	*PieceAttr
}

func (p *KingPiece) ToString() string {
	if p.PieceColor == White {
		return "K"
	}
	return "k"
}

func (p *KingPiece) CalculateMovableCoords(board [][]ChessPiece) [][]bool {
	ml := make([][]bool, 8)
	for i := range ml {
		ml[i] = make([]bool, 8)
	}

	adders := [8]*ChessCoord{
		{Row: 1, Column: 1},
		{Row: 1, Column: -1},
		{Row: -1, Column: 1},
		{Row: -1, Column: -1},
		{Row: 0, Column: 1},
		{Row: 0, Column: -1},
		{Row: 1, Column: 0},
		{Row: -1, Column: 0},
	}

	for _, adder := range adders {
		cc := p.Coord().Add(adder)
		if cc != nil {
			p.canMoveSetTrue(board, cc, ml, true, true)
		}
	}

	return ml
}

func (p *KingPiece) CalculateCastlableLocations(
	board [][]ChessPiece,
	ml [][]bool,
	castlingRights []CastleSide,
) {
	for _, cs := range castlingRights {
		switch cs {
		case "k":
			if ml[7][5] && board[7][6] == nil {
				ml[7][6] = true
			}
		case "q":
			if ml[7][3] && board[7][2] == nil && board[7][1] == nil {
				ml[7][2] = true
			}
		case "K":
			if ml[0][5] && board[0][6] == nil {
				ml[0][6] = true
			}
		case "Q":
			if ml[0][3] && board[7][2] == nil && board[7][1] == nil {
				ml[0][2] = true
			}
		}
	}
}
