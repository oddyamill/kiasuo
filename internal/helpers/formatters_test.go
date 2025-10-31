package helpers

import "testing"

var (
	telegramFormatter = Formatter{}
)

func TestTelegramFormatterTitle(t *testing.T) {
	result := telegramFormatter.Title("test")

	if result != "*test*\n" {
		t.Errorf("Formatter.Title() = %s; want *test*\n", result)
	}
}

func TestTelegramFormatterItem(t *testing.T) {
	result := telegramFormatter.Item("test")

	if result != "`•` test\n" {
		t.Errorf("Formatter.Item() = %s; want `•` test\n", result)
	}
}

func TestTelegramFormatterBold(t *testing.T) {
	result := telegramFormatter.Bold("test")

	if result != "*test*" {
		t.Errorf("Formatter.Bold() = %s; want *test*", result)
	}
}

func TestTelegramFormatterCode(t *testing.T) {
	result := telegramFormatter.Code("test")

	if result != "`test`" {
		t.Errorf("Formatter.Code() = %s; want `test`", result)
	}
}

func TestTelegramFormatterBlock(t *testing.T) {
	result := telegramFormatter.Block("test")

	if result != "`test`" {
		t.Errorf("Formatter.Block() = %s; want `test`", result)
	}

	result = telegramFormatter.Block("test\ntest")

	if result != "`test\n test`" {
		t.Errorf("Formatter.Block() = %s; want `test\n test`", result)
	}
}

func TestTelegramFormatterLine(t *testing.T) {
	result := telegramFormatter.Line("test")

	if result != "test\n" {
		t.Errorf("Formatter.Line() = %s; want test\n", result)
	}
}

func TestTelegramFormatterLink(t *testing.T) {
	result := telegramFormatter.Link("test", "http://example.com")

	if result != "[test](http://example.com)" {
		t.Errorf("Formatter.Link() = %s; want [test](http://example.com)", result)
	}
}
