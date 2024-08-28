package commands

var StartCommand = Command(func(_ Context, responder Responder, _ Formatter) error {
	return responder.Respond("Привет!")
})
