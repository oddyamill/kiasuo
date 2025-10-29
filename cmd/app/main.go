package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kiasuo/bot/internal/client"
	"github.com/kiasuo/bot/internal/database"
	"github.com/kiasuo/bot/internal/helpers"
	"github.com/kiasuo/bot/internal/version"
	"github.com/kiasuo/bot/internal/webapp"
)

type App struct {
	db                *database.DB
	bot               *tgbotapi.BotAPI
	updates           chan tgbotapi.Update
	httpTelegramToken string
}

func (app *App) internalWebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Telegram-Bot-Api-Secret-Token") != app.httpTelegramToken {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	var update tgbotapi.Update

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error(err.Error())
		return
	}

	app.updates <- update
}

func (app *App) webappHandler(w http.ResponseWriter, r *http.Request) {
	data := webapp.IndexPage{
		Version: version.Version,
	}

	if err := webapp.Templates.ExecuteTemplate(w, "IndexPage", data); err != nil {
		slog.Error(err.Error())
	}
}

func (app *App) webappMarksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("App-Version") != version.Version {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	initToken := r.Header.Get("Telegram-Init")

	if initToken == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	appData, ok := webapp.ValidateTelegramInit(app.bot.Token, initToken)

	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var tgUser tgbotapi.User

	if err := json.Unmarshal(helpers.StringToBytes(appData.Get("user")), &tgUser); err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := app.db.GetUser(r.Context(), tgUser.ID)

	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if user.State != database.UserStateReady || user.StudentID == 0 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	c := client.New(user)

	studyPeriods, err := c.GetStudyPeriods()

	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var studyPeriod *client.StudyPeriod

	for _, p := range *studyPeriods {
		if p.Match(time.Now()) {
			studyPeriod = &p
			break
		}
	}

	if studyPeriod == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	lessons, err := c.GetLessons(studyPeriod.ID)

	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := webapp.MarksPage{
		Header: webapp.Header{
			StudentNameAcronym: user.GetStudentNameAcronym(),
		},
		StudyPeriod:      *studyPeriod,
		Lessons:          *lessons,
		HidePasses:       !user.HasFlag(database.UserFlagShowPasses),
		HideEmptyLessons: !user.HasFlag(database.UserFlagShowEmptyLessons),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := webapp.Templates.ExecuteTemplate(w, "MarksPage", data); err != nil {
		slog.Error(err.Error())
		return
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI(helpers.GetEnv("TELEGRAM_TOKEN"))

	if err != nil {
		log.Fatal(err)
	}

	app := &App{
		database.New(helpers.GetEnv("REDIS_URL")),
		bot,
		make(chan tgbotapi.Update, bot.Buffer),
		helpers.GetEnv("HTTP_TELEGRAM_TOKEN"),
	}

	go handleBot(app)
	http.HandleFunc("GET /webapp", app.webappHandler)
	http.HandleFunc("GET /webapp/marks", app.webappMarksHandler)
	http.HandleFunc("POST /internal/webhook", app.internalWebhookHandler)
	log.Panic(http.ListenAndServe(":39814", nil))
}
