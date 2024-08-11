package helpers

func If[T any](condition bool, truth, falsity T) T {
	if condition {
		return truth
	}
	return falsity
}
