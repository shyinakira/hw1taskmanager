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
	ErrSameStatus         = errors.New("new status cant be same")
)

type InvalidField struct {
	Err   error
	Field string
}

type ValidationError struct {
	Err    error
	Fields []InvalidField
}

func (e *ValidationError) Error() string {
	errStr := ""
	for _, f := range e.Fields {
		errStr += fmt.Sprintf("%s: %v\n", f.Field, f.Err)
	}

	return errStr
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}
