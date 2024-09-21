package commands

import (
	"github.com/kiasuo/bot/internal/client"
	"github.com/kiasuo/bot/internal/helpers"
)

var StartCommand = Command(func(context Context, responder Responder, _ helpers.Formatter) error {
	token := context.Arguments

	if len(token) == 32 {
		context.User.RefreshToken = token
		err := client.RefreshToken(context.GetClient())

		if err != nil {
			return err
		}

		return responder.Write("Токен успешно обновлен!").Respond()
	}

	return responder.Write("Привет!").Respond()
})
