package helpers

import (
	"os"
	"testing"
)

func TestIfTrue(t *testing.T) {
	result := If(true, "true", "false")

	if result != "true" {
		t.Errorf("If(true, _, _) = %s; want true", result)
	}
}

func TestIfFalse(t *testing.T) {
	result := If(false, "true", "false")

	if result != "false" {
		t.Errorf("If(false, _, _) = %s; want false", result)
	}
}

func TestGetEnvPath(t *testing.T) {
	// I think PATH variable is set on every system
	GetEnv("PATH")
}

const UnknownVariable = "RANDOM_VARIABLE_THAT_DOES_NOT_EXIST"

func TestGetEnvNonExistent(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("GetEnv(%s) did not panic", UnknownVariable)
		}
	}()

	GetEnv(UnknownVariable)
}

func TestGetEnvFile(t *testing.T) {
	file, err := os.CreateTemp("", "random-file.*.txt")
	defer os.Remove(file.Name())

	if err != nil {
		panic(err)
	}

	_, err = file.WriteString("test")

	if err != nil {
		panic(err)
	}

	err = os.Setenv(UnknownVariable+"_FILE", file.Name())

	if err != nil {
		panic(err)
	}

	value := GetEnv(UnknownVariable)

	if value != "test" {
		t.Errorf("GetEnv(%s) = %s; want test", UnknownVariable, value)
	}
}
