package commands

import (
	"github.com/bwmarrin/discordgo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Keyboard []KeyboardRow

type KeyboardRow []KeyboardButton

type KeyboardButton struct {
	Text     string
	Callback string
}

func ParseTelegramKeyboard(keyboard Keyboard) (tgKeyboard tgbotapi.InlineKeyboardMarkup) {
	for _, row := range keyboard {
		var buttons []tgbotapi.InlineKeyboardButton

		for _, button := range row {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(button.Text, button.Callback))
		}

		tgKeyboard.InlineKeyboard = append(tgKeyboard.InlineKeyboard, buttons)
	}

	return
}

func ParseDiscordKeyboard(keyboard Keyboard) (dgKeyboard []discordgo.MessageComponent) {
	for _, row := range keyboard {
		var components []discordgo.MessageComponent

		for _, button := range row {
			components = append(components, discordgo.Button{
				Label:    button.Text,
				CustomID: button.Callback,
			})
		}

		dgKeyboard = append(dgKeyboard, discordgo.ActionsRow{
			Components: components,
		})
	}

	return
}
