package commands

import (
	"slices"

	"github.com/kiasuo/bot/internal/helpers"
)

var ClassmatesCommand = Command(func(ctx Context, resp *Responder, formatter helpers.Formatter) error {
	recipients, err := ctx.GetClient().GetRecipients()

	if err != nil {
		return err
	}

	students := make([]string, 0, len(recipients.Students))

	for student := range recipients.Students {
		students = append(students, student)
	}

	slices.Sort(students)
	resp.Write(formatter.Title("Список учеников (%d)"), len(students))

	for _, student := range students {
		resp.Write(formatter.Item(student))
	}

	return resp.Respond()
})
