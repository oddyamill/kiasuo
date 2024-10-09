package commands

import (
	"github.com/kiasuo/bot/internal/crypto"
	"github.com/kiasuo/bot/internal/helpers"
	"github.com/kiasuo/bot/internal/users"
)

const StartCommandName = "start"

var StartCommand = Command(func(context Context, responder Responder, _ helpers.Formatter) error {
	token := context.Arguments

	if len(token) == 32 && helpers.IsHexUnsafe(token) {
		context.User.RefreshToken = crypto.Encrypt(token)

		if err := context.GetClient().RefreshToken(); err != nil {
			return err
		}

		context.User.UpdateState(users.Ready)
		return responder.Write("Токен успешно обновлен!").Respond()
	}

	return responder.Write("Привет!").Respond()
})
