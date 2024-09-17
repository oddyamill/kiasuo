package commands

import "github.com/kiasuo/bot/internal/helpers"

var StartCommand = Command(func(_ Context, responder Responder, _ helpers.Formatter) error {
	return responder.Respond("Привет!")
})
