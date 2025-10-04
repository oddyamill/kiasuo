package main

import (
	"log/slog"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kiasuo/bot/internal/client"
	"github.com/kiasuo/bot/internal/commands"
	"github.com/kiasuo/bot/internal/helpers"
	"github.com/kiasuo/bot/internal/users"
	"github.com/kiasuo/bot/internal/version"
)

const AdminID int64 = 6135991898

func main() {
	token := helpers.GetEnv("TELEGRAM")
	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		panic(err)
	}

	slog.Info("authorized on account", "username", "hi")

	if _, err = bot.Request(commands.ParseTelegramCommands()); err != nil {
		panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		switch {
		case update.Message != nil:
			handleMessage(bot, update)
		case update.CallbackQuery != nil:
			handleCallbackQuery(bot, update)
		case update.MyChatMember != nil:
			handleMyChatMember(update)
		}
	}
}

func handleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var (
		user      *users.User
		command   string
		arguments string
	)

	resp := commands.TelegramResponder{
		Bot:    *bot,
		Update: update,
	}

	if update.Message.ForwardFrom != nil {
		if update.Message.From.ID != AdminID {
			return
		}

		user = users.GetByTelegramID(update.Message.ForwardFrom.ID)

		if user == nil {
			_ = resp.Write("Пользователь не зарегистрирован").Respond()
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
			// TODO:
			return
		case users.Ready:
			break
		case users.Pending:
			if commands.IsSystemCommand(command) {
				break
			}
			_ = resp.Write("Токен обнови.").Respond()
			return
		case users.Blacklisted:
			_ = resp.Write("Ты заблокирован. Клоун.").Respond()
			return
		default:
			return
		}

		user = users.GetByID(id)
		arguments = update.Message.CommandArguments()
	} else {
		return
	}

	ctx := commands.NewContext(command, arguments, *user)

	formatter := helpers.TelegramFormatter{}

	commands.HandleCommand(ctx, &resp, &formatter)
}

func handleCallbackQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	data := strings.Split(update.CallbackQuery.Data, ":")

	if len(data) < 2 {
		return
	}

	var user *users.User

	resp := commands.TelegramResponder{
		Bot:    *bot,
		Update: update,
	}

	if data[0] != version.Version {
		_ = resp.Write("Меню устарело. Используйте команду повторно.").Respond()
	}

	if data[1] == commands.AdminCommandName {
		user = users.GetByID(data[3])

		if user == nil {
			_ = resp.Write("Пользователь не зарегистрирован").Respond()
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

	ctx := commands.NewContext(data[1], "", *user)

	formatter := helpers.TelegramFormatter{}

	commands.HandleCallback(ctx, &resp, &formatter, data[2:])
}

func handleMyChatMember(update tgbotapi.Update) {
	if !update.MyChatMember.Chat.IsPrivate() || !update.MyChatMember.NewChatMember.WasKicked() {
		return
	}

	user := users.GetByTelegramID(update.MyChatMember.NewChatMember.User.ID)

	if user == nil {
		return
	}

	c := client.New(user)

	if user.State == users.Ready {
		_ = c.RevokeToken()
	}

	if user.Cache {
		_ = c.PurgeCache()
	}

	user.Delete()
}
