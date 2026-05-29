package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

func FormatValidationError(
	err error,
) string {

	errors := err.(validator.ValidationErrors)

	for _, e := range errors {

		switch e.Tag() {

		case "required":
			return fmt.Sprintf(
				"%s is required",
				e.Field(),
			)

		case "email":
			return fmt.Sprintf(
				"%s must be valid email",
				e.Field(),
			)

		case "min":
			return fmt.Sprintf(
				"%s too short",
				e.Field(),
			)

		case "max":
			return fmt.Sprintf(
				"%s too long",
				e.Field(),
			)
		}
	}

	return err.Error()
}