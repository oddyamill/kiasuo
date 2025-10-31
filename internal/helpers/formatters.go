package helpers

import "strings"

type Formatter struct{}

func (f Formatter) Title(title string) string {
	return "*" + title + "*\n"
}

func (f Formatter) Item(item string) string {
	return "`•` " + item + "\n"
}

func (f Formatter) Code(text string) string {
	return "`" + text + "`"
}

func (f Formatter) Block(text string) string {
	return "`" + strings.ReplaceAll(text, "\n", "\n ") + "`"
}

func (f Formatter) Bold(text string) string {
	return "*" + text + "*"
}

func (f Formatter) Line(text string) string {
	return text + "\n"
}

func (f Formatter) Link(text, url string) string {
	return "[" + text + "](" + url + ")"
}
