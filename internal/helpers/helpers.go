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
