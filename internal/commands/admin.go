package commands

import (
	"strconv"

	"github.com/kiasuo/bot/internal/helpers"
)

var adminID int64

func init() {
	adminID, _ = strconv.ParseInt(helpers.GetEnv("ADMIN_ID"), 10, 64)
}

const AdminCommandName string = "admin"

var AdminCommand = Command(func(ctx Context, resp Responder, formatter helpers.Formatter) error {
	if ctx.User.TelegramID != adminID {
		return nil
	}

	return resp.Write("OK").Respond()
})
