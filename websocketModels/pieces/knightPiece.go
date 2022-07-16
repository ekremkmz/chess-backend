package pieces

type KnightPiece struct {
	*PieceAttr
}

func (p *KnightPiece) ToString() string {
	if p.PieceColor == White {
		return "N"
	}
	return "n"
}

func (p *KnightPiece) CalculateMovableCoords(board [][]ChessPiece) [][]bool {
	ml := make([][]bool, 8)
	for i := range ml {
		ml[i] = make([]bool, 8)
	}

	adders := [8]*ChessCoord{
		{Row: 1, Column: 2},
		{Row: 1, Column: -2},
		{Row: -1, Column: 2},
		{Row: -1, Column: -2},
		{Row: 2, Column: 1},
		{Row: 2, Column: -1},
		{Row: -2, Column: 1},
		{Row: -2, Column: -1},
	}

	for _, adder := range adders {
		cc := p.Coord().Add(adder)
		if cc != nil {
			p.canMoveSetTrue(board, cc, ml, true, true)
		}
	}
	return ml
}
