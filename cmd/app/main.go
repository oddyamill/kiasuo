package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kiasuo/bot/internal/client"
	"github.com/kiasuo/bot/internal/database"
	"github.com/kiasuo/bot/internal/helpers"
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

func (app *App) internalWebappCors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "https://brand-new-kiasuo-webapp-indev.oddya.ru")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Telegram-Init")
}

func (app *App) authorizeWebappUserAndPutCors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.internalWebappCors(w, r)

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

		next(w, r.WithContext(context.WithValue(r.Context(), "user", user)))
	}
}

func (app *App) internalWebappStudent(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*database.User)

	result := webapp.StudentPage{
		StudentNameAcronym: user.GetStudentNameAcronym(),
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		slog.Error("JSON encoding error", "error", err)
	}
}

func (app *App) internalWebppMarks(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*database.User)
	c := client.New(user)

	studyPeriods, err := c.GetStudyPeriods()
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//studyPeriodID := r.PathValue("studyPeriodID")

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

	lastMarksSeenAt, err := user.GetLastMarksCommand(r.Context(), studyPeriod.ID)

	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result := webapp.MarksPage{
		StudyPeriod:      *studyPeriod,
		StudyPeriods:     *studyPeriods,
		Lessons:          *lessons,
		ShowPasses:       user.HasFlag(database.UserFlagShowPasses),
		ShowEmptyLessons: user.HasFlag(database.UserFlagShowEmptyLessons),
		LastMarksSeenAt:  lastMarksSeenAt.UnixMilli(),
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		slog.Error("JSON encoding error", "error", err)
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
	http.HandleFunc("POST /internal/webhook", app.internalWebhookHandler)
	http.HandleFunc("OPTIONS /internal/webapp/student", app.internalWebappCors)
	http.HandleFunc("GET /internal/webapp/student", app.authorizeWebappUserAndPutCors(app.internalWebappStudent))
	http.HandleFunc("OPTIONS /internal/webapp/marks", app.internalWebappCors)
	http.HandleFunc("GET /internal/webapp/marks", app.authorizeWebappUserAndPutCors(app.internalWebppMarks))
	log.Panic(http.ListenAndServe(":39814", nil))
}
