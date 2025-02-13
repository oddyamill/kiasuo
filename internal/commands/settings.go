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
			NewKeyboardButton(
				"Выбрать ученика",
				"settings:userStudents",
			),
			NewKeyboardButton(
				helpers.If(user.DiscordID.Valid, "Отвязать", "Привязать")+" Discord",
				"settings:discord",
			),
			NewKeyboardButton(
				helpers.If(user.Cache, "Отключить", "Включить")+" кэширование",
				"settings:cache",
			),
		},
	}

	return responder.
		Write("Ученик: " + formatter.Bold(user.StudentNameAcronym.Decrypt())).
		RespondWithKeyboard(keyboard)
})

var SettingsCallback = Callback(func(context Context, responder Responder, formatter helpers.Formatter, data []string) error {
	// shitcode
	switch data[0] {
	case "userStudents":
		return getUserStudents(context, responder)
	case "userStudent":
		return updateUserStudent(context, responder, formatter, data)
	case "discord":
		return getDiscord(context, responder)
	case "cache":
		return updateCache(context, responder)
	case "marks":
		return getMarks(context, responder, formatter)
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
			NewKeyboardButton(
				getName(child), "settings:userStudent:"+strconv.Itoa(child.ID)+":"+getNameAcronym(child),
			),
		})
	}

	return responder.Write("Выберите ученика из списка:").RespondWithKeyboard(keyboard)
}

func updateUserStudent(context Context, responder Responder, formatter helpers.Formatter, data []string) error {
	studentID, err := strconv.Atoi(data[1])
	studentNameAcronym := data[2]

	if err != nil {
		return err
	}

	context.User.UpdateStudent(studentID, studentNameAcronym)

	return responder.
		Write("Ученик %s успешно выбран!", formatter.Bold(studentNameAcronym)).
		Respond()
}

func getDiscord(context Context, responder Responder) error {
	if !context.User.DiscordID.Valid {
		return responder.Write("Привязка аккаунта Discord пока не доступна.").Respond()
	}

	context.User.UpdateDiscord("")
	return responder.Write("Аккаунт Discord успешно отвязан!").Respond()
}

func updateCache(context Context, response Responder) error {
	context.User.UpdateCache(!context.User.Cache)

	if context.User.Cache {
		ok := context.GetClient().PurgeCache()
		return response.Write("Кэширование успешно отключено." + helpers.If(ok, " Кэш очищен.", "")).Respond()
	}

	return response.Write("Кэширование успешно включено!").Respond()
}

func getMarks(_ Context, responder Responder, formatter helpers.Formatter) error {
	keyboard := Keyboard{
		KeyboardRow{
			NewKeyboardButton(
				"Скрывать пропуски", "settings:hide_passes",
			),
			NewKeyboardButton(
				"Скрывать пустые предметы", "settings:hide_empty_lines",
			),
		},
		KeyboardRow{
			NewKeyboardButton("Назад", "marks:0"),
		},
	}

	return responder.
		Write("Настройки команды " + formatter.Code("/marks")).
		RespondWithKeyboard(keyboard)
}
