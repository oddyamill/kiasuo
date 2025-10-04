package commands

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Responder interface {
	Write(template string, a ...any) Responder
	Respond() error
	RespondWithKeyboard(keyboard Keyboard) error
}

type TelegramResponder struct {
	Builder  strings.Builder
	Keyboard Keyboard
	Bot      tgbotapi.BotAPI
	Update   tgbotapi.Update
}

func (r *TelegramResponder) Write(template string, a ...any) Responder {
	text := fmt.Sprintf(template, a...)
	r.Builder.WriteString(text)
	return r
}

func (r *TelegramResponder) Respond() error {
	var msg tgbotapi.Chattable
	markup := ParseTelegramKeyboard(r.Keyboard)

	if r.Update.Message != nil {
		msg = tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:      r.Update.Message.Chat.ID,
				ReplyMarkup: markup,
			},
			Text:      r.Builder.String(),
			ParseMode: tgbotapi.ModeMarkdown,
		}
	} else if r.Update.CallbackQuery != nil {
		msg = tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:      r.Update.CallbackQuery.Message.Chat.ID,
				MessageID:   r.Update.CallbackQuery.Message.MessageID,
				ReplyMarkup: markup,
			},
			Text:      r.Builder.String(),
			ParseMode: tgbotapi.ModeMarkdown,
		}
	}

	_, err := r.Bot.Send(msg)
	return err
}

func (r *TelegramResponder) RespondWithKeyboard(keyboard Keyboard) error {
	r.Keyboard = keyboard
	return r.Respond()
}
