package commands

import (
	"log"
	"strconv"

	"github.com/kiasuo/bot/internal/database"
	"github.com/kiasuo/bot/internal/helpers"
)

var adminID int64

func init() {
	var err error

	if adminID, err = strconv.ParseInt(helpers.GetEnv("ADMIN_ID"), 10, 64); err != nil {
		log.Panic(err)
	}
}

const AdminCommandName string = "admin"

var AdminCommand = Command(func(ctx Context, resp Responder, formatter helpers.Formatter) error {
	if ctx.User.TelegramID != adminID {
		return nil
	}

	return resp.Write("OK").Respond()

	//user := ctx.User
	//
	//text := formatter.Title("Панель управления") +
	//	formatter.Item("Telegram: "+strconv.FormatInt(user.TelegramID, 10)) +
	//	formatter.Item("Статус: "+user.State.String()) +
	//	formatter.Item("Версия: "+version.Version)
	//
	//keyboard := Keyboard{
	//	KeyboardRow{
	//		NewKeyboardButton(
	//			helpers.If(user.State == database.UserStateBlacklisted, "Разблокировать", "Заблокировать"),
	//			AdminCommandName+":blacklist:"+strconv.FormatInt(user.TelegramID, 10),
	//		),
	//	},
	//}
	//
	//return resp.Write(text).RespondWithKeyboard(keyboard)
})

var AdminCallback = Callback(func(ctx Context, resp Responder, formatter helpers.Formatter, data []string) error {
	switch data[0] {
	case "blacklist":
		isBlacklisted := ctx.User.State == database.UserStateBlacklisted
		ctx.User.SetState(ctx.Context(), helpers.If(isBlacklisted, database.UserStateReady, database.UserStateBlacklisted))
	default:
		return resp.Write("Неизвестная команда. Меню устарело?").Respond()
	}

	return AdminCommand(ctx, resp, formatter)
})
