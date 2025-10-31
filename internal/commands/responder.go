package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Responder struct {
	builder  strings.Builder
	keyboard Keyboard
	Bot      *bot.Bot
	Update   models.Update
}

func (r *Responder) Write(template string, a ...any) *Responder {
	text := fmt.Sprintf(template, a...)
	r.builder.WriteString(text)
	return r
}

func (r *Responder) Respond() error {
	return r.respond()
}

func (r *Responder) RespondWithKeyboard(keyboard Keyboard) error {
	r.keyboard = keyboard
	return r.respond()
}

func (r *Responder) respond() error {
	switch {
	case r.Update.Message != nil:
		return r.sendMessage()
	case r.Update.CallbackQuery != nil:
		return r.editMessage()
	default:
		return errors.New("unsupported Update type")
	}
}

func (r *Responder) sendMessage() error {
	msg := &bot.SendMessageParams{
		ChatID:      r.Update.Message.Chat.ID,
		Text:        r.builder.String(),
		ReplyMarkup: r.keyboard.Parse(),
		ParseMode:   models.ParseModeMarkdownV1,
	}

	_, err := r.Bot.SendMessage(context.Background(), msg)
	return err
}

func (r *Responder) editMessage() error {
	msg := &bot.EditMessageTextParams{
		ChatID:      r.Update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:   r.Update.CallbackQuery.Message.Message.ID,
		Text:        r.builder.String(),
		ReplyMarkup: r.keyboard.Parse(),
		ParseMode:   models.ParseModeMarkdownV1,
	}

	_, err := r.Bot.EditMessageText(context.Background(), msg)
	return err
}
