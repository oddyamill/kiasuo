package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
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

type DiscordResponder struct {
	Keyboard    Keyboard
	Builder     strings.Builder
	Interaction discordgo.Interaction
	Session     *discordgo.Session
}

func (r *DiscordResponder) Write(template string, a ...any) Responder {
	text := fmt.Sprintf(template, a...)
	r.Builder.WriteString(text)
	return r
}

func (r *DiscordResponder) Respond() error {
	content := r.Builder.String()
	components := ParseDiscordKeyboard(r.Keyboard)

	_, err := r.Session.InteractionResponseEdit(&r.Interaction, &discordgo.WebhookEdit{
		Content:    &content,
		Components: &components,
	})

	return err
}

func (r *DiscordResponder) RespondWithKeyboard(keyboard Keyboard) error {
	r.Keyboard = keyboard
	return r.Respond()
}

func (r *DiscordResponder) RespondWithDefer() error {
	if r.Interaction.Type == discordgo.InteractionApplicationCommand {
		return r.Session.InteractionRespond(&r.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
	} else {
		for _, component := range r.Interaction.Message.Components {
			row := component.(*discordgo.ActionsRow)

			for _, button := range row.Components {
				btn := button.(*discordgo.Button)
				btn.Disabled = true
			}
		}

		return r.Session.InteractionRespond(&r.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content:    r.Interaction.Message.Content,
				Components: r.Interaction.Message.Components,
			},
		})
	}
}
