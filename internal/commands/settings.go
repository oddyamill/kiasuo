package commands

import (
	"strconv"

	"github.com/kiasuo/bot/internal/client"
	"github.com/kiasuo/bot/internal/helpers"
)

var SettingsCommand = Command(func(ctx Context, resp Responder, formatter helpers.Formatter) error {
	user := ctx.User

	keyboard := Keyboard{
		KeyboardRow{
			NewKeyboardButton(
				"Выбрать ученика",
				"settings:userStudents",
			),
			NewKeyboardButton(
				helpers.If(user.Cache, "Отключить", "Включить")+" кэширование",
				"settings:cache",
			),
		},
	}

	return resp.
		Write("Ученик: " + formatter.Bold(user.StudentNameAcronym.Decrypt())).
		RespondWithKeyboard(keyboard)
})

var SettingsCallback = Callback(func(ctx Context, resp Responder, formatter helpers.Formatter, data []string) error {
	// shitcode
	switch data[0] {
	case "userStudents":
		return getUserStudents(ctx, resp)
	case "userStudent":
		return updateUserStudent(ctx, resp, formatter, data)
	case "cache":
		return updateCache(ctx, resp)
	case "marks":
		return getMarks(ctx, resp, formatter)
	}

	return nil
})

func getName(child client.Child) string {
	return child.LastName + " " + child.FirstName + " " + child.MiddleName
}

func getNameAcronym(child client.Child) string {
	return child.LastName + " " + string([]rune(child.FirstName)[0]) + ". " + string([]rune(child.MiddleName)[0]) + "."
}

func getUserStudents(ctx Context, resp Responder) error {
	user, err := ctx.GetClient().GetUser()

	if err != nil {
		return err
	}

	if len(user.Children) == 0 {
		return resp.Write("У вас нет детей.").Respond()
	}

	keyboard := Keyboard{}

	for _, child := range user.Children {
		keyboard = append(keyboard, KeyboardRow{
			NewKeyboardButton(
				getName(child), "settings:userStudent:"+strconv.Itoa(child.ID)+":"+getNameAcronym(child),
			),
		})
	}

	return resp.Write("Выберите ученика из списка:").RespondWithKeyboard(keyboard)
}

func updateUserStudent(ctx Context, resp Responder, formatter helpers.Formatter, data []string) error {
	studentID, err := strconv.Atoi(data[1])
	studentNameAcronym := data[2]

	if err != nil {
		return err
	}

	ctx.User.UpdateStudent(studentID, studentNameAcronym)

	return resp.
		Write("Ученик %s успешно выбран!", formatter.Bold(studentNameAcronym)).
		Respond()
}

func updateCache(ctx Context, response Responder) error {
	ctx.User.UpdateCache(!ctx.User.Cache)

	if ctx.User.Cache {
		ok := ctx.GetClient().PurgeCache()
		return response.Write("Кэширование успешно отключено." + helpers.If(ok, " Кэш очищен.", "")).Respond()
	}

	return response.Write("Кэширование успешно включено!").Respond()
}

func getMarks(_ Context, resp Responder, formatter helpers.Formatter) error {
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

	return resp.
		Write("Настройки команды " + formatter.Code("/marks")).
		RespondWithKeyboard(keyboard)
}
