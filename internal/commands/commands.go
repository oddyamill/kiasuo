package commands

import "log"

type Command func(Context, Responder, Formatter) error

var CommandMap = map[string]Command{
	AdminCommandName: AdminCommand,
	"start":          StartCommand,
	"stop":           StopCommand,
	"settings":       SettingsCommand,
	"students":       StudentsCommand,
	"staff":          StaffCommand,
}

type Callback func(Context, Responder, Formatter, []string) error

var CallbackMap = map[string]Callback{
	AdminCommandName: AdminCallback,
	"stop":           StopCallback,
	"settings":       SettingsCallback,
}

func HandleCommand(context Context, responder Responder, formatter Formatter) {
	command := CommandMap[context.Command]

	if command == nil {
		log.Println("Someone tried to execute unknown command", context.Command)
		return
	}

	log.Println("Handling command", context.Command)
	handleError(command(context, responder, formatter))
}

func HandleCallback(context Context, responder Responder, formatter Formatter, data []string) {
	callback := CallbackMap[context.Command]

	if callback == nil {
		log.Println("Someone tried to execute unknown callback", context.Command)
		return
	}

	log.Println("Handling callback", context.Command)
	handleError(callback(context, responder, formatter, data))
}

func handleError(err error) {
	if err != nil {
		log.Printf("Error: %v", err)
	}
}
