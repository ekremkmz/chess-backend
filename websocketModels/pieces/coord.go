package pieces

import "strconv"

type ChessCoord struct {
	Row    int64
	Column int64
}

func (c *ChessCoord) Add(other *ChessCoord) *ChessCoord {
	new := &ChessCoord{Row: c.Row + other.Row, Column: c.Column + other.Column}

	if new.Row < 0 || new.Row > 7 || new.Column < 0 || new.Column > 7 {
		return nil
	}
	return new
}

func (c *ChessCoord) ToString() string {
	return intToChar[c.Row] + strconv.FormatInt(c.Column+1, 10)
}

func NewChessCoordFromString(str string) *ChessCoord {
	if str == "" {
		return nil
	}
	x, _ := charToInt[string(str[0])]
	y, _ := strconv.ParseInt(string(str[1]), 10, 64)
	return &ChessCoord{Row: 8 - y, Column: x - 1}
}

var charToInt = map[string]int64{
	"a": 1,
	"b": 2,
	"c": 3,
	"d": 4,
	"e": 5,
	"f": 6,
	"g": 7,
	"h": 8,
}

var intToChar = map[int64]string{
	1: "a",
	2: "b",
	3: "c",
	4: "d",
	5: "e",
	6: "f",
	7: "g",
	8: "h",
}
