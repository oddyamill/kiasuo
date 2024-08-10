package commands

import "log"

func HandleCommand(context Context, responder Responder, formatter Formatter) {
	log.Println("Handling command", context.Command)

	switch context.Command {
	case "start":
		StartCommand(context, responder)
	case "settings":
		SettingsCommand(context, responder, formatter)
	case "students":
		StudentsCommand(context, responder, formatter)
	}
}
