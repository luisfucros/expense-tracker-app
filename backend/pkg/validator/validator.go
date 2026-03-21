package validator

import (
	"github.com/go-playground/validator/v10"

	"github.com/luisfucros/expense-tracker-app/internal/domain/apierror"
)

var validate = validator.New()

// Validate runs struct-level validation on s and returns an APIError on failure.
func Validate(s any) error {
	if err := validate.Struct(s); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			return apierror.BadRequest(apierror.CodeValidation, formatValidationErrors(ve))
		}
		return err
	}
	return nil
}

func formatValidationErrors(ve validator.ValidationErrors) string {
	for _, fe := range ve {
		return fe.Field() + ": " + fe.Tag() + " validation failed"
	}
	return "validation failed"
}
