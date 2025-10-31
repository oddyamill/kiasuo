package commands

import (
	"testing"

	"github.com/go-telegram/bot/models"
)

func TestKeyboardParse(t *testing.T) {
	keyboard := Keyboard{
		{
			{
				Text:     "text",
				Callback: "callback",
			},
		},
	}

	result := keyboard.Parse().(*models.InlineKeyboardMarkup)

	if len(result.InlineKeyboard) != 1 {
		t.Errorf("len(keyboard.Parse()) = %d; want 1", len(result.InlineKeyboard))
	}

	row := result.InlineKeyboard[0]

	if len(row) != 1 {
		t.Errorf("len(keyboard.Parse()[0]) = %d; want 1", len(result.InlineKeyboard[0]))
	}

	button := row[0]

	if button.Text != "text" {
		t.Errorf("keyboard.Parse().Text = %s; want text", button.Text)
	}

	if button.CallbackData != "callback" {
		t.Errorf("keyboard.Parse().CallbackData = %s; want callback", button.CallbackData)
	}
}
