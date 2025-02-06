package client

import (
	"strconv"
	"strings"
)

const (
	PublicUrl = "https://diaryapi.kiasuo.ru"
	ApiUrl    = "https://kiasuo.oddya.ru/diary"
)

var (
	refreshURL      = ApiUrl + "/refresh"
	userURL         = ApiUrl + "/api/user"
	recipientsURL   = ApiUrl + "/api/recipients"
	studyPeriodsURL = ApiUrl + "/api/study_periods"
	revokeURL       = ApiUrl + "/pwa_logout"

	purgeCacheURL = ApiUrl + "/../internal/purge-cache"
)

func lessonMarksURL(id int) string {
	return ApiUrl + "/api/lesson_marks/" + strconv.Itoa(id)
}

func scheduleURL(year, week int) string {
	return ApiUrl + "/api/schedule?year=" + strconv.Itoa(year) + "&week=" + strconv.Itoa(week)
}

func appendID(rawUrl string, id *int) string {
	if id == nil {
		return rawUrl
	}

	if strings.Contains(rawUrl, "?") {
		return rawUrl + "&id=" + strconv.Itoa(*id)
	}

	return rawUrl + "?id=" + strconv.Itoa(*id)
}
