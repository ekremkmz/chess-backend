package modelErrors

type PlayerNotFoundError struct {
	Nick string
}

func (e *PlayerNotFoundError) Error() string {
	return "Player with" + e.Nick + "id not found."
}
