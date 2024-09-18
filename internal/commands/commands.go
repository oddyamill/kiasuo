package commands

import (
	"github.com/bwmarrin/discordgo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kiasuo/bot/internal/helpers"
	"log"
)

type Command func(Context, Responder, helpers.Formatter) error

var commandMap = map[string]Command{
	AdminCommandName: AdminCommand,
	"start":          StartCommand,
	"stop":           StopCommand,
	"settings":       SettingsCommand,
	"schedule":       ScheduleCommand,
	"marks":          MarksCommand,
	"classmates":     ClassmatesCommand,
	"teachers":       TeachersCommand,
}

type commandConfig struct {
	Name         string
	Description  string
	TelegramOnly bool
}

var publicCommands = []commandConfig{
	{
		Name:         "start",
		Description:  "Начать работу",
		TelegramOnly: true,
	},
	{
		Name:         "stop",
		Description:  "Завершить работу с ботом",
		TelegramOnly: true,
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

func ParseDiscordCommands() []*discordgo.ApplicationCommand {
	commands := make([]*discordgo.ApplicationCommand, 0)

	for _, config := range publicCommands {
		if config.TelegramOnly {
			continue
		}

		commands = append(commands, &discordgo.ApplicationCommand{
			Name:        config.Name,
			Description: config.Description,
		})
	}

	return commands
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
	AdminCommandName: AdminCallback,
	"stop":           StopCallback,
	"settings":       SettingsCallback,
	"schedule":       ScheduleCallback,
	"marks":          MarksCallback,
}

func HandleCommand(context Context, responder Responder, formatter helpers.Formatter) {
	command := commandMap[context.Command]

	if command == nil {
		log.Println("Someone tried to execute unknown command", context.Command)
		return
	}

	log.Println("Handling command", context.Command)
	handleError(responder, command(context, responder, formatter))
}

func HandleCallback(context Context, responder Responder, formatter helpers.Formatter, data []string) {
	callback := callbackMap[context.Command]

	if callback == nil {
		log.Println("Someone tried to execute unknown callback", context.Command)
		return
	}

	log.Println("Handling callback", context.Command)
	handleError(responder, callback(context, responder, formatter, data))
}

func handleError(responder Responder, err error) {
	if err == nil {
		//err = responder.Respond()

		//if err == nil {
		//	return
		//}

		return
	}

	log.Printf("Error: %v", err)
}
