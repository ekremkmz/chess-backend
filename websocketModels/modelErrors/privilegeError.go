package modelErrors

type PrivilegeError struct{}

func (e *PrivilegeError) Error() string {
	return "You do not have the required privilege"
}
