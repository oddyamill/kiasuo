package commands

import (
	"github.com/kiasuo/bot/internal/crypto"
	"github.com/kiasuo/bot/internal/helpers"
)

var StartCommand = Command(func(context Context, responder Responder, _ helpers.Formatter) error {
	token := context.Arguments

	if len(token) == 32 {
		// crutch
		client := context.GetClient()
		client.User.RefreshToken = crypto.Encrypt(token)

		if err := client.RefreshToken(); err != nil {
			return err
		}

		return responder.Write("Токен успешно обновлен!").Respond()
	}

	return responder.Write("Привет!").Respond()
})
