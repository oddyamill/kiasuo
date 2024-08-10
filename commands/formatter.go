package commands

type Formatter interface {
	Title(title string) string
	Item(item string) string
	Bold(text string) string
}

type TelegramFormatter struct{}

func (f *TelegramFormatter) Title(title string) string {
	return "*" + title + "*\n"
}

func (f *TelegramFormatter) Item(item string) string {
	return "`•` " + item + "\n"
}

func (f *TelegramFormatter) Bold(text string) string {
	return "*" + text + "*"
}

type DiscordFormatter struct{}

func (f *DiscordFormatter) Title(title string) string {
	return "# " + title + "\n"
}

func (f *DiscordFormatter) Item(item string) string {
	return "- " + item + "\n"
}

func (f *DiscordFormatter) Bold(text string) string {
	return "**" + text + "**"
}
