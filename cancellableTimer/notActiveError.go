package cancellableTimer

type NotActiveError struct{}

func (e *NotActiveError) Error() string {
	return "Timer is not active."
}
