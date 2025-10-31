package commands

import (
	"slices"
	"strings"
	"time"

	"github.com/kiasuo/bot/internal/helpers"
)

func scheduleCommand(ctx Context, resp *Responder, formatter helpers.Formatter, t time.Time) error {
	data, err := ctx.GetClient().GetSchedule(t)

	if err != nil {
		return err
	}

	keyboard := Keyboard{
		KeyboardRow{
			NewCallbackButton(
				"Предыдущая неделя",
				"schedule:"+t.AddDate(0, 0, -7).Format(time.DateOnly),
			),
			NewCallbackButton(
				"Следующая неделя",
				"schedule:"+t.AddDate(0, 0, 7).Format(time.DateOnly),
			),
		},
	}

	if len(data.Schedule) == 0 {
		return resp.Write("Расписания нет. Отдыхаем?").RespondWithKeyboard(keyboard)
	}

	date := ""
	var checked []string

	for _, event := range data.Schedule {
		if event.LessonDate != date {
			resp.Write(formatter.Title(formatDate(event.Date())))
			date = event.LessonDate
			checked = nil
		}

		resp.Write(formatter.Line(event.String()))

		if len(event.Slots) > 0 {
			marks := ""

			for i, slot := range event.Slots {
				if i > 0 {
					marks += ", "
				}

				marks += slot.Mark.Value
			}

			resp.Write(formatter.Item("Оценки: " + formatter.Code(marks)))
		}

		for _, homeworkID := range event.Homeworks {
			for _, homework := range data.Homeworks {
				if homework.ID != homeworkID {
					continue
				}

				text := homework.String()

				if !slices.Contains(checked, text) {
					if text != "" {
						if strings.Contains(text, "\n") {
							resp.Write(formatter.Item(formatter.Block(text)))
						} else {
							resp.Write(formatter.Item(text))
						}
					}

					checked = append(checked, text)
				}

				for _, file := range homework.Files {
					resp.Write(formatter.Item(file.String(formatter)))
				}

				for _, link := range homework.Links {
					resp.Write(formatter.Item(link.String(formatter)))
				}
			}
		}
	}

	return resp.RespondWithKeyboard(keyboard)
}

var ScheduleCommand = Command(func(ctx Context, resp *Responder, formatter helpers.Formatter) error {
	return scheduleCommand(ctx, resp, formatter, time.Now())
})

var ScheduleCallback = Callback(func(ctx Context, resp *Responder, formatter helpers.Formatter, data []string) error {
	time, _ := time.Parse(time.DateOnly, data[0])
	return scheduleCommand(ctx, resp, formatter, time)
})

var weekdays = map[time.Weekday]string{
	time.Monday:    "Понедельник",
	time.Tuesday:   "Вторник",
	time.Wednesday: "Среда",
	time.Thursday:  "Четверг",
	time.Friday:    "Пятница",
	time.Saturday:  "Суббота",
	time.Sunday:    "Воскресенье",
}

func formatDate(t time.Time) string {
	return weekdays[t.Weekday()] + ", " + t.Format("02.01")
}
