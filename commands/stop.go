package commands

func StopCommand(_ Context, responder Responder) {
	keyboard := Keyboard{
		KeyboardRow{
			KeyboardButton{
				Text:     "Пойти на хуй",
				Callback: "stop:yes",
			},
		},
	}

	responder.RespondWithKeyboard(keyboard, "Нам очень жаль, что Вы решили покинуть нас. Нажмите ниже, чтобы удалить все данные.")
}

func StopCallback(_ Context, responder Responder, formatter Formatter, data []string) {
	responder.Respond("Скоро сделаю!")
}
