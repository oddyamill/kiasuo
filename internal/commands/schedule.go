package commands

import (
	"github.com/kiasuo/bot/internal/helpers"
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
		return responder.RespondWithKeyboard(keyboard, "Расписания нет. Отдыхаем?")
	}

	var result strings.Builder
	result.WriteString(formatter.Title("Расписание"))
	date := ""

	for _, event := range data.Schedule {
		if event.LessonDate != date {
			result.WriteString(formatter.Title(formatDate(event.Date())))
			date = event.LessonDate
		}

		result.WriteString(formatter.Line(event.String()))

		if len(event.Marks) > 0 {
			marks := "Оценки: "

			for i, mark := range event.Marks {
				if i > 0 {
					marks += ", "
				}

				marks += mark.Mark
			}

			result.WriteString(formatter.Item(marks))
		}

		for _, homeworkId := range event.Homeworks {
			for _, homework := range data.Homeworks {
				if homework.ID != homeworkId {
					continue
				}

				if homework.Text != "" {
					result.WriteString(formatter.Item(homework.Text))
				}

				for _, file := range homework.Files {
					result.WriteString(formatter.Item(file.String(formatter)))
				}

				for _, link := range homework.Links {
					result.WriteString(formatter.Item(link.String(formatter)))
				}
			}
		}
	}

	return responder.RespondWithKeyboard(keyboard, result.String())
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
