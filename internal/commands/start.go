package commands

import (
	"github.com/kiasuo/bot/internal/crypto"
	"github.com/kiasuo/bot/internal/database"
	"github.com/kiasuo/bot/internal/helpers"
)

const StartCommandName = "start"

var StartCommand = Command(func(ctx Context, resp Responder, formatter helpers.Formatter) error {
	if ctx.User.State == database.UserStateReady {
		return resp.Write("Привет!").Respond()
	}

	token := ctx.Arguments

	if len(token) == 32 && helpers.IsHexUnsafe(token) {
		ctx.User.RefreshToken = crypto.Encrypt(token).Encrypted

		if err := ctx.GetClient().RefreshToken(); err != nil {
			return err
		}

		if err := ctx.User.SetState(ctx.Context(), database.UserStateReady); err != nil {
			return err
		}

		return resp.Write("Токен успешно обновлен!").Respond()
	}

	return resp.
		Write("Привет! Для работы бота необходимо пройти регистрацию: https://kiasuo.oddya.ru/registration").
		Respond()
})
