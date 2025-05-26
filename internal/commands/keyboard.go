package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kiasuo/bot/internal/version"
)

type Keyboard []KeyboardRow

type KeyboardRow []KeyboardButton

type KeyboardButton struct {
	Text     string
	Callback string
}

func NewKeyboardButton(text, callback string) KeyboardButton {
	return KeyboardButton{Text: text, Callback: version.Version + ":" + callback}
}

func ParseTelegramKeyboard(keyboard Keyboard) *tgbotapi.InlineKeyboardMarkup {
	if len(keyboard) == 0 {
		return nil
	}

	result := &tgbotapi.InlineKeyboardMarkup{}

	for _, row := range keyboard {
		var buttons []tgbotapi.InlineKeyboardButton

		for _, button := range row {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(button.Text, button.Callback))
		}

		result.InlineKeyboard = append(result.InlineKeyboard, buttons)
	}

	return result
}
