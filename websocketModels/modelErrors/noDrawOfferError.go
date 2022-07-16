package modelErrors

type NoDrawOfferError struct{}

func (e *NoDrawOfferError) Error() string {
	return "There is no draw offer"
}
