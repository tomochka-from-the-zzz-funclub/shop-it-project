package errors

type basicError struct {
	Cause string `json:"cause"`
	err   string
	code  int
}

func (e *basicError) Error() string {
	result := e.err
	if e.Cause != "" {
		result = result + ": " + e.Cause
	}

	return result
}

func (e *basicError) SetCause(cause string) *basicError {
	e.Cause = cause
	return e
}

func (e *basicError) SetCode(code int) *basicError {
	e.code = code
	return e
}

func (e *basicError) Code() int {
	return e.code
}

func new(err string, cause string, code int) *basicError {
	return &basicError{Cause: cause, err: err, code: code}
}
