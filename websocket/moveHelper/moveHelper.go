package moveHelper

import (
	"chess-backend/websocketModels/pieces"

	"golang.org/x/exp/slices"
)

type MoveHelper struct {
	Board          [][]pieces.ChessPiece
	Source         *pieces.ChessCoord
	Target         *pieces.ChessCoord
	CastlingRights []pieces.CastleSide
	PlayerColor    pieces.Color
	EnPassant      *pieces.ChessCoord
	Promote        pieces.Promote
}

func (m *MoveHelper) MakeMove(simulate bool) {
	piece := m.Board[m.Source.Row][m.Source.Column]
	targetPiece := m.Board[m.Target.Row][m.Target.Column]

	isMoveEnPassant := false
	isPawnMove := false
	var castleSide *pieces.CastleSide

	switch piece.(type) {
	case *pieces.PawnPiece:
		isPawnMove = true
		if m.EnPassant != nil {
			isMoveEnPassant = *piece.Coord() == *m.EnPassant
		}
	case *pieces.KingPiece:
		for _, cs := range m.CastlingRights {
			if cs.KingCoord() == *m.Target {
				cpycs := cs
				castleSide = &cpycs
			}
		}
	}

	isMoveCastling := castleSide != nil

	switch {
	case isMoveEnPassant:
		m.handleEnPassantMove()
	case isMoveCastling:
		m.handleCastlingMove(*castleSide)
	default:
		m.handleMove(m.Target, m.Source)
	}
	// End of move

	// If it's only for simulating the move we don't need other calculations
	if simulate {
		return
	}

	// Consequences of move
	m.EnPassant = nil

	_, isPieceRook := piece.(*pieces.RookPiece)
	_, isTargetRook := targetPiece.(*pieces.RookPiece)

	switch {
	case isMoveCastling:
		m.castlingConsequences(piece)
	case len(m.CastlingRights) > 0 && isPieceRook && isTargetRook:
		m.rookInteractionConsequences()
	case isPawnMove:
		m.checkPossibleEnPassant()
		m.checkPossiblePromote()
	}

}

func (m *MoveHelper) castlingConsequences(piece pieces.ChessPiece) {
	var willRemove []pieces.CastleSide
	if piece.Color() == pieces.White {
		willRemove = []pieces.CastleSide{"K", "Q"}
	} else {
		willRemove = []pieces.CastleSide{"k", "q"}
	}

	for _, cs := range willRemove {
		index := slices.Index(m.CastlingRights, cs)
		if index != -1 {
			m.CastlingRights = append(m.CastlingRights[:index], m.CastlingRights[index+1:]...)
		}
	}
}

func (m *MoveHelper) rookInteractionConsequences() {
	index := -1
	list := []pieces.ChessCoord{*m.Source, *m.Target}

	for i, cs := range m.CastlingRights {
		if slices.Contains(list, cs.RookCoord()) {
			index = i
		}
	}
	if index != -1 {
		m.CastlingRights = append(m.CastlingRights[:index], m.CastlingRights[index+1:]...)
	}
}

func (m *MoveHelper) checkPossibleEnPassant() {
	rowDelta := m.Source.Row - m.Target.Row
	if rowDelta == 2 || rowDelta == -2 {
		m.EnPassant = (&pieces.ChessCoord{Row: rowDelta / 2}).Add(m.Target)
	}
}

func (m *MoveHelper) checkPossiblePromote() {
	if m.Target.Row == 0 || m.Target.Row == 7 {
		piece := m.Board[m.Target.Row][m.Target.Column]
		var promotedPiece pieces.ChessPiece
		switch m.Promote {
		case pieces.UnknownPromotion:
			m.Promote = pieces.QueenPromotion
			fallthrough
		case pieces.QueenPromotion:
			promotedPiece = &pieces.QueenPiece{PieceAttr: piece.Attr()}
		case pieces.RookPromotion:
			promotedPiece = &pieces.RookPiece{PieceAttr: piece.Attr()}
		case pieces.BishopPromotion:
			promotedPiece = &pieces.BishopPiece{PieceAttr: piece.Attr()}
		case pieces.KnightPromotion:
			promotedPiece = &pieces.KnightPiece{PieceAttr: piece.Attr()}
		}
		m.Board[m.Target.Row][m.Target.Column] = promotedPiece
	}
}

func (m *MoveHelper) handleEnPassantMove() {
	adder := -1
	piece := m.Board[m.Source.Row][m.Source.Column]
	if piece.Color() == pieces.White {
		adder = 1
	}

	pawnLoc := &pieces.ChessCoord{Row: m.Source.Row + int64(adder), Column: m.Source.Column}
	m.handleMove(m.Target, m.Source)
	m.Board[pawnLoc.Row][pawnLoc.Column] = nil
}

func (m *MoveHelper) handleCastlingMove(cs pieces.CastleSide) {
	m.handleMove(m.Target, m.Source)

	rc := cs.RookCoord()
	racc := cs.RookAfterCastlingCoord()
	m.handleMove(&racc, &rc)
}

func (m *MoveHelper) handleMove(target *pieces.ChessCoord, source *pieces.ChessCoord) {
	m.Board[target.Row][target.Column] = m.Board[source.Row][source.Column]
	m.Board[source.Row][source.Column] = nil
	m.Board[target.Row][target.Column].Move(target)
}
