package commands

import (
	"github.com/kiasuo/bot/internal/helpers"
	"github.com/kiasuo/bot/internal/users_sql"
	"strconv"
)

const AdminCommandName string = "-internal-admin"

func AdminCommand(context Context, responder Responder, formatter Formatter) {
	user := context.User

	state := map[users_sql.UserState]string{
		users_sql.Unknown:     "неизвестно",
		users_sql.Ready:       "готов",
		users_sql.Pending:     "ожидает",
		users_sql.Blacklisted: "заблокирован",
	}[user.State]

	template := formatter.Title("Панель управления") +
		formatter.Item("Telegram: %s") +
		formatter.Item("Discord: %s") +
		formatter.Item("Статус: %s")

	keyboard := Keyboard{
		KeyboardRow{
			KeyboardButton{
				Text:     helpers.If(user.State == users_sql.Blacklisted, "Разблокировать", "Заблокировать"),
				Callback: AdminCommandName + ":blacklist:" + string(rune(user.ID)),
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
		isBlacklisted := context.User.State == users_sql.Blacklisted
		context.User.UpdateState(helpers.If(isBlacklisted, users_sql.Ready, users_sql.Blacklisted))
	default:
		responder.Respond("Неизвестная команда. Меню устарело?")
		return
	}

	AdminCommand(context, responder, formatter)
}
