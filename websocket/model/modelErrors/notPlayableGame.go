package modelErrors

type NotPlayableError struct{}

func (e *NotPlayableError) Error() string {
	return "Game is not playable now"
}
