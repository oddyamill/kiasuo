package commands

import (
	"sort"
	"strings"

	"github.com/kiasuo/bot/internal/helpers"
)

var TeachersCommand = Command(func(ctx Context, resp Responder, formatter helpers.Formatter) error {
	recipients, err := ctx.GetClient().GetRecipients()

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

		resp.Write(formatter.Title(role))

		for staff := range recipients.Staff[role] {
			resp.Write(formatter.Item(formatTeacher(staff)))
		}
	}

	return resp.Respond()
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
