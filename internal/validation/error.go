package validation

import (
	"errors"
	"fmt"
)

var (
	ErrValidation = errors.New("validation error")

	ErrEmptyDesc          = errors.New("description cant be empty")
	ErrMissedDeadline     = errors.New("deadline cant be already missed")
	ErrUnchangeableStatus = errors.New("cant change done or canceled tasks")
)

type InvalidFiled struct {
	Err   error
	Field string
}

type ValidationError struct {
	Err    error
	Fields []InvalidFiled
}

func (e *ValidationError) Error() string {
	errStr := ""
	for _, f := range e.Fields {
		errStr += fmt.Sprintf("%s: %v\n", f.Field, f.Err)
	}

	return errStr
}
