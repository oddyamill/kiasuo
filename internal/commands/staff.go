package commands

var StaffCommand = Command(func(context Context, responder Responder, formatter Formatter) error {
	recipients, err := context.GetClient().GetRecipients()

	if err != nil {
		return err
	}

	result := ""

	for role := range recipients.Staff {
		if len(recipients.Staff[role]) == 0 {
			continue
		}

		result += formatter.Title(role)

		for staff := range recipients.Staff[role] {
			result += formatter.Item(staff)
		}
	}

	return responder.Respond(result)
})
