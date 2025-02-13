package commands

import (
	"github.com/kiasuo/bot/internal/helpers"
	"github.com/kiasuo/bot/internal/users"
)

const defaultStopMessage = "Удалить данные"

func getStopMessage(context Context) string {
	if context.User.State != users.Ready {
		return defaultStopMessage
	}

	user, err := context.GetClient().GetUser()

	if err != nil {
		return defaultStopMessage
	}

	if user.Parent {
		return "Сдать ребенка в детдом"
	}

	for _, child := range user.Children {
		if context.User.StudentID != nil && child.ID != *context.User.StudentID {
			continue
		}

		if child.Age >= 18 {
			return "Пойти на хуй"
		}
	}

	return defaultStopMessage
}

var StopCommand = Command(func(context Context, responder Responder, _ helpers.Formatter) error {
	keyboard := Keyboard{
		KeyboardRow{
			NewKeyboardButton(getStopMessage(context), "stop"),
		},
	}

	return responder.
		Write("Нам очень жаль, что Вы решили покинуть нас. Нажмите ниже, чтобы удалить все данные.").
		RespondWithKeyboard(keyboard)
})

var StopCallback = Callback(func(context Context, responder Responder, formatter helpers.Formatter, data []string) error {
	context.User.Delete()
	responder.Write("Все данные удалены.")

	if err := context.GetClient().RevokeToken(); err == nil {
		responder.Write(" Токен был отозван.")
	}

	return responder.Respond()
})
