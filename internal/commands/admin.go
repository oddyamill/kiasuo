package commands

import (
	"strconv"

	"github.com/kiasuo/bot/internal/helpers"
	"github.com/kiasuo/bot/internal/users"
	"github.com/kiasuo/bot/internal/version"
)

const AdminCommandName string = "-internal-admin"

var AdminCommand = Command(func(ctx Context, resp Responder, formatter helpers.Formatter) error {
	user := ctx.User

	text := formatter.Title("Панель управления") +
		formatter.Item("Telegram: "+strconv.FormatInt(user.TelegramID, 10)) +
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

	return resp.Write(text).RespondWithKeyboard(keyboard)
})

var AdminCallback = Callback(func(ctx Context, resp Responder, formatter helpers.Formatter, data []string) error {
	switch data[0] {
	case "blacklist":
		isBlacklisted := ctx.User.State == users.Blacklisted
		ctx.User.UpdateState(helpers.If(isBlacklisted, users.Ready, users.Blacklisted))
	default:
		return resp.Write("Неизвестная команда. Меню устарело?").Respond()
	}

	return AdminCommand(ctx, resp, formatter)
})
