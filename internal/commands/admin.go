package commands

import (
	"strconv"

	"github.com/kiasuo/bot/internal/helpers"
	"github.com/kiasuo/bot/internal/users"
	"github.com/kiasuo/bot/internal/version"
)

const AdminCommandName string = "-internal-admin"

var AdminCommand = Command(func(context Context, responder Responder, formatter helpers.Formatter) error {
	user := context.User

	text := formatter.Title("Панель управления") +
		formatter.Item("Telegram: "+strconv.FormatInt(user.TelegramID, 10)) +
		formatter.Item("Discord: "+helpers.If(user.DiscordID.Valid, user.DiscordID.String, "не указан")) +
		formatter.Item("Статус: "+user.State.String()) +
		formatter.Item("Версия: "+version.Version)

	keyboard := Keyboard{
		KeyboardRow{
			NewKeyboardButton(
				helpers.If(user.State == users.Blacklisted, "Разблокировать", "Заблокировать"),
				AdminCommandName+":blacklist:"+strconv.Itoa(user.ID),
			),
		},
	}

	return responder.Write(text).RespondWithKeyboard(keyboard)
})

var AdminCallback = Callback(func(context Context, responder Responder, formatter helpers.Formatter, data []string) error {
	switch data[0] {
	case "blacklist":
		isBlacklisted := context.User.State == users.Blacklisted
		context.User.UpdateState(helpers.If(isBlacklisted, users.Ready, users.Blacklisted))
	default:
		return responder.Write("Неизвестная команда. Меню устарело?").Respond()
	}

	return AdminCommand(context, responder, formatter)
})
