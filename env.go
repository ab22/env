package env

import (
	"os"
	"reflect"
	"strconv"
)

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
		return nil, &InvalidInterfaceError{}
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

	if interfaceKind != reflect.Ptr {
		return true
	}

	interfaceKind = interfaceValue.Elem().Kind()

	return interfaceKind != reflect.Struct
}

func setStructValues(structElem *reflect.Value) error {
	numFields := structElem.NumField()
	structType := structElem.Type()

	for i := 0; i < numFields; i++ {
		structField := structType.Field(i)
		fieldValue := structElem.Field(i)
		envValue := getEnvOrDefaultValue(&structField)

		if envValue == "" {
			continue
		}

		if err := setValue(&fieldValue, structField.Name, envValue); err != nil {
			return err
		}
	}

	return nil
}

func getEnvOrDefaultValue(field *reflect.StructField) string {
	envKey := field.Tag.Get("env")
	defaultValue := field.Tag.Get("envDefault")

	if envKey == "" {
		return defaultValue
	}

	envValue := os.Getenv(envKey)

	if envValue == "" {
		return defaultValue
	}

	return envValue
}

func setValue(field *reflect.Value, fieldName string, envValue string) error {
	if !field.CanSet() {
		return &FieldMustBeAssignableError{FieldName: fieldName}
	}

	fieldKind := field.Kind()

	switch fieldKind {
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
		return &UnsupportedFieldKindError{
			FieldName: fieldName,
			FieldKind: fieldKind.String(),
		}
	}

	return nil
}
