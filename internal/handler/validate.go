package handler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// Validate validates a struct and returns a human-readable error string
// suitable for use as a 400 Bad Request message, or nil if valid.
func Validate(v any) error {
	err := validate.Struct(v)
	if err == nil {
		return nil
	}

	var errs validator.ValidationErrors
	if !errors.As(err, &errs) {
		return err
	}

	msgs := make([]string, 0, len(errs))
	for _, fe := range errs {
		msgs = append(msgs, fmt.Sprintf("%s: failed %s validation", fe.Field(), fe.Tag()))
	}
	return fmt.Errorf("%s", strings.Join(msgs, "; "))
}
