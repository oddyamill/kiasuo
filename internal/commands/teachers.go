package commands

import (
	"github.com/kiasuo/bot/internal/helpers"
	"sort"
	"strings"
)

var TeachersCommand = Command(func(context Context, responder Responder, formatter helpers.Formatter) error {
	recipients, err := context.GetClient().GetRecipients()

	if err != nil {
		return err
	}

	roles := make([]string, 0, len(recipients.Staff))

	for role := range recipients.Staff {
		roles = append(roles, role)
	}

	sort.Strings(roles)

	for _, role := range roles {
		if len(recipients.Staff[role]) == 0 {
			continue
		}

		responder.Write(formatter.Title(role))

		for staff := range recipients.Staff[role] {
			responder.Write(formatter.Item(formatTeacher(staff)))
		}
	}

	return responder.Respond()
})

func formatTeacher(staff string) string {
	if !strings.Contains(staff, "(") {
		return staff
	}

	// my javascript mind is telling me to use regex
	chunks := strings.Split(staff, " (")
	lessons := strings.Split(strings.TrimSuffix(chunks[1], ")"), ", ")
	result := chunks[0] + " ("

	for i, lesson := range lessons {
		if i > 0 {
			result += ", "
		}

		result += helpers.HumanizeLesson(lesson)
	}

	return result + ")"
}
