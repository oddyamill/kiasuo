package main

import (
	"context"
	"log"
	"log/slog"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kiasuo/bot/internal/client"
	"github.com/kiasuo/bot/internal/commands"
	"github.com/kiasuo/bot/internal/database"
	"github.com/kiasuo/bot/internal/helpers"
	"github.com/kiasuo/bot/internal/version"
)

func handleBot(app *App) {
	slog.Info("authorized on account", "username", app.bot.Self.UserName)

	if _, err := app.bot.Request(commands.ParseTelegramCommands()); err != nil {
		log.Panic(err)
	}

	for update := range app.updates {
		switch {
		case update.Message != nil:
			handleMessage(app, update)
		case update.CallbackQuery != nil:
			handleCallbackQuery(app, update)
		case update.MyChatMember != nil:
			handleMyChatMember(app, update.MyChatMember)
		}
	}
}

func handleMessage(app *App, update tgbotapi.Update) {
	message := update.Message

	if !message.IsCommand() {
		return
	}

	resp := commands.TelegramResponder{
		Bot:    *app.bot,
		Update: update,
	}

	command := update.Message.Command()

	if command == "" {
		return
	}

	user, err := app.db.GetUser(context.Background(), message.Chat.ID)

	if err != nil {
		slog.Error(err.Error(), "event", "message")
		return
	}

	if user == nil {
		if command != commands.StartCommandName {
			return
		}

		user, err = app.db.NewUser(context.Background(), message.Chat.ID)

		if err != nil {
			slog.Error(err.Error(), "event", "message")
			return
		}
	}

	switch user.State {
	case database.UserStateReady:
		if user.StudentID == 0 && !commands.IsSystemCommand(command) {
			_ = resp.Write("Необходимо выбрать ученика. Используйте /settings").Respond()
			return
		}
	case database.UserStatePending:
		if command != commands.StartCommandName {
			_ = resp.Write("Токен обнови.").Respond()
			return
		}
	case database.UserStateBlacklisted:
		_ = resp.Write("Ты заблокирован. Клоун.").Respond()
		return
	default:
		break
	}

	arguments := update.Message.CommandArguments()

	ctx := commands.NewContext(command, arguments, *user)

	formatter := helpers.TelegramFormatter{}

	commands.HandleCommand(ctx, &resp, &formatter)
}

func handleCallbackQuery(app *App, update tgbotapi.Update) {
	callbackQuery := update.CallbackQuery
	data := strings.Split(callbackQuery.Data, ":")

	if len(data) < 2 {
		return
	}

	resp := commands.TelegramResponder{
		Bot:    *app.bot,
		Update: update,
	}

	if data[0] != version.Version {
		_ = resp.Write("Меню устарело. Используйте команду повторно.").Respond()
		return
	}

	user, err := app.db.GetUser(context.Background(), callbackQuery.From.ID)

	if err != nil {
		slog.Error(err.Error(), "event", "callbackQuery")
	}

	if user == nil {
		return
	}

	switch user.State {
	case database.UserStateReady:
		break
	case database.UserStatePending:
		if commands.IsSystemCommand(data[1]) {
			break
		}
	default:
		return
	}

	ctx := commands.NewContext(data[1], "", *user)

	formatter := helpers.TelegramFormatter{}

	commands.HandleCallback(ctx, &resp, &formatter, data[2:])
}

func handleMyChatMember(app *App, chatMember *tgbotapi.ChatMemberUpdated) {
	if !chatMember.Chat.IsPrivate() || !chatMember.NewChatMember.WasKicked() {
		return
	}

	user, err := app.db.GetUser(context.Background(), chatMember.NewChatMember.User.ID)

	if err != nil {
		slog.Error(err.Error(), "event", "myChatMember")
		return
	}

	if user == nil {
		return
	}

	c := client.New(user)

	if user.State == database.UserStateReady {
		_ = c.RevokeToken()
	}

	if user.HasFlag(database.UserFlagCache) {
		_ = c.PurgeCache()
	}

	if err = user.Delete(context.Background()); err != nil {
		slog.Error(err.Error(), "event", "myChatMember", "method", "User#Delete()")
	}
}
