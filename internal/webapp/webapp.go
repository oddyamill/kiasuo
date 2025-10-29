package webapp

import (
	"crypto/hmac"
	"crypto/sha256"
	"embed"
	_ "embed"
	"encoding/hex"
	"html/template"
	"net/url"
	"strings"

	"github.com/kiasuo/bot/internal/client"
)

type Header struct {
	StudentNameAcronym string
}

type IndexPage struct {
	Version string
}

type MarksPage struct {
	Header
	StudyPeriod      client.StudyPeriod
	Lessons          []client.Lesson
	HidePasses       bool
	HideEmptyLessons bool
}

//go:embed *.gohtml
var templateFS embed.FS

var Templates = template.Must(template.ParseFS(templateFS, "*.gohtml"))

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
