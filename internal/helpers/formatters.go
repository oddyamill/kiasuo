package helpers

import "strings"

type Formatter interface {
	Title(title string) string
	Item(item string) string
	Bold(text string) string
	Code(text string) string
	Block(text string) string
	Line(text string) string
	Link(text, url string) string
}

type TelegramFormatter struct{}

func (f TelegramFormatter) Title(title string) string {
	return "*" + title + "*\n"
}

func (f TelegramFormatter) Item(item string) string {
	return "`•` " + item + "\n"
}

func (f TelegramFormatter) Code(text string) string {
	return "`" + text + "`"
}

func (f TelegramFormatter) Block(text string) string {
	return "`" + strings.ReplaceAll(text, "\n", "\n ") + "`"
}

func (f TelegramFormatter) Bold(text string) string {
	return "*" + text + "*"
}

func (f TelegramFormatter) Line(text string) string {
	return text + "\n"
}

func (f TelegramFormatter) Link(text, url string) string {
	return "[" + text + "](" + url + ")"
}
