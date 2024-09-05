package commands

import "testing"

func TestFormatTeacher(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"Иванов Иван Иванович (Физическая культура, Основы безопасности жизнедеятельности, Алгебра и начала математического анализа, География)",
			"Иванов Иван Иванович (Физ-ра, ОБЖ, Алгебра, География)",
		},
		{
			input:    "Дмитриев Дмитрий Дмитриевич",
			expected: "Дмитриев Дмитрий Дмитриевич",
		},
	}

	for _, test := range tests {
		result := formatTeacher(test.input)

		if result != test.expected {
			t.Errorf("formatTeacher(%s) = %s; want %s", test.input, result, test.expected)
		}
	}
}
