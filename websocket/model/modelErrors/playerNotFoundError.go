package modelErrors

type PlayerNotFoundError struct {
	Id string
}

func (e *PlayerNotFoundError) Error() string {
	return "Player with" + e.Id + "id not found."
}
