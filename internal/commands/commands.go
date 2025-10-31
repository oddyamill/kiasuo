package commands

import (
	"errors"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/kiasuo/bot/internal/client"
	"github.com/kiasuo/bot/internal/database"
	"github.com/kiasuo/bot/internal/helpers"
)

type Command func(Context, *Responder, helpers.Formatter) error

var commandMap = map[string]Command{
	AdminCommandName:    AdminCommand,
	StartCommandName:    StartCommand,
	SettingsCommandName: SettingsCommand,
	"schedule":          ScheduleCommand,
	"marks":             MarksCommand,
	"classmates":        ClassmatesCommand,
	"teachers":          TeachersCommand,
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
}

func ParseTelegramCommands() *bot.SetMyCommandsParams {
	commands := make([]models.BotCommand, 0)

	for _, config := range publicCommands {
		commands = append(commands, models.BotCommand{
			Command:     config.Name,
			Description: config.Description,
		})
	}

	return &bot.SetMyCommandsParams{Commands: commands}
}

type Callback func(Context, *Responder, helpers.Formatter, []string) error

var callbackMap = map[string]Callback{
	"settings": SettingsCallback,
	"schedule": ScheduleCallback,
	"marks":    MarksCallback,
}

func HandleCommand(ctx Context, resp *Responder, fmt helpers.Formatter) {
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

func HandleCallback(ctx Context, resp *Responder, fmt helpers.Formatter, data []string) {
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

func handleError(ctx Context, resp *Responder, err error) {
	switch {
	case errors.Is(err, client.ErrExpiredToken):
		err = ctx.User.SetState(ctx.Context(), database.UserStatePending)

		if err != nil {
			slog.Warn("failed to set user state", "error", err)
		}

		_ = resp.Write("Ваш токен истек. Используйте /start для получении информации").Respond()
		break
	case errors.Is(err, client.ErrServerError):
		_ = resp.Write("Сервер КИАСУО не отвечает… Попробуйте позже…").Respond()
	default:
		slog.Error("error while handling command", "command", ctx.Command, "error", err)
	}
}
