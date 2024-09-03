package commands

import "sort"

var ClassmatesCommand = Command(func(context Context, responder Responder, formatter Formatter) error {
	recipients, err := context.GetClient().GetRecipients()

	if err != nil {
		return err
	}

	students := make([]string, 0, len(recipients.Students))

	for student := range recipients.Students {
		students = append(students, student)
	}

	sort.Strings(students)
	result := formatter.Title("Список учеников (%d)")

	for _, student := range students {
		result += formatter.Item(student)
	}

	return responder.Respond(result, len(students))
})
