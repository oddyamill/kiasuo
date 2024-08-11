package commands

import (
	"github.com/kiasuo/bot/helpers"
	"github.com/kiasuo/bot/users"
	"strconv"
)

const AdminCommandName string = "teachers"

func AdminCommand(context Context, responder Responder, formatter Formatter) {
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
				Callback: AdminCommandName + ":blacklist:" + user.ID.Hex(),
			},
		},
	}

	responder.RespondWithKeyboard(
		keyboard,
		template,
		formatter.Code(strconv.FormatInt(user.TelegramID, 10)),
		formatter.Code(helpers.If(user.DiscordID == "", "не указан", user.DiscordID)),
		formatter.Code(state),
	)
}

func AdminCallback(context Context, responder Responder, formatter Formatter, data []string) {
	switch data[1] {
	case "blacklist":
		isBlacklisted := context.User.State == users.Blacklisted
		users.UpdateState(context.User, helpers.If(isBlacklisted, users.Ready, users.Blacklisted))
	default:
		responder.Respond("Неизвестная команда. Меню устарело?")
		return
	}

	AdminCommand(context, responder, formatter)
}
