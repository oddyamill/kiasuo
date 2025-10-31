package commands

import (
	"log"
	"strconv"

	"github.com/kiasuo/bot/internal/helpers"
)

var adminID int64

func init() {
	if helpers.IsTesting() {
		return
	}

	var err error

	if adminID, err = strconv.ParseInt(helpers.GetEnv("ADMIN_ID"), 10, 64); err != nil {
		log.Fatal(err)
	}
}

const AdminCommandName string = "admin"

var AdminCommand = Command(func(ctx Context, resp Responder, formatter helpers.Formatter) error {
	if ctx.User.TelegramID != adminID {
		return nil
	}

	return resp.Write("OK").Respond()
})
