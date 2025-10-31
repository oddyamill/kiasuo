package main

import (
	"context"
	"log"
	"log/slog"
	"strings"

	"github.com/go-telegram/bot/models"
	"github.com/kiasuo/bot/internal/client"
	"github.com/kiasuo/bot/internal/commands"
	"github.com/kiasuo/bot/internal/database"
	"github.com/kiasuo/bot/internal/helpers"
	"github.com/kiasuo/bot/internal/version"
)

func handleBot(app *App) {
	me, err := app.bot.GetMe(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	slog.Info("authorized on account", "username", me.Username)

	if _, err := app.bot.SetMyCommands(context.Background(), commands.ParseTelegramCommands()); err != nil {
		log.Fatal(err)
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

func resolveCommand(message *models.Message) (string, string) {
	for _, e := range message.Entities {
		if e.Type == models.MessageEntityTypeBotCommand {
			command := message.Text[e.Offset+1 : e.Offset+e.Length]

			if i := strings.Index(command, "@"); i != -1 {
				command = command[:i]
			}

			if len(message.Text) < e.Offset+e.Length+1 {
				return command, ""
			}

			return command, message.Text[e.Offset+e.Length+1:]
		}
	}

	return "", ""
}

func handleMessage(app *App, update models.Update) {
	message := update.Message

	command, arguments := resolveCommand(message)

	if command == "" {
		return
	}

	resp := &commands.Responder{
		Bot:    app.bot,
		Update: update,
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
		if user.StudentID == 0 && command != commands.SettingsCommandName {
			//todo: auto!!!
			_ = resp.Write("Необходимо выбрать ученика. Используйте /settings.").Respond()
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

	ctx := commands.NewContext(command, arguments, *user)

	formatter := helpers.Formatter{}

	commands.HandleCommand(ctx, resp, formatter)
}

func handleCallbackQuery(app *App, update models.Update) {
	callbackQuery := update.CallbackQuery
	data := strings.Split(callbackQuery.Data, ":")

	if len(data) < 2 {
		return
	}

	resp := &commands.Responder{
		Bot:    app.bot,
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
		//TODO:?
		if data[1] != commands.StartCommandName {
			return
		}
	default:
		return
	}

	ctx := commands.NewContext(data[1], "", *user)

	formatter := helpers.Formatter{}

	commands.HandleCallback(ctx, resp, formatter, data[2:])
}

func handleMyChatMember(app *App, chatMember *models.ChatMemberUpdated) {
	if chatMember.Chat.Type != models.ChatTypePrivate {
		return
	}

	if chatMember.NewChatMember.Member.Status != models.ChatMemberTypeBanned {
		return
	}

	user, err := app.db.GetUser(context.Background(), chatMember.NewChatMember.Member.User.ID)

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
