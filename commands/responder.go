package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Responder interface {
	Respond(template string, a ...any)
}

type TelegramResponder struct {
	Bot    tgbotapi.BotAPI
	Update tgbotapi.Update
}

func (r *TelegramResponder) Respond(template string, a ...any) {
	text := fmt.Sprintf(template, a...)

	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: r.Update.Message.Chat.ID,
		},
		Text:      text,
		ParseMode: tgbotapi.ModeMarkdown,
	}

	r.Bot.Send(msg)
}

type DiscordResponder struct {
	Interaction discordgo.Interaction
	Session     *discordgo.Session
}

func (r *DiscordResponder) Respond(template string, a ...any) {
	text := fmt.Sprintf(template, a...)

	r.Session.InteractionRespond(&r.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: text,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
