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

	template := formatter.Title("Панель управления") +
		formatter.Item("Telegram: %s") +
		formatter.Item("Discord: %s") +
		formatter.Item("Статус: %s")

	keyboard := Keyboard{
		KeyboardRow{
			KeyboardButton{
				Text:     helpers.If(user.State == users.Blacklisted, "Разблокировать", "Заблокировать"),
				Callback: AdminCommandName + ":blacklist:" + strconv.Itoa(user.ID),
			},
		},
	}

	return responder.RespondWithKeyboard(
		keyboard,
		template,
		formatter.Code(strconv.FormatInt(user.TelegramID, 10)),
		formatter.Code(helpers.If(user.DiscordID == "", "не указан", user.DiscordID)),
		formatter.Code(state),
	)
})

var AdminCallback = Callback(func(context Context, responder Responder, formatter helpers.Formatter, data []string) error {
	switch data[1] {
	case "blacklist":
		isBlacklisted := context.User.State == users.Blacklisted
		context.User.UpdateState(helpers.If(isBlacklisted, users.Ready, users.Blacklisted))
	default:
		return responder.Respond("Неизвестная команда. Меню устарело?")
	}

	return AdminCommand(context, responder, formatter)
})
