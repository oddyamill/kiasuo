package commands

import (
	"github.com/kiasuo/bot/internal/client"
	"github.com/kiasuo/bot/internal/helpers"
	"strconv"
	"time"
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
		row = append(row, KeyboardButton{
			Text:     p.Text,
			Callback: "marks:" + strconv.Itoa(p.ID),
		})

		if periodID == p.ID || (periodID == 0 && p.Match(now)) {
			period = &p
		}
	}

	keyboard := Keyboard{row}

	if period == nil {
		return responder.Write("Каникулы?").RespondWithKeyboard(keyboard)
	}

	marks, err := context.GetClient().GetLessons(period.ID)

	if err != nil {
		return err
	}

	responder.Write(formatter.Title("Оценки за " + period.Text))

	for _, lesson := range *marks {
		line := lesson.String()

		if len(lesson.Marks) > 0 {
			marksLine := ""

			for i, mark := range lesson.Marks {
				if i > 0 {
					marksLine += ", "
				}

				marksLine += mark.Mark

				if mark.UpdatedAt.After(context.User.LastMarksUpdate) {
					marksLine += "⁺"
				}
			}

			line += ": " + formatter.Code(marksLine)
		} else {
			line += ": " + formatter.Code("-")
		}

		responder.Write(formatter.Item(line))
	}

	context.User.UpdateLastMarksUpdate()

	return responder.RespondWithKeyboard(keyboard)
}

var MarksCommand = Command(func(context Context, responder Responder, formatter helpers.Formatter) error {
	return marksCommand(context, responder, formatter, 0)
})

var MarksCallback = Callback(func(context Context, responder Responder, formatter helpers.Formatter, data []string) error {
	id, _ := strconv.Atoi(data[1])
	return marksCommand(context, responder, formatter, id)
})
