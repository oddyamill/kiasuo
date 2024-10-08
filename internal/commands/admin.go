package commands

import (
	"github.com/kiasuo/bot/internal/helpers"
	"github.com/kiasuo/bot/internal/users"
	"strconv"
)

const AdminCommandName string = "-internal-admin"

var AdminCommand = Command(func(context Context, responder Responder, formatter helpers.Formatter) error {
	user := context.User

	state := map[users.UserState]string{
		users.Unknown:     "неизвестно",
		users.Ready:       "готов",
		users.Pending:     "ожидает",
		users.Blacklisted: "заблокирован",
	}[user.State]

	text := formatter.Title("Панель управления") +
		formatter.Item("Telegram: "+strconv.FormatInt(user.TelegramID, 10)) +
		formatter.Item("Discord: "+helpers.If(user.DiscordID.Valid, user.DiscordID.String, "не указан")) +
		formatter.Item("Статус: "+state)

	keyboard := Keyboard{
		KeyboardRow{
			KeyboardButton{
				Text:     helpers.If(user.State == users.Blacklisted, "Разблокировать", "Заблокировать"),
				Callback: AdminCommandName + ":blacklist:" + strconv.Itoa(user.ID),
			},
		},
	}

	return responder.Write(text).RespondWithKeyboard(keyboard)
})

var AdminCallback = Callback(func(context Context, responder Responder, formatter helpers.Formatter, data []string) error {
	switch data[1] {
	case "blacklist":
		isBlacklisted := context.User.State == users.Blacklisted
		context.User.UpdateState(helpers.If(isBlacklisted, users.Ready, users.Blacklisted))
	default:
		return responder.Write("Неизвестная команда. Меню устарело?").Respond()
	}

	return AdminCommand(context, responder, formatter)
})
