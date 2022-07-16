package cancellableTimer

type Timeoutable interface {
	WhenTimeout()
}
