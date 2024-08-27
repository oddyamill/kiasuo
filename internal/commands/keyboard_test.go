package commands

import (
	"github.com/bwmarrin/discordgo"
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

func TestParseDiscordKeyboard(t *testing.T) {
	keyboard := Keyboard{
		{
			{
				Text:     "text",
				Callback: "callback",
			},
		},
	}

	result := ParseDiscordKeyboard(keyboard)

	if len(result) != 1 {
		t.Errorf("len(ParseDiscordKeyboard()) = %d; want 1", len(result))
	}

	row := result[0].(discordgo.ActionsRow)

	if len(row.Components) != 1 {
		t.Errorf("len(ParseDiscordKeyboard()[0].Components) = %d; want 1", len(result[0].(discordgo.ActionsRow).Components))
	}

	button := row.Components[0].(discordgo.Button)

	if button.Label != "text" {
		t.Errorf("ParseDiscordKeyboard() = %s; want text", button.Label)
	}

	if button.CustomID != "callback" {
		t.Errorf("ParseDiscordKeyboard() = %s; want callback", button.CustomID)
	}
}
