package commands

import "log"

func HandleCommand(context Context, responder Responder, formatter Formatter) {
	log.Println("Handling command", context.Command)

	var err error

	switch context.Command {
	case AdminCommandName:
		AdminCommand(context, responder, formatter)
	case "start":
		StartCommand(context, responder)
	case "stop":
		StopCommand(context, responder)
	case "settings":
		SettingsCommand(context, responder, formatter)
	case "students":
		err = StudentsCommand(context, responder, formatter)
	case "staff":
		err = StaffCommand(context, responder, formatter)
	}

	handleError(err)
}

func HandleCallback(context Context, responder Responder, formatter Formatter, data []string) {
	log.Println("Handling callback", context.Command)

	var err error

	switch context.Command {
	case AdminCommandName:
		AdminCallback(context, responder, formatter, data)
	case "stop":
		StopCallback(context, responder, formatter, data)
	case "settings":
		err = SettingsCallback(context, responder, formatter, data)
	}

	handleError(err)
}

func handleError(err error) {
	if err != nil {
		log.Printf("Error: %v", err)
	}
}
