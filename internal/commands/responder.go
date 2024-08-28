package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TODO

type Responder interface {
	Respond(template string, a ...any) error
	RespondWithKeyboard(keyboard Keyboard, template string, a ...any) error
}

type TelegramResponder struct {
	Bot    tgbotapi.BotAPI
	Update tgbotapi.Update
}

func (r TelegramResponder) Respond(template string, a ...any) error {
	text := fmt.Sprintf(template, a...)

	var msg tgbotapi.Chattable

	if r.Update.Message != nil {
		msg = tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: r.Update.Message.Chat.ID,
			},
			Text:      text,
			ParseMode: tgbotapi.ModeMarkdown,
		}
	} else if r.Update.CallbackQuery != nil {
		msg = tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:    r.Update.CallbackQuery.Message.Chat.ID,
				MessageID: r.Update.CallbackQuery.Message.MessageID,
			},
			Text:      text,
			ParseMode: tgbotapi.ModeMarkdown,
		}
	}

	_, err := r.Bot.Send(msg)
	return err
}

func (r TelegramResponder) RespondWithKeyboard(keyboard Keyboard, template string, a ...any) error {
	text := fmt.Sprintf(template, a...)
	markup := ParseTelegramKeyboard(keyboard)

	var msg tgbotapi.Chattable

	if r.Update.Message != nil {
		msg = tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:      r.Update.Message.Chat.ID,
				ReplyMarkup: markup,
			},
			Text:      text,
			ParseMode: tgbotapi.ModeMarkdown,
		}
	} else if r.Update.CallbackQuery != nil {
		msg = tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:      r.Update.CallbackQuery.Message.Chat.ID,
				MessageID:   r.Update.CallbackQuery.Message.MessageID,
				ReplyMarkup: &markup,
			},
			Text:      text,
			ParseMode: tgbotapi.ModeMarkdown,
		}
	}

	_, err := r.Bot.Send(msg)
	return err
}

type DiscordResponder struct {
	Interaction discordgo.Interaction
	Session     *discordgo.Session
}

func (r DiscordResponder) Respond(template string, a ...any) error {
	text := fmt.Sprintf(template, a...)

	return r.Session.InteractionRespond(&r.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: text,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func (r DiscordResponder) RespondWithKeyboard(keyboard Keyboard, template string, a ...any) error {
	text := fmt.Sprintf(template, a...)

	return r.Session.InteractionRespond(&r.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    text,
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: ParseDiscordKeyboard(keyboard),
		},
	})
}
