package pieces

type PawnPiece struct {
	*PieceAttr
}

func (p *PawnPiece) ToString() string {
	if p.PieceColor == White {
		return "P"
	}
	return "p"
}

func (p *PawnPiece) CalculateMovableCoords(board [][]ChessPiece) [][]bool {
	ml := make([][]bool, 8)
	for i := range ml {
		ml[i] = make([]bool, 8)
	}
	adder := p.adder()

	// Double
	if p.isFirstMove() {
		cc := &ChessCoord{Row: adder * 2, Column: 0}
		p.canMoveSetTrue(board, p.Coord().Add(cc), ml, true, false)
	}

	// Regular
	cc := &ChessCoord{Row: adder, Column: 0}
	cc = p.Coord().Add(cc)
	if cc != nil {
		p.canMoveSetTrue(board, cc, ml, true, false)
	}

	// Capture
	cc = &ChessCoord{Row: adder, Column: 1}
	cc = p.Coord().Add(cc)

	if cc != nil {
		p.canMoveSetTrue(board, cc, ml, false, true)
	}

	cc = &ChessCoord{Row: adder, Column: -1}
	cc = p.Coord().Add(cc)

	if cc != nil {
		p.canMoveSetTrue(board, cc, ml, false, true)
	}

	return ml
}

func (p *PawnPiece) isFirstMove() bool {
	if p.PieceColor == White {
		return p.Coord().Row == 6
	}
	return p.Coord().Row == 1
}

func (p *PawnPiece) adder() int64 {
	if p.PieceColor == White {
		return -1
	}
	return 1
}
func (p *PawnPiece) CalculateEnPassantLocations(
	ml [][]bool,
	enPassant *ChessCoord,
) {
	if enPassant == nil {
		return
	}

	cc := p.Coord().Add(&ChessCoord{Row: p.adder(), Column: 1})
	if cc != nil && *cc == *enPassant {
		ml[cc.Row][cc.Column] = true
	}

	cc = p.Coord().Add(&ChessCoord{Row: p.adder(), Column: -1})
	if cc != nil && *cc == *enPassant {
		ml[cc.Row][cc.Column] = true
	}
}
