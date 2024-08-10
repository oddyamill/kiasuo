package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kiasuo/bot/commands"
	"github.com/kiasuo/bot/users"
	"log"
	"os"
)

var bot tgbotapi.BotAPI

func init() {
	token, ok := os.LookupEnv("TELEGRAM")

	if !ok {
		panic("TELEGRAM not set")
	}

	botApi, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		panic(err)
	}

	bot = *botApi
	log.Println("Authorized on account", bot.Self.UserName)
}

func main() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		user := users.GetByTelegramID(update.Message.From.ID)

		responder := commands.TelegramResponder{
			Bot:    bot,
			Update: update,
		}

		if user == nil {
			responder.Respond("Ты кто такой? Cъебал.")
			continue
		}

		if user.State != users.Ready {
			responder.Respond("Пошел нахуй.")
			continue
		}

		context := commands.Context{
			Command: update.Message.Command(),
			User:    *user,
		}

		formatter := commands.TelegramFormatter{}

		commands.HandleCommand(context, &responder, &formatter)
	}
}
