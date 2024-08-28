package commands

import "sort"

var StudentsCommand = Command(func(context Context, responder Responder, formatter Formatter) error {
	recipients, err := context.GetClient().GetRecipients()

	if err != nil {
		return err
	}

	students := make([]string, 0, len(recipients.Students))

	for student := range recipients.Students {
		students = append(students, student)
	}

	sort.Strings(students)
	result := formatter.Title("Список учеников")

	for _, student := range students {
		result += formatter.Item(student)
	}

	return responder.Respond(result)
})
