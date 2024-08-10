package commands

import "strconv"

func SettingsCommand(context Context, responder Responder, formatter Formatter) {
	responder.Respond("Выбранный ученик: %v", formatter.Bold(strconv.Itoa(context.User.StudentID)))
}
