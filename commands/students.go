package commands

func StudentsCommand(context Context, responder Responder, formatter Formatter) error {
	recipients, err := context.GetClient().GetRecipients()

	if err != nil {
		return err
	}

	result := formatter.Title("Список учеников")

	for student := range recipients.Students {
		result += formatter.Item(student)
	}

	responder.Respond(result)
	return nil
}
