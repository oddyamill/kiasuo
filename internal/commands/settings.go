package commands

import (
	"github.com/kiasuo/bot/internal/client"
	"github.com/kiasuo/bot/internal/helpers"
	"strconv"
)

var SettingsCommand = Command(func(context Context, responder Responder, formatter helpers.Formatter) error {
	user := context.User

	keyboard := Keyboard{
		KeyboardRow{
			KeyboardButton{
				Text:     "Выбрать ученика",
				Callback: "settings:userStudents",
			},
			KeyboardButton{
				Text:     helpers.If(user.DiscordID == "", "Привязать", "Отвязать") + " Discord",
				Callback: "settings:discord",
			},
		},
	}

	return responder.
		Write("Ученик: " + formatter.Bold(user.StudentNameAcronym.Decrypt())).
		RespondWithKeyboard(keyboard)
})

var SettingsCallback = Callback(func(context Context, responder Responder, formatter helpers.Formatter, data []string) error {
	switch data[1] {
	case "userStudents":
		return getUserStudents(context, responder)
	case "userStudent":
		return updateUserStudent(context, responder, formatter, data)
	case "discord":
		return getDiscord(context, responder)
	}

	return nil
})

func getName(child client.Child) string {
	return child.LastName + " " + child.FirstName + " " + child.MiddleName
}

func getNameAcronym(child client.Child) string {
	return child.LastName + " " + string([]rune(child.FirstName)[0]) + ". " + string([]rune(child.MiddleName)[0]) + "."
}

func getUserStudents(context Context, responder Responder) error {
	user, err := context.GetClient().GetUser()

	if err != nil {
		return err
	}

	if len(user.Children) == 0 {
		return responder.Write("У вас нет детей.").Respond()
	}

	keyboard := Keyboard{}

	for _, child := range user.Children {
		keyboard = append(keyboard, KeyboardRow{
			KeyboardButton{
				Text:     getName(child),
				Callback: "settings:userStudent:" + strconv.Itoa(child.ID) + ":" + getNameAcronym(child),
			},
		})
	}

	return responder.Write("Выберите ученика из списка:").RespondWithKeyboard(keyboard)
}

func updateUserStudent(context Context, responder Responder, formatter helpers.Formatter, data []string) error {
	studentID, err := strconv.Atoi(data[2])
	studentNameAcronym := data[3]

	if err != nil {
		return err
	}

	context.User.UpdateStudent(studentID, studentNameAcronym)

	return responder.
		Write("Ученик %s успешно выбран!", formatter.Bold(studentNameAcronym)).
		Respond()
}

func getDiscord(context Context, responder Responder) error {
	if context.User.DiscordID == "" {
		return responder.Write("Привязка аккаунта Discord пока не доступна.").Respond()
	}

	context.User.UpdateDiscord("")
	return responder.Write("Аккаунт Discord успешно отвязан!").Respond()
}
