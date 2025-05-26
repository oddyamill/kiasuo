package commands

import (
	"testing"
)

func TestParseTelegramKeyboard(t *testing.T) {
	keyboard := Keyboard{
		{
			{
				Text:     "text",
				Callback: "callback",
			},
		},
	}

	result := ParseTelegramKeyboard(keyboard)

	if len(result.InlineKeyboard) != 1 {
		t.Errorf("len(ParseTelegramKeyboard()) = %d; want 1", len(result.InlineKeyboard))
	}

	row := result.InlineKeyboard[0]

	if len(row) != 1 {
		t.Errorf("len(ParseTelegramKeyboard()[0]) = %d; want 1", len(result.InlineKeyboard[0]))
	}

	button := row[0]

	if button.Text != "text" {
		t.Errorf("ParseTelegramKeyboard() = %s; want text", button.Text)
	}

	if *button.CallbackData != "callback" {
		t.Errorf("ParseTelegramKeyboard() = %s; want callback", *button.CallbackData)
	}
}
