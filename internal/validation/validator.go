package validation

import (
	"strings"
	"time"

	"hw1taskmanager/internal/domain"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) CanEdit(status domain.Status) error {
	if status == domain.CANCELED || status == domain.DONE {
		return &ValidationError{
			Err: ErrValidation,
			Fields: []InvalidField{
				{
					Err:   ErrUnchangeableStatus,
					Field: "status",
				},
			},
		}
	}

	return nil
}

func (v *Validator) ValidateTaskData(title string, deadline time.Time) error {
	invalidFields := make([]InvalidField, 0)
	if strings.TrimSpace(title) == "" {
		invalidFields = append(invalidFields, InvalidField{
			Err:   ErrEmptyDesc,
			Field: "title",
		})
	}

	if !deadline.IsZero() && deadline.Before(time.Now()) {
		invalidFields = append(invalidFields, InvalidField{
			Err:   ErrMissedDeadline,
			Field: "deadline",
		})
	}

	if len(invalidFields) != 0 {
		return &ValidationError{
			Err:    ErrValidation,
			Fields: invalidFields,
		}
	}

	return nil
}

func (v *Validator) CanChangeStatus(oldStatus, newStatus domain.Status) error {
	if newStatus == oldStatus {
		return &ValidationError{
			Err: ErrValidation,
			Fields: []InvalidField{
				{
					Err:   ErrSameStatus,
					Field: "status",
				},
			},
		}
	}

	if oldStatus == domain.DONE || oldStatus == domain.CANCELED {
		return &ValidationError{
			Err: ErrValidation,
			Fields: []InvalidField{
				{
					Err:   ErrUnchangeableStatus,
					Field: "status",
				},
			},
		}
	}

	return nil
}
