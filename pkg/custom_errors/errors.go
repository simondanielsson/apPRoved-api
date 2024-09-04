package customerrors

import (
	"fmt"
	"log"
)

type ValidationError struct {
	Field string
	Msg   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Msg)
}

func NewValidationError(field, msg string) error {
	return &ValidationError{Field: field, Msg: msg}
}

func IgnoreError(e error) {
	if e != nil {
		log.Printf("Ignored error %v", e.Error())
	}
}
