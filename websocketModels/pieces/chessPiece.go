package pieces

type Color int

const (
	White Color = iota
	Black
)

type ChessPiece interface {
	ToString() string
	Color() Color
	Coord() *ChessCoord
	Move(*ChessCoord)
	Attr() *PieceAttr
	CalculateMovableCoords(board [][]ChessPiece) [][]bool
}

type PieceAttr struct {
	PieceColor  Color
	Coordinates *ChessCoord
}

func (p *PieceAttr) Color() Color {
	return p.PieceColor
}

func (p *PieceAttr) Coord() *ChessCoord {
	return p.Coordinates
}

func (p *PieceAttr) Move(cc *ChessCoord) {
	p.Coordinates = cc
}
func (p *PieceAttr) Attr() *PieceAttr {
	return p
}

func (p *PieceAttr) canMoveSetTrue(
	board [][]ChessPiece,
	cc *ChessCoord,
	ml [][]bool,
	move bool,
	capture bool,
) {
	if board[cc.Row][cc.Column] == nil {
		if move {
			ml[cc.Row][cc.Column] = true
		}
	} else if capture && board[cc.Row][cc.Column].Color() != p.Color() {
		ml[cc.Row][cc.Column] = true
	}
}
