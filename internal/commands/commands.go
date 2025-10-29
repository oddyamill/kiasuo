package commands

import (
	"errors"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kiasuo/bot/internal/client"
	"github.com/kiasuo/bot/internal/database"
	"github.com/kiasuo/bot/internal/helpers"
)

type Command func(Context, Responder, helpers.Formatter) error

var commandMap = map[string]Command{
	AdminCommandName: AdminCommand,
	StartCommandName: StartCommand,
	"settings":       SettingsCommand,
	"schedule":       ScheduleCommand,
	"marks":          MarksCommand,
	"classmates":     ClassmatesCommand,
	"teachers":       TeachersCommand,
}

func IsSystemCommand(command string) bool {
	return command == StartCommandName || command == "settings"
}

type commandConfig struct {
	Name        string
	Description string
}

var publicCommands = []commandConfig{
	{
		Name:        StartCommandName,
		Description: "Начать",
	},
	{
		Name:        "settings",
		Description: "Настройки",
	},
	{
		Name:        "schedule",
		Description: "Расписание",
	},
	{
		Name:        "marks",
		Description: "Оценки",
	},
	{
		Name:        "classmates",
		Description: "Список одноклассников",
	},
	{
		Name:        "teachers",
		Description: "Список учителей",
	},
	{
		Name:        "stop",
		Description: "Остановить",
	},
}

func ParseTelegramCommands() tgbotapi.SetMyCommandsConfig {
	commands := make([]tgbotapi.BotCommand, 0)

	for _, config := range publicCommands {
		commands = append(commands, tgbotapi.BotCommand{
			Command:     config.Name,
			Description: config.Description,
		})
	}

	return tgbotapi.NewSetMyCommands(commands...)
}

type Callback func(Context, Responder, helpers.Formatter, []string) error

var callbackMap = map[string]Callback{
	"settings": SettingsCallback,
	"schedule": ScheduleCallback,
	"marks":    MarksCallback,
}

func HandleCommand(ctx Context, resp Responder, fmt helpers.Formatter) {
	command := commandMap[ctx.Command]

	if command == nil {
		slog.Warn("someone tried to execute unknown command", "command", ctx.Command)
		return
	}

	slog.Info("handling command", "command", ctx.Command)

	if err := command(ctx, resp, fmt); err != nil {
		handleError(ctx, resp, err)
	}
}

func HandleCallback(ctx Context, resp Responder, fmt helpers.Formatter, data []string) {
	callback := callbackMap[ctx.Command]

	if callback == nil {
		slog.Warn("someone tried to execute unknown callback", "callback", ctx.Command)
		return
	}

	slog.Info("handling callback", "callback", ctx.Command)

	if err := callback(ctx, resp, fmt, data); err != nil {
		handleError(ctx, resp, err)
	}
}

func handleError(ctx Context, resp Responder, err error) {
	if !errors.Is(err, client.ErrExpiredToken) {
		slog.Error("error while handling command", "command", ctx.Command, "error", err)
		return
	}

	err = ctx.User.SetState(ctx.Context(), database.UserStatePending)

	if err != nil {
		slog.Warn("failed to set user state", "error", err)
		return
	}

	_ = resp.Write("Ваш токен истек. Используйте /start для получении информации").Respond()
}
