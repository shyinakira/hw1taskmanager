package domain

import (
	"errors"
)

var (
	ErrUnknowmCommand = errors.New("unknowm command")
	ErrWrongArgsCount = errors.New("invalid arguments provided")
	ErrInvalidArgument = errors.New("invalid argumnet")
	
	ErrTaskNotFound     = errors.New("task with this ID not found")

	ErrSaving    = errors.New("failed to save tasks")
	ErrRestoring = errors.New("failed to restore tasks")
)