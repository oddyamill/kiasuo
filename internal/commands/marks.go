package commands

import (
	"strconv"
	"time"

	"github.com/kiasuo/bot/internal/client"
	"github.com/kiasuo/bot/internal/helpers"
)

func marksCommand(context Context, responder Responder, formatter helpers.Formatter, periodID int) error {
	periods, err := context.GetClient().GetStudyPeriods()

	if err != nil {
		return err
	}

	row := KeyboardRow{}
	now := time.Now()
	var period *client.StudyPeriod

	for _, p := range *periods {
		row = append(row, NewKeyboardButton(p.Text, "marks:"+strconv.Itoa(p.ID)))

		if periodID == p.ID || (periodID == 0 && p.Match(now)) {
			period = &p
		}
	}

	row = append(row, NewKeyboardButton("Настройки", "settings:marks"))

	keyboard := Keyboard{row}

	if period == nil {
		return responder.Write("Каникулы?").RespondWithKeyboard(keyboard)
	}

	marks, err := context.GetClient().GetLessons(period.ID)

	if err != nil {
		return err
	}

	responder.Write(formatter.Title("Оценки за " + period.Text))

	// TODO!
	hidePasses, hideEmptyLessons := true, true

	for _, lesson := range *marks {
		line := ""

		for i, slot := range lesson.Slots {
			mark := slot.Mark

			if hidePasses && mark.IsPass() {
				continue
			}

			if i > 0 && line != "" {
				line += ", "
			}

			line += mark.Value

			if slot.UpdatedAt.After(context.User.LastMarksUpdate) {
				line += "⁺"
			}
		}

		if line == "" {
			if hideEmptyLessons {
				continue
			}

			line = "-"
		}

		responder.Write(formatter.Item(lesson.String() + ": " + formatter.Code(line)))
	}

	context.User.UpdateLastMarksUpdate()

	return responder.RespondWithKeyboard(keyboard)
}

var MarksCommand = Command(func(context Context, responder Responder, formatter helpers.Formatter) error {
	return marksCommand(context, responder, formatter, 0)
})

var MarksCallback = Callback(func(context Context, responder Responder, formatter helpers.Formatter, data []string) error {
	id, _ := strconv.Atoi(data[0])
	return marksCommand(context, responder, formatter, id)
})
