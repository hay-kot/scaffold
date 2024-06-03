package scaffoldrc

type RCValidationError struct {
	Key   string
	Cause error
}

type RcValidationErrors []RCValidationError

func (e RcValidationErrors) Error() string {
	return "invalid scaffold rc"
}
