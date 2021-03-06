package env_test

import (
	"fmt"
	"os"
	"reflect"
	"syscall"
	"testing"

	"github.com/ab22/env"
)

const (
	defaultString = "DefaultStringVal"
	defaultInt    = 1234
	defaultFloat  = float32(4321.12)
	defaultBool   = true
)

// Struct to test env variables and default env values.
type SupportedTypes struct {
	StringType string  `env:"STRING_VAR" envDefault:"DefaultStringVal"`
	IntType    int     `env:"INT_VAR" envDefault:"1234"`
	FloatType  float32 `env:"FLOAT_VAR" envDefault:"4321.12"`
	BoolType   bool    `env:"BOOL_VAR" envDefault:"true"`
}

// Struct with no env tags set.
type NoTagValues struct {
	StringType string
	IntType    int
	FloatType  float32
	BoolType   bool

	AnotherStruct struct {
		IntType   int
		FloatType float32
	}
}

// Struct to test recursive structs.
type RecursiveStruct struct {
	StringType string `env:"STRING_VAR" envDefault:"DefaultStringVal"`

	AnotherStruct struct {
		IntType   int     `env:"INT_VAR" envDefault:"1234"`
		FloatType float32 `env:"FLOAT_VAR" envDefault:"4321.12"`
	}
}

// Struct to test unaccesible struct
type RecursiveUnexportedStruct struct {
	anotherStruct struct {
		intType int
	}
}

// Struct with no accesible fields.
type NonaccesibleFields struct {
	unexportedField string `env:"STRING_VAR" envDefault:"DefaultStringVal"`
}

// Struct with a unsupported field kind.
type UnSupportedField struct {
	Numbers []int `env:"NUMBERS" envDefault:"1,2,3"`
}

// Test when parsing invalid types.
// env.Parse should return an *env.ErrInvalidInterface.
func TestInvalidInterfaces(t *testing.T) {
	values := []interface{}{
		"string type",
		1234,
		float32(4321.12),
		true,
		nil,
		SupportedTypes{},
	}

	for _, v := range values {
		// Parse by value
		if err := env.Parse(v); err == nil {
			t.Errorf("Accepting invalid type: %v", reflect.TypeOf(v))
		} else if err != env.ErrInvalidInterface {
			t.Fatalf("Error should be of *env.ErrInvalidInterface but was: %v", reflect.TypeOf(v))
		}

		// Parse by reference
		if err := env.Parse(&v); err == nil {
			t.Errorf("env.Parse accepting invalid type by reference: %v", reflect.TypeOf(&v))
		} else if err != env.ErrInvalidInterface {
			t.Fatalf("Error should be of *env.ErrInvalidInterface but was: %v", reflect.TypeOf(v))
		}
	}
}

// Test when parsing env variables into a struct.
// It should set the env values into the structure.
func TestEnvironmentValues(t *testing.T) {
	stringVar := "string value"
	intVar := 6789
	floatVar := float32(1234.56)
	boolVar := true
	s := &SupportedTypes{}

	os.Setenv("STRING_VAR", stringVar)
	os.Setenv("INT_VAR", fmt.Sprintf("%d", intVar))
	os.Setenv("FLOAT_VAR", fmt.Sprintf("%f", floatVar))
	os.Setenv("BOOL_VAR", fmt.Sprintf("%t", boolVar))
	defer func() {
		syscall.Clearenv()
	}()

	if err := env.Parse(s); err != nil {
		t.Fatal(err.Error())
	}

	if s.BoolType != boolVar {
		t.Errorf("Test env values: bool value was not set properly. Expected: [%v] but was [%v]", boolVar, s.BoolType)
	}

	if s.FloatType != floatVar {
		t.Errorf("Test env values: float value was not set properly. Expected: [%v] but was [%v]", floatVar, s.FloatType)
	}

	if s.IntType != intVar {
		t.Errorf("Test env values: int value was not set properly. Expected: [%v] but was [%v]", intVar, s.IntType)
	}

	if s.StringType != stringVar {
		t.Errorf("Test env values: string value was not set properly. Expected: [%v] but was [%v]", stringVar, s.StringType)
	}

}

// Test a struct with no env value set.
// It should set the default values of the structure to the fields.
func TestEnvDefaultValues(t *testing.T) {
	s := &SupportedTypes{}

	if err := env.Parse(s); err != nil {
		t.Fatalf("Parsing struct by reference: %v", err.Error())
	}

	if s.BoolType != defaultBool {
		t.Errorf("Test envDefault values: bool value was not set properly. Expected: [%v] but was [%v]", defaultBool, s.BoolType)
	}

	if s.FloatType != defaultFloat {
		t.Errorf("Test envDefault values: float value was not set properly. Expected: [%v] but was [%v]", defaultFloat, s.FloatType)
	}

	if s.IntType != defaultInt {
		t.Errorf("Test envDefault values: int value was not set properly. Expected: [%v] but was [%v]", defaultInt, s.IntType)
	}

	if s.StringType != defaultString {
		t.Errorf("Test envDefault values: string value was not set properly. Expected: [%v] but was [%v]", defaultString, s.StringType)
	}
}

// Test a struct with no env or envDefault tags set.
// It should not set any values to it and the values in the struct
// should be the default values that Go sets up.
func TestNoTagsSet(t *testing.T) {
	s := &NoTagValues{}

	if err := env.Parse(s); err != nil {
		t.Fatalf("Error parsing struct with no tags set: %v", err.Error())
	}

	if s.BoolType != false {
		t.Errorf("Test default values: bool value was not set properly. Expected: [%v] but was [%v]", false, s.BoolType)
	}

	if s.FloatType != 0 {
		t.Errorf("Test default values: float value was not set properly. Expected: [%v] but was [%v]", 0, s.FloatType)
	}

	if s.IntType != 0 {
		t.Errorf("Test default values: int value was not set properly. Expected: [%v] but was [%v]", 0, s.IntType)
	}

	if s.StringType != "" {
		t.Errorf("Test default values: string value was not set properly. Expected: [%v] but was [%v]", "", s.StringType)
	}

	if s.AnotherStruct.FloatType != 0 {
		t.Errorf("Test default values: float value was not set properly. Expected: [%v] but was [%v]", 0, s.AnotherStruct.FloatType)
	}

	if s.AnotherStruct.IntType != 0 {
		t.Errorf("Test default values: int value was not set properly. Expected: [%v] but was [%v]", 0, s.AnotherStruct.IntType)
	}
}

// Test a struct with no accesible fields.
// It should return a *ErrFieldMustBeAssignable.
func TestNonAssignableFields(t *testing.T) {
	var err error
	s := &NonaccesibleFields{}

	if err = env.Parse(s); err != nil {
		_, ok := err.(env.ErrFieldMustBeAssignable)

		if ok {
			return
		}
	}

	t.Errorf("Expected: [*env.FieldMustBeAssignable] but got: [%s] ", err)
}

// Test a struct with a unsupported type
// It should return a *ErrUnsupportedFieldKind.
func TestUnSupportedField(t *testing.T) {
	var err error
	s := &UnSupportedField{}

	if err = env.Parse(s); err != nil {
		_, ok := err.(*env.ErrUnsupportedFieldKind)

		if ok {
			return
		}
	}

	t.Errorf("Expected: [*env.ErrUnsupportedFieldKind] but got: [%s] ", err)
}

// Test when parsing a struct that has an struct. Values should be set
// recursively into the embedded struct.
func TestRecursiveStructs(t *testing.T) {
	s := &RecursiveStruct{}

	if err := env.Parse(s); err != nil {
		t.Fatalf("Parsing struct by reference: %v", err.Error())
	}

	if s.StringType != defaultString {
		t.Errorf("Test recursive structs: string value was not set properly. Expected: [%v] but was [%v]", defaultString, s.StringType)
	}

	if s.AnotherStruct.IntType != defaultInt {
		t.Errorf("Test recursive structs: int value was not set properly. Expected: [%v] but was [%v]", defaultInt, s.AnotherStruct.IntType)
	}

	if s.AnotherStruct.FloatType != defaultFloat {
		t.Errorf("Test recursive structs: float value was not set properly. Expected: [%v] but was [%v]", defaultFloat, s.AnotherStruct.FloatType)
	}
}

// Test when parsing a struct that has an unaccesible/unexported field struct.
func TestRecursiveUnaccesibleStruct(t *testing.T) {
	var err error
	s := &RecursiveUnexportedStruct{}

	if err = env.Parse(s); err != nil {
		_, ok := err.(env.ErrFieldMustBeAssignable)

		if ok {
			return
		}
	}

	t.Errorf("Expected: [*env.FieldMustBeAssignable] but got: [%s] ", err)
}
