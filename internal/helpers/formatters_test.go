package helpers

import "testing"

var (
	discordFormatter  = DiscordFormatter{}
	telegramFormatter = TelegramFormatter{}
)

func TestTelegramFormatterTitle(t *testing.T) {
	result := telegramFormatter.Title("test")

	if result != "*test*\n" {
		t.Errorf("TelegramFormatter.Title() = %s; want *test*\n", result)
	}
}

func TestTelegramFormatterItem(t *testing.T) {
	result := telegramFormatter.Item("test")

	if result != "`•` test\n" {
		t.Errorf("TelegramFormatter.Item() = %s; want `•` test\n", result)
	}
}

func TestTelegramFormatterBold(t *testing.T) {
	result := telegramFormatter.Bold("test")

	if result != "*test*" {
		t.Errorf("TelegramFormatter.Bold() = %s; want *test*", result)
	}
}

func TestTelegramFormatterCode(t *testing.T) {
	result := telegramFormatter.Code("test")

	if result != "`test`" {
		t.Errorf("TelegramFormatter.Code() = %s; want `test`", result)
	}
}

func TestTelegramFormatterLine(t *testing.T) {
	result := telegramFormatter.Line("test")

	if result != "test\n" {
		t.Errorf("TelegramFormatter.Line() = %s; want test\n", result)
	}
}

func TestTelegramFormatterLink(t *testing.T) {
	result := telegramFormatter.Link("test", "http://example.com")

	if result != "[test](http://example.com)" {
		t.Errorf("TelegramFormatter.Link() = %s; want [test](http://example.com)", result)
	}
}

func TestDiscordFormatterTitle(t *testing.T) {
	result := discordFormatter.Title("test")

	if result != "## test\n" {
		t.Errorf("DiscordFormatter.Title() = %s; want ## test\n", result)
	}
}

func TestDiscordFormatterItem(t *testing.T) {
	result := discordFormatter.Item("test")

	if result != "- test\n" {
		t.Errorf("DiscordFormatter.Item() = %s; want - test\n", result)
	}
}

func TestDiscordFormatterBold(t *testing.T) {
	result := discordFormatter.Bold("test")

	if result != "**test**" {
		t.Errorf("DiscordFormatter.Bold() = %s; want **test**", result)
	}
}

func TestDiscordFormatterCode(t *testing.T) {
	result := discordFormatter.Code("test")

	if result != "`test`" {
		t.Errorf("DiscordFormatter.Code() = %s; want `test`", result)
	}
}

func TestDiscordFormatterLine(t *testing.T) {
	result := discordFormatter.Line("test")

	if result != "test\n" {
		t.Errorf("DiscordFormatter.Line() = %s; want test\n", result)
	}
}

func TestDiscordFormatterLink(t *testing.T) {
	result := discordFormatter.Link("test", "http://example.com")

	if result != "[test](http://example.com)" {
		t.Errorf("DiscordFormatter.Link() = %s; want [test](http://example.com)", result)
	}
}
