package helpers

import "testing"

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
