package helpers

import (
	"os"
	"strings"
)

func If[T any](condition bool, truth, falsity T) T {
	if condition {
		return truth
	}

	return falsity
}

func GetEnv(key string) string {
	value, ok := os.LookupEnv(key)

	if ok {
		return value
	}

	value, ok = os.LookupEnv(key + "_FILE")

	if ok {
		buffer, err := os.ReadFile(value)

		if err != nil {
			panic(err)
		}

		return strings.TrimSpace(string(buffer))
	}

	panic("Environment variable " + key + " not set")
}

func IsTesting() bool {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			return true
		}
	}
	return false
}

func HumanizeLesson(lesson string) string {
	lesson = strings.TrimSpace(lesson)

	switch strings.ToLower(lesson) {
	case "основы безопасности жизнедеятельности", "основы безопасности и защиты родины":
		return "ОБЖ"
	case "физическая культура":
		return "Физ-ра"
	case "алгебра и начала математического анализа":
		return "Алгебра"
	case "искусственный интеллект":
		return "ИИ"
	default:
		return lesson
	}
}
