package commands

import (
	"fmt"
	"strconv"

	"github.com/kiasuo/bot/internal/client"
	"github.com/kiasuo/bot/internal/database"
	"github.com/kiasuo/bot/internal/helpers"
)

const SettingsCommandName = "settings"

var SettingsCommand = Command(func(ctx Context, resp Responder, formatter helpers.Formatter) error {
	user := ctx.User

	keyboard := Keyboard{
		KeyboardRow{
			NewKeyboardButton(
				"Выбрать ученика",
				"settings:userStudents",
			),
			NewKeyboardButton(
				helpers.If(user.HasFlag(database.UserFlagCache), "Отключить", "Включить")+" кэширование",
				"settings:cache",
			),
		},
	}

	return resp.
		Write("Ученик: " + formatter.Bold(user.GetStudentNameAcronym())).
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
	case "showPasses", "showEmptyLessons":
		return updateMarks(ctx, resp, formatter, data)
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

	if err = ctx.User.SetStudent(ctx.Context(), studentID, studentNameAcronym); err != nil {
		return err
	}

	return resp.
		Write("Ученик %s успешно выбран!", formatter.Bold(studentNameAcronym)).
		Respond()
}

func updateCache(ctx Context, response Responder) error {
	val := !ctx.User.HasFlag(database.UserFlagCache)
	err := ctx.User.SetFlag(ctx.Context(), database.UserFlagCache, val)

	if err != nil {
		return err
	}

	if val {
		return response.Write("Кэширование успешно включено!").Respond()
	}

	ok := ctx.GetClient().PurgeCache()
	return response.Write("Кэширование успешно отключено." + helpers.If(ok, " Кэш очищен.", "")).Respond()
}

func getMarks(ctx Context, resp Responder, formatter helpers.Formatter) error {
	user := ctx.User

	keyboard := Keyboard{
		KeyboardRow{
			NewKeyboardButton(
				helpers.If(user.HasFlag(database.UserFlagShowPasses), "Скрывать", "Отображать")+" пропуски",
				"settings:showPasses",
			),
			NewKeyboardButton(
				helpers.If(user.HasFlag(database.UserFlagShowEmptyLessons), "Скрывать", "Отображать")+" пустые предметы",
				"settings:showEmptyLessons",
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

func updateMarks(ctx Context, resp Responder, formatter helpers.Formatter, data []string) error {
	var flag database.UserFlag

	switch data[0] {
	case "showEmptyLessons":
		flag = database.UserFlagShowEmptyLessons
	case "showPasses":
		flag = database.UserFlagShowPasses
	default:
		return fmt.Errorf("unknown flag: %s", data[1])
	}

	if err := ctx.User.SetFlag(ctx.Context(), flag, !ctx.User.HasFlag(flag)); err != nil {
		return err
	}

	return getMarks(ctx, resp, formatter)
}
