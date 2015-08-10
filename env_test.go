package env_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/ab22/env"
)

const (
	defaultString = "DefaultStringVal"
	defaultInt    = 1234
	defaultFloat  = float32(4321.12)
	defaultBool   = true
)

// Struct to test env variables and default env values
type SupportedTypesStruct struct {
	StringType string  `env:"STRING_VAR" envDefault:"DefaultStringVal"`
	IntType    int     `env:"INT_VAR" envDefault:"1234"`
	FloatType  float32 `env:"FLOAT_VAR" envDefault:"4321.12"`
	BoolType   bool    `env:"BOOL_VAR" envDefault:"true"`
}

// Struct with no env tags set
type NoTagValuesStruct struct {
	StringType string
	IntType    int
	FloatType  float32
	BoolType   bool
}

// Test when parsing invalid types.
// env.Parse should return an *env.InvalidInterfaceError
func TestInvalidInterfaces(t *testing.T) {
	values := []interface{}{
		"string type",
		1234,
		float32(4321.12),
		true,
		nil,
		SupportedTypesStruct{},
	}

	for _, v := range values {
		// Parse by value
		if err := env.Parse(v); err == nil {
			t.Errorf("Accepting invalid type:", reflect.TypeOf(v))
		} else if _, ok := err.(*env.InvalidInterfaceError); !ok {
			t.Fatal("Error should be of *env.InvalidInterfaceError but was:", reflect.TypeOf(v))
		}

		// Parse by reference
		if err := env.Parse(&v); err == nil {
			t.Errorf("env.Parse accepting invalid type by reference:", reflect.TypeOf(&v))
		} else if _, ok := err.(*env.InvalidInterfaceError); !ok {
			t.Fatal("Error should be of *env.InvalidInterfaceError but was:", reflect.TypeOf(v))
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
	s := &SupportedTypesStruct{}

	os.Setenv("STRING_VAR", stringVar)
	os.Setenv("INT_VAR", fmt.Sprintf("%d", intVar))
	os.Setenv("FLOAT_VAR", fmt.Sprintf("%f", floatVar))
	os.Setenv("BOOL_VAR", fmt.Sprintf("%t", boolVar))
	defer func() {
		os.Unsetenv("STRING_VAR")
		os.Unsetenv("INT_VAR")
		os.Unsetenv("FLOAT_VAR")
		os.Unsetenv("BOOL_VAR")
	}()

	if err := env.Parse(s); err != nil {
		t.Fatal(err.Error())
	}

	if s.BoolType != boolVar {
		t.Errorf("Test default values: bool value was not set properly. Expected: [%v] but was [%v]", boolVar, s.BoolType)
	}

	if s.FloatType != floatVar {
		t.Errorf("Test default values: float value was not set properly. Expected: [%v] but was [%v]", floatVar, s.FloatType)
	}

	if s.IntType != intVar {
		t.Errorf("Test default values: int value was not set properly. Expected: [%v] but was [%v]", intVar, s.IntType)
	}

	if s.StringType != stringVar {
		t.Errorf("Test default values: string value was not set properly. Expected: [%v] but was [%v]", stringVar, s.StringType)
	}

}

// Test a struct with no env value set.
// It should set the default values of the structure to the fields.
func TestDefaultValues(t *testing.T) {
	s := &SupportedTypesStruct{}

	if err := env.Parse(s); err != nil {
		t.Fatalf("Parsing struct by reference: %v", err.Error())
	}

	if s.BoolType != defaultBool {
		t.Errorf("Test default values: bool value was not set properly. Expected: [%v] but was [%v]", defaultBool, s.BoolType)
	}

	if s.FloatType != defaultFloat {
		t.Errorf("Test default values: float value was not set properly. Expected: [%v] but was [%v]", defaultFloat, s.FloatType)
	}

	if s.IntType != defaultInt {
		t.Errorf("Test default values: int value was not set properly. Expected: [%v] but was [%v]", defaultInt, s.IntType)
	}

	if s.StringType != defaultString {
		t.Errorf("Test default values: string value was not set properly. Expected: [%v] but was [%v]", defaultString, s.StringType)
	}
}

// Test a struct with no env or envDefault tags set.
// It should not set any values to it and the values in the struct
// should be the default values that Go sets up.
func TestNoTagsSet(t *testing.T) {
	s := &NoTagValuesStruct{}

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
}
