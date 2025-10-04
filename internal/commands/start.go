package commands

import (
	"github.com/kiasuo/bot/internal/crypto"
	"github.com/kiasuo/bot/internal/helpers"
	"github.com/kiasuo/bot/internal/users"
)

const StartCommandName = "start"

var StartCommand = Command(func(ctx Context, resp Responder, _ helpers.Formatter) error {
	token := ctx.Arguments

	if len(token) == 32 && helpers.IsHexUnsafe(token) {
		ctx.User.RefreshToken = crypto.Encrypt(token)

		if err := ctx.GetClient().RefreshToken(); err != nil {
			return err
		}

		ctx.User.UpdateState(users.Ready)
		return resp.Write("Токен успешно обновлен!").Respond()
	}

	return resp.Write("Привет!").Respond()
})
