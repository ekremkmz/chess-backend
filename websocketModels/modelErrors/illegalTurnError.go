package modelErrors

type IllegalTurn struct {
}

func (e *IllegalTurn) Error() string {
	return "Illegal player turn."
}
