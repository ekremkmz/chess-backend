package modelErrors

type TimeOverError struct{}

func (e *TimeOverError) Error() string {
	return "Time out."
}
