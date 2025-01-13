package commands

import (
	"github.com/kiasuo/bot/internal/helpers"
	"slices"
	"strings"
	"time"
)

func scheduleCommand(context Context, responder Responder, formatter helpers.Formatter, t time.Time) error {
	data, err := context.GetClient().GetSchedule(t)

	if err != nil {
		return err
	}

	keyboard := Keyboard{
		KeyboardRow{
			KeyboardButton{
				Text:     "Предыдущая неделя",
				Callback: "schedule:" + t.AddDate(0, 0, -7).Format(time.DateOnly),
			},
			{
				Text:     "Следующая неделя",
				Callback: "schedule:" + t.AddDate(0, 0, 7).Format(time.DateOnly),
			},
		},
	}

	if len(data.Schedule) == 0 {
		return responder.Write("Расписания нет. Отдыхаем?").RespondWithKeyboard(keyboard)
	}

	date := ""
	var checked []string

	for _, event := range data.Schedule {
		if event.LessonDate != date {
			responder.Write(formatter.Title(formatDate(event.Date())))
			date = event.LessonDate
			checked = nil
		}

		responder.Write(formatter.Line(event.String()))

		if len(event.Slots) > 0 {
			marks := ""

			for i, slot := range event.Slots {
				if i > 0 {
					marks += ", "
				}

				marks += slot.Mark.Value
			}

			responder.Write(formatter.Item("Оценки: " + formatter.Code(marks)))
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
							responder.Write(formatter.Item(formatter.Block(text)))
						} else {
							responder.Write(formatter.Item(text))
						}
					}

					checked = append(checked, text)
				}

				for _, file := range homework.Files {
					responder.Write(formatter.Item(file.String(formatter)))
				}

				for _, link := range homework.Links {
					responder.Write(formatter.Item(link.String(formatter)))
				}
			}
		}
	}

	return responder.RespondWithKeyboard(keyboard)
}

var ScheduleCommand = Command(func(context Context, responder Responder, formatter helpers.Formatter) error {
	return scheduleCommand(context, responder, formatter, time.Now())
})

var ScheduleCallback = Callback(func(context Context, responder Responder, formatter helpers.Formatter, data []string) error {
	time, _ := time.Parse(time.DateOnly, data[1])
	return scheduleCommand(context, responder, formatter, time)
})

func formatDate(t time.Time) string {
	return map[time.Weekday]string{
		time.Monday:    "Понедельник",
		time.Tuesday:   "Вторник",
		time.Wednesday: "Среда",
		time.Thursday:  "Четверг",
		time.Friday:    "Пятница",
		time.Saturday:  "Суббота",
		time.Sunday:    "Воскресенье",
	}[t.Weekday()] + ", " + t.Format("02.01")
}
