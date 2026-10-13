package validation

import (
	"hw1taskmanager/internal/domain"
	"time"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateTaskData(desc string, deadline time.Time) error {
	invalidFields := make([]InvalidFiled, 0)
	if desc == "" {
		invalidFields = append(invalidFields, InvalidFiled{
			Err:   ErrEmptyDesc,
			Field: "description",
		})
	}

	if deadline.Before(time.Now()) {
		invalidFields = append(invalidFields, InvalidFiled{
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

func (v *Validator) CanChangeStatus(status domain.Status) error {
	if status == domain.DONE || status == domain.CANCELED {
		return &ValidationError{
			Err: ErrValidation,
			Fields: []InvalidFiled{
				{
					Err:   ErrUnchangeableStatus,
					Field: "status",
				},
			},
		}
	}

	return nil
}
