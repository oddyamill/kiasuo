package commands

import (
	"github.com/kiasuo/bot/internal/client"
	"github.com/kiasuo/bot/internal/helpers"
	"github.com/kiasuo/bot/internal/users"
	"strconv"
)

func SettingsCommand(context Context, responder Responder, formatter Formatter) {
	user := context.User

	keyboard := Keyboard{
		KeyboardRow{
			KeyboardButton{
				Text:     "Выбрать ученика",
				Callback: "settings:userStudents",
			},
		},
		KeyboardRow{
			KeyboardButton{
				Text:     helpers.If(user.DiscordID == "", "Привязать Discord", "Отвязать Discord"),
				Callback: "settings:discord",
			},
		},
	}

	responder.RespondWithKeyboard(keyboard, "Ученик: %s", formatter.Bold(user.StudentNameAcronym))
}

func SettingsCallback(context Context, responder Responder, formatter Formatter, data []string) error {
	switch data[1] {
	case "userStudents":
		return getUserStudents(context, responder)
	case "userStudent":
		return updateUserStudent(context, responder, formatter, data)
	case "discord":
		getDiscord(context, responder)
	}

	return nil
}

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
		responder.Respond("У вас нет детей.")
		return nil
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

	responder.RespondWithKeyboard(keyboard, "Выберите ребенка из списка:")
	return nil
}

func updateUserStudent(context Context, responder Responder, formatter Formatter, data []string) error {
	studentID, err := strconv.Atoi(data[2])
	studentNameAcronym := data[3]

	if err != nil {
		return err
	}

	users.UpdateStudent(context.User, studentID, studentNameAcronym)
	responder.Respond("Ученик %s успешно выбран!", formatter.Bold(studentNameAcronym))
	return nil
}

func getDiscord(context Context, responder Responder) {
	if context.User.DiscordID == "" {
		responder.Respond("Привязка аккаунта Discord пока не доступна.")
		return
	}

	users.UpdateDiscord(context.User, "")
	responder.Respond("Аккаунт Discord успешно отвязан!")
}
