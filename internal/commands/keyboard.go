package commands

import (
	"github.com/go-telegram/bot/models"
	"github.com/kiasuo/bot/internal/version"
)

type Keyboard []KeyboardRow

func (k Keyboard) Parse() models.ReplyMarkup {
	if len(k) == 0 {
		return nil
	}

	result := &models.InlineKeyboardMarkup{}

	for _, row := range k {
		var buttons []models.InlineKeyboardButton

		for _, button := range row {
			inlineKeyboardButton := models.InlineKeyboardButton{
				Text: button.Text,
			}

			if button.Callback != "" {
				inlineKeyboardButton.CallbackData = button.Callback
			}

			if button.WebAppURL != "" {
				inlineKeyboardButton.WebApp = &models.WebAppInfo{
					URL: button.WebAppURL,
				}
			}

			buttons = append(buttons, inlineKeyboardButton)
		}

		result.InlineKeyboard = append(result.InlineKeyboard, buttons)
	}

	return result
}

type KeyboardRow []KeyboardButton

type KeyboardButton struct {
	Text      string
	Callback  string
	WebAppURL string
}

func NewCallbackButton(text, callback string) KeyboardButton {
	return KeyboardButton{Text: text, Callback: version.Version + ":" + callback}
}

func NewWebappButton(text, url string) KeyboardButton {
	return KeyboardButton{Text: text, WebAppURL: url}
}
