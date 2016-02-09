package env

import (
	"errors"
	"fmt"
)

var ErrInvalidInterface = errors.New("env: struct parsing: expected pointer to struct")

type ErrUnsupportedFieldKind struct {
	FieldName string
	FieldKind string
}

func (e ErrUnsupportedFieldKind) Error() string {
	return fmt.Sprintf("env: set value '%s': unsupported field kind '%s'", e.FieldName, e.FieldKind)
}

type ErrFieldMustBeAssignable string

func (field ErrFieldMustBeAssignable) Error() string {
	return fmt.Sprintf("env: set value '%s': cannot set value to unexported field", string(field))
}
