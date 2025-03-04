package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kiasuo/bot/internal/commands"
	"github.com/kiasuo/bot/internal/helpers"
	"github.com/kiasuo/bot/internal/users"
	"github.com/kiasuo/bot/internal/version"
	"log"
	"strings"
)

const AdminID int64 = 6135991898

var bot tgbotapi.BotAPI

func init() {
	token := helpers.GetEnv("TELEGRAM")
	botApi, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		panic(err)
	}

	bot = *botApi
	log.Println("Authorized on account", bot.Self.UserName)

	if _, err = bot.Request(commands.ParseTelegramCommands()); err != nil {
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
		user      *users.User
		command   string
		arguments string
	)

	responder := commands.TelegramResponder{
		Bot:    bot,
		Update: update,
	}

	if update.Message.ForwardFrom != nil {
		if update.Message.From.ID != AdminID {
			return
		}

		user = users.GetByTelegramID(update.Message.ForwardFrom.ID)

		if user == nil {
			_ = responder.Write("Пользователь не зарегистрирован").Respond()
			return
		}

		command = commands.AdminCommandName
	} else if update.Message.IsCommand() {
		command = update.Message.Command()

		if command == "" {
			return
		}

		id, state := users.GetPartialByTelegramID(update.Message.From.ID)

		switch state {
		case users.Unknown:
			if command == commands.StartCommandName {
				users.CreateWithTelegramID(update.Message.From.ID)
			}
			return
		case users.Ready:
			break
		case users.Pending:
			if commands.IsSystemCommand(command) {
				break
			}
			_ = responder.Write("Токен обнови.").Respond()
			return
		case users.Blacklisted:
			_ = responder.Write("Ты заблокирован. Клоун.").Respond()
			return
		default:
			return
		}

		user = users.GetByID(id)
		arguments = update.Message.CommandArguments()
	} else {
		return
	}

	context := commands.Context{
		Command:   command,
		Arguments: arguments,
		User:      *user,
	}

	formatter := helpers.TelegramFormatter{}

	commands.HandleCommand(context, &responder, &formatter)
}

func handleCallbackQuery(update tgbotapi.Update) {
	data := strings.Split(update.CallbackQuery.Data, ":")

	if len(data) < 2 {
		return
	}

	var user *users.User

	responder := commands.TelegramResponder{
		Bot:    bot,
		Update: update,
	}

	if data[0] != version.Version {
		_ = responder.Write("Меню устарело. Используйте команду повторно.").Respond()
	}

	if data[1] == commands.AdminCommandName {
		user = users.GetByID(data[3])

		if user == nil {
			_ = responder.Write("Пользователь не зарегистрирован").Respond()
			return
		}
	} else {
		user = users.GetByTelegramID(update.CallbackQuery.From.ID)

		if user == nil {
			return
		}

		switch user.State {
		case users.Ready:
			break
		case users.Pending:
			if commands.IsSystemCommand(data[1]) {
				break
			}
		default:
			return
		}
	}

	context := commands.Context{
		Command: data[1],
		User:    *user,
	}

	formatter := helpers.TelegramFormatter{}

	commands.HandleCallback(context, &responder, &formatter, data[2:])
}
