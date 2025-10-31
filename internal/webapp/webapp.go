package webapp

import (
	"github.com/kiasuo/bot/internal/client"
	"github.com/kiasuo/bot/internal/helpers"
)

type StudentPage struct {
	StudentNameAcronym string `json:"studentNameAcronym"`
}

type MarksPage struct {
	StudyPeriod      client.StudyPeriod   `json:"studyPeriod"`
	StudyPeriods     []client.StudyPeriod `json:"studyPeriods"`
	Lessons          []client.Lesson      `json:"lessons"`
	ShowPasses       bool                 `json:"showPasses"`
	ShowEmptyLessons bool                 `json:"showEmptyLessons"`
	LastMarksSeenAt  int64                `json:"lastMarksSeenAt"`
}

var publicUrl string

func init() {
	if helpers.IsTesting() {
		return
	}

	publicUrl = helpers.GetEnv("WEBAPP_URL")
}

func URL() string {
	return publicUrl
}

func MarksURL() string {
	return publicUrl + "/webapp/marks"
}
