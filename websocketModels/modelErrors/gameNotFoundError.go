package modelErrors

type GameNotFoundError struct{}

func (e *GameNotFoundError) Error() string {
	return "Game not found."
}
