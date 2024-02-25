package password

type PasswordError struct {
	message string
	error   error
}

func (e *PasswordError) Error() string {
	return e.message
}

func (e *PasswordError) Unwrap() error {
	return e.error
}
