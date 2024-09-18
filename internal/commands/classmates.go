package commands

import (
	"github.com/kiasuo/bot/internal/helpers"
	"slices"
)

var ClassmatesCommand = Command(func(context Context, responder Responder, formatter helpers.Formatter) error {
	recipients, err := context.GetClient().GetRecipients()

	if err != nil {
		return err
	}

	students := make([]string, 0, len(recipients.Students))

	for student := range recipients.Students {
		students = append(students, student)
	}

	slices.Sort(students)
	responder.Write(formatter.Title("Список учеников (%d)"), len(students))

	for _, student := range students {
		responder.Write(formatter.Item(student))
	}

	return responder.Respond()
})
