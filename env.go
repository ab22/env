package env

import (
	"errors"
	"os"
	"reflect"
	"strconv"
)

var InvalidInterfaceError = errors.New("env parse: expected struct or pointer to struct")
var UnsupportedFieldKindError = errors.New("env parse: unsupported field kind")
var FieldMustBeAssignableError = errors.New("env parse: cannot set value to unexported field")

func Parse(i interface{}) error {
	var elem *reflect.Value
	var err error

	if elem, err = getStructureElement(i); err != nil {
		return err
	}

	err = setStructValues(elem)
	return err
}

func getStructureElement(i interface{}) (*reflect.Value, error) {
	if isInvalidInterface(i) {
		return nil, InvalidInterfaceError
	}

	elem := reflect.ValueOf(i).Elem()
	return &elem, nil
}

func isInvalidInterface(i interface{}) bool {
	if i == nil {
		return true
	}

	interfaceValue := reflect.ValueOf(i)
	interfaceKind := interfaceValue.Kind()

	if interfaceKind == reflect.Ptr {
		interfaceKind = interfaceValue.Elem().Kind()
	}

	return interfaceKind != reflect.Struct
}

func setStructValues(structElem *reflect.Value) error {
	numFields := structElem.NumField()
	structType := structElem.Type()

	for i := 0; i < numFields; i++ {
		structField := structType.Field(i)
		fieldValue := structElem.Field(i)
		envValue := getEnvValue(&structField)

		if envValue == "" {
			continue
		}

		if err := setEnvVariable(&fieldValue, envValue); err != nil {
			return err
		}
	}

	return nil
}

func getEnvValue(field *reflect.StructField) string {
	envVariableName := field.Tag.Get("env")
	defaultValue := field.Tag.Get("envDefault")

	if envVariableName == "" {
		return defaultValue
	}

	envVariable := os.Getenv(envVariableName)

	if envVariable == "" {
		return defaultValue
	}

	return envVariable
}

func setEnvVariable(field *reflect.Value, envValue string) error {
	if !field.CanSet() {
		return FieldMustBeAssignableError
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(envValue)
	case reflect.Int:
		intValue, err := strconv.ParseInt(envValue, 10, 32)

		if err != nil {
			return err
		}

		field.SetInt(intValue)
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(envValue)

		if err != nil {
			return err
		}

		field.SetBool(boolValue)
	case reflect.Float32:
		floatValue, err := strconv.ParseFloat(envValue, 32)

		if err != nil {
			return err
		}

		field.SetFloat(floatValue)
	default:
		return UnsupportedFieldKindError
	}

	return nil
}
