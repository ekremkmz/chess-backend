package pieces

type RookPiece struct {
	*PieceAttr
}

func (p *RookPiece) ToString() string {
	if p.PieceColor == White {
		return "R"
	}
	return "r"
}

func (p *RookPiece) CalculateMovableCoords(board [][]ChessPiece) [][]bool {
	ml := make([][]bool, 8)
	for i := range ml {
		ml[i] = make([]bool, 8)
	}
	list := []*ChessCoord{}

	adders := [4]*ChessCoord{
		{Row: 1, Column: 0},
		{Row: -1, Column: 0},
		{Row: 0, Column: 1},
		{Row: 0, Column: -1},
	}

	for _, adder := range adders {
		lastLoc := p.Coord()
		run := true
		for run {
			lastLoc = lastLoc.Add(adder)
			if lastLoc != nil {
				if board[lastLoc.Row][lastLoc.Column] != nil {
					run = false
				}
				list = append(list, lastLoc)
			} else {
				run = false
			}
		}
	}

	for _, loc := range list {
		p.canMoveSetTrue(board, loc, ml, true, true)
	}
	return ml
}
