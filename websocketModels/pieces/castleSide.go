package pieces

type CastleSide string

func (c CastleSide) String() string {
	return string(c)
}

func (c CastleSide) KingCoord() ChessCoord {
	return _cs2king[c]
}

func (c CastleSide) RookCoord() ChessCoord {
	return _cs2rook[c]
}

func (c CastleSide) RookAfterCastlingCoord() ChessCoord {
	return _cs2rookAfterCastling[c]
}

var _cs2king = map[CastleSide]ChessCoord{
	"K": {Row: 7, Column: 6},
	"Q": {Row: 7, Column: 2},
	"k": {Row: 0, Column: 6},
	"q": {Row: 0, Column: 2},
}

var _cs2rook = map[CastleSide]ChessCoord{
	"K": {Row: 7, Column: 7},
	"Q": {Row: 7, Column: 0},
	"k": {Row: 0, Column: 7},
	"q": {Row: 0, Column: 0},
}

var _cs2rookAfterCastling = map[CastleSide]ChessCoord{
	"K": {Row: 7, Column: 5},
	"Q": {Row: 7, Column: 3},
	"k": {Row: 0, Column: 5},
	"q": {Row: 0, Column: 3},
}
