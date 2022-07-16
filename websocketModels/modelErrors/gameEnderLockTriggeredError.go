package modelErrors

type GameEnderLockTriggeredError struct {
}

func (e *GameEnderLockTriggeredError) Error() string {
	return "Game ended."
}
