package commands

import (
	"strconv"
	"time"

	"github.com/kiasuo/bot/internal/client"
	"github.com/kiasuo/bot/internal/helpers"
)

func marksCommand(ctx Context, resp Responder, formatter helpers.Formatter, periodID int) error {
	periods, err := ctx.GetClient().GetStudyPeriods()

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
		return resp.Write("Каникулы?").RespondWithKeyboard(keyboard)
	}

	marks, err := ctx.GetClient().GetLessons(period.ID)

	if err != nil {
		return err
	}

	resp.Write(formatter.Title("Оценки за " + period.Text))

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

			if slot.UpdatedAt.After(ctx.User.LastMarksUpdate) {
				line += "⁺"
			}
		}

		if line == "" {
			if hideEmptyLessons {
				continue
			}

			line = "-"
		}

		resp.Write(formatter.Item(lesson.String() + ": " + formatter.Code(line)))
	}

	ctx.User.UpdateLastMarksUpdate()

	return resp.RespondWithKeyboard(keyboard)
}

var MarksCommand = Command(func(ctx Context, resp Responder, formatter helpers.Formatter) error {
	return marksCommand(ctx, resp, formatter, 0)
})

var MarksCallback = Callback(func(ctx Context, resp Responder, formatter helpers.Formatter, data []string) error {
	id, _ := strconv.Atoi(data[0])
	return marksCommand(ctx, resp, formatter, id)
})
