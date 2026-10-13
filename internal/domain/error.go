package domain

import (
	"errors"
)

var (
	ErrUnclosedQuote = errors.New("quote is not closed")

	ErrTimezone        = errors.New("unknowm timezone")
	ErrUnknownCommand  = errors.New("unknowm command")
	ErrWrongArgsCount  = errors.New("invalid arguments provided")
	ErrInvalidArgument = errors.New("invalid argument")

	ErrTaskNotFound = errors.New("task with this ID not found")

	ErrSaving    = errors.New("failed to save tasks")
	ErrRestoring = errors.New("failed to restore tasks")
)
