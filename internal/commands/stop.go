package commands

import "github.com/kiasuo/bot/internal/helpers"

var StopCommand = Command(func(_ Context, responder Responder, _ helpers.Formatter) error {
	keyboard := Keyboard{
		KeyboardRow{
			KeyboardButton{
				Text:     "Пойти на хуй",
				Callback: "stop:yes",
			},
		},
	}

	return responder.RespondWithKeyboard(keyboard, "Нам очень жаль, что Вы решили покинуть нас. Нажмите ниже, чтобы удалить все данные.")
})

var StopCallback = Callback(func(context Context, responder Responder, formatter helpers.Formatter, data []string) error {
	return responder.Respond("Скоро сделаю!")
})
