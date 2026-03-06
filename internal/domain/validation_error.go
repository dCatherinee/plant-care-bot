package domain

import "fmt"

type ValidationError struct {
	Field   string
	Problem string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation: %s: %s", e.Field, e.Problem)
}

func (e ValidationError) Unwrap() error {
	return ErrInvalidArgument
}
