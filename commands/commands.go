package commands

import "log"

func HandleCommand(context Context, responder Responder, formatter Formatter) {
	log.Println("Handling command", context.Command)

	switch context.Command {
	case AdminCommandName:
		AdminCommand(context, responder, formatter)
	case "start":
		StartCommand(context, responder)
	case "settings":
		SettingsCommand(context, responder, formatter)
	case "students":
		StudentsCommand(context, responder, formatter)
	case "staff":
		StaffCommand(context, responder, formatter)
	}
}

func HandleCallback(context Context, responder Responder, formatter Formatter, data []string) {
	log.Println("Handling callback", context.Command)

	switch context.Command {
	case AdminCommandName:
		AdminCallback(context, responder, formatter, data)
	}
}
