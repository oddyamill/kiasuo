package helpers

import "testing"

var (
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

func TestTelegramFormatterBlock(t *testing.T) {
	result := telegramFormatter.Block("test")

	if result != "`test`" {
		t.Errorf("TelegramFormatter.Block() = %s; want `test`", result)
	}

	result = telegramFormatter.Block("test\ntest")

	if result != "`test\n test`" {
		t.Errorf("TelegramFormatter.Block() = %s; want `test\n test`", result)
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
