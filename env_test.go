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

type SupportedTypesStruct struct {
	StringType string  `env:"STRING_VAR" envDefault:"DefaultStringVal"`
	IntType    int     `env:"INT_VAR" envDefault:"1234"`
	FloatType  float32 `env:"FLOAT_VAR" envDefault:"4321.12"`
	BoolType   bool    `env:"BOOL_VAR" envDefault:"true"`
}

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
