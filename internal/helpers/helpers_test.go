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

func TestStringToBytes(t *testing.T) {
	result := StringToBytes("test")

	if string(result) != "test" {
		t.Errorf("StringToBytes() = %s; want test", result)
	}
}

func TestBytesToString(t *testing.T) {
	result := BytesToString([]byte("test"))

	if result != "test" {
		t.Errorf("BytesToString() = %s; want test", result)
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
			t.Errorf("GetEnv() with non-existent variable did not panic")
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
		t.Errorf("GetEnv() = %s; want test", value)
	}
}

func TestIsTesting(t *testing.T) {
	if !IsTesting() {
		t.Errorf("IsTesting() = false; want true")
	}
}

func TestHumanizeLesson(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Основы безопасности жизнедеятельности", "ОБЖ"},
		{"Основы безопасности и защиты родины", "ОБЖ"},
		{"Физическая культура", "Физкультура"},
		{"Алгебра и начала математического анализа", "Алгебра"},
		{"Россия - мои горизонты", "Профориентация"},
		{"Unknown", "Unknown"},
	}

	for _, test := range tests {
		result := HumanizeLesson(test.input)

		if result != test.expected {
			t.Errorf("HumanizeLesson(%s) = %s; want %s", test.input, result, test.expected)
		}
	}
}
