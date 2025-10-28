package main

import (
	"encoding/json"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kiasuo/bot/internal/database"
	"github.com/kiasuo/bot/internal/helpers"
	"github.com/kiasuo/bot/internal/response"
)

type App struct {
	db                *database.DB
	bot               *tgbotapi.BotAPI
	updates           chan tgbotapi.Update
	httpTelegramToken string
}

func (app *App) telegramHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Telegram-Bot-Api-Secret-Token") != app.httpTelegramToken {
		response.ErrUnauthorized(w)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		response.ErrUnsupportedMediaType(w)
		return
	}

	var update tgbotapi.Update

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		response.ErrBadRequest(w, err.Error())
		return
	}

	app.updates <- update
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
	http.HandleFunc("POST /internal/webhook", app.telegramHandler)
	log.Panic(http.ListenAndServe(":39814", nil))
}
