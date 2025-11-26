package helpers

import (
	"os"
	"strings"
	"unsafe"
)

func If[T any](condition bool, truth, falsity T) T {
	if condition {
		return truth
	}

	return falsity
}

func StringToBytes(str string) []byte {
	return unsafe.Slice(unsafe.StringData(str), len(str))
}

func BytesToString(bytes []byte) string {
	return unsafe.String(unsafe.SliceData(bytes), len(bytes))
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

		return strings.TrimSpace(BytesToString(buffer))
	}

	panic("environment variable " + key + " not set")
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
	case "основы безопасности жизнедеятельности":
		return "ОБЖ"
	case "основы безопасности и защиты родины", "основы безопасности и защиты родины (обзр)":
		return "ОБЗР"
	case "физическая культура":
		return "Физкультура"
	case "алгебра и начала математического анализа":
		return "Алгебра"
	case "искусственный интеллект":
		return "Информатика (ИИ)"
	case "россия - мои горизонты":
		return "Профориентация"
	case "учимся писать сочинение":
		return "Русский язык (сочинение)"
	case "трудные вопросы орфографии и пунктуации", "секреты русского языка":
		return "Русский язык (консультация)"
	case "родной язык (русский язык)":
		return "Родной язык"
	case "труд (технология)":
		return "Труд"
	case "решение задач повышенной сложности":
		return "Математика (консультация)"
	case "основы духовно-нравственной культуры народов россии":
		return "ОДНКНР"
	case "изобразительное искусство":
		return "ИЗО"
	case "мировая художественная культура":
		return "МХК"
	case "основы религиозных культур и светской этики":
		return "ОРКСЭ"
	case "иностранный язык (английский язык)", "иностранный (английский) язык":
		return "Английский язык"
	default:
		return lesson
	}
}

// IsHexUnsafe only checks if the string contains valid hex characters
func IsHexUnsafe(hex string) bool {
	for i := 0; i < len(hex); i++ {
		if (hex[i] < '0' || hex[i] > '9') && (hex[i] < 'a' || hex[i] > 'f') && (hex[i] < 'A' || hex[i] > 'F') {
			return false
		}
	}

	return true
}
