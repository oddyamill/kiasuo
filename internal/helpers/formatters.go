package helpers

type Formatter interface {
	Title(title string) string
	Item(item string) string
	Bold(text string) string
	Code(text string) string
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

func (f TelegramFormatter) Bold(text string) string {
	return "*" + text + "*"
}

func (f TelegramFormatter) Line(text string) string {
	return text + "\n"
}

func (f TelegramFormatter) Link(text, url string) string {
	return "[" + text + "](" + url + ")"
}

type DiscordFormatter struct{}

func (f DiscordFormatter) Title(title string) string {
	return "## " + title + "\n"
}

func (f DiscordFormatter) Item(item string) string {
	return "- " + item + "\n"
}

func (f DiscordFormatter) Bold(text string) string {
	return "**" + text + "**"
}

func (f DiscordFormatter) Code(text string) string {
	return "`" + text + "`"
}

func (f DiscordFormatter) Line(text string) string {
	return text + "\n"
}

func (f DiscordFormatter) Link(text, url string) string {
	return "[" + text + "](" + url + ")"
}
