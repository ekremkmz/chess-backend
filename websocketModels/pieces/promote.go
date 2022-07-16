package pieces

type Promote string

const (
	UnknownPromotion Promote = ""
	QueenPromotion   Promote = "q"
	RookPromotion    Promote = "r"
	BishopPromotion  Promote = "b"
	KnightPromotion  Promote = "n"
)
