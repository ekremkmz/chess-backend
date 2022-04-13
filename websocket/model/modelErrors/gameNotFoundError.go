package modelErrors

type GameNotFoundError struct{}

func (g *GameNotFoundError) Error() string {
	return "Game not found."
}
