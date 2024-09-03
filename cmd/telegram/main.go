package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kiasuo/bot/internal/commands"
	"github.com/kiasuo/bot/internal/helpers"
	"github.com/kiasuo/bot/internal/users"
	"log"
	"strings"
)

const AdminId int64 = 6135991898

var bot tgbotapi.BotAPI

func init() {
	token := helpers.GetEnv("TELEGRAM")
	botApi, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		panic(err)
	}

	bot = *botApi
	log.Println("Authorized on account", bot.Self.UserName)

	_, err = bot.Request(commands.ParseTelegramCommands())

	if err != nil {
		panic(err)
	}
}

func main() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handleMessage(update)
		} else if update.CallbackQuery != nil {
			handleCallbackQuery(update)
		}
	}
}

func handleMessage(update tgbotapi.Update) {
	var (
		user    *users.User
		command string
	)

	responder := commands.TelegramResponder{
		Bot:    bot,
		Update: update,
	}

	if update.Message.ForwardFrom != nil {
		if update.Message.From.ID != AdminId {
			return
		}

		user = users.GetByTelegramID(update.Message.ForwardFrom.ID)

		if user == nil {
			responder.Respond("Пользователь не зарегистрирован")
			return
		}

		command = commands.AdminCommandName
	} else if update.Message.IsCommand() {
		user = users.GetByTelegramID(update.Message.From.ID)

		if user == nil {
			responder.Respond("Ты кто такой? Уйди.")
			return
		}

		if user.State != users.Ready {
			responder.Respond("Пошел отсюда.")
			return
		}

		command = update.Message.Command()
	}

	if command == "" {
		return
	}

	context := commands.Context{
		Command: command,
		User:    *user,
	}

	formatter := commands.TelegramFormatter{}

	commands.HandleCommand(context, &responder, &formatter)
}

func handleCallbackQuery(update tgbotapi.Update) {
	data := strings.Split(update.CallbackQuery.Data, ":")

	if len(data) < 2 {
		return
	}

	var user *users.User

	if data[0] == commands.AdminCommandName {
		user = users.GetByID(data[2])
	} else {
		user = users.GetByTelegramID(update.CallbackQuery.From.ID)
	}

	context := commands.Context{
		Command: data[0],
		User:    *user,
	}

	responder := commands.TelegramResponder{
		Bot:    bot,
		Update: update,
	}

	formatter := commands.TelegramFormatter{}

	commands.HandleCallback(context, &responder, &formatter, data)
}
