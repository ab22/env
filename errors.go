package env

import (
	"errors"
	"fmt"
)

var InvalidInterfaceError = errors.New("env: struct parsing: expected struct or pointer to struct")

type UnsupportedFieldKindError struct {
	FieldName string
	FieldKind string
}

func (e UnsupportedFieldKindError) Error() string {
	return fmt.Sprintf("env: set value '%s': unsupported field kind '%s'", e.FieldName, e.FieldKind)
}

type FieldMustBeAssignableError struct {
	FieldName string
}

func (e FieldMustBeAssignableError) Error() string {
	return fmt.Sprintf("env: set value '%s': cannot set value to unexported field", e.FieldName)
}
