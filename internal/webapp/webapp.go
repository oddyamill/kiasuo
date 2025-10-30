package webapp

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"strings"

	"github.com/kiasuo/bot/internal/client"
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

func ValidateTelegramInit(botToken string, data string) (url.Values, bool) {
	appData, err := url.ParseQuery(data)

	if err != nil {
		return nil, false
	}

	hash := appData.Get("hash")

	if hash == "" {
		return nil, false
	}

	appData.Del("hash")
	appDataToCheck, _ := url.QueryUnescape(strings.ReplaceAll(appData.Encode(), "&", "\n"))
	secretKey := hmacHash([]byte(botToken), []byte("WebAppData"))

	if hex.EncodeToString(hmacHash([]byte(appDataToCheck), secretKey)) != hash {
		return nil, false
	}

	return appData, true
}

func hmacHash(data, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	_, _ = h.Write(data)
	return h.Sum(nil)
}
