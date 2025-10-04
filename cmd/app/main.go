package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kiasuo/bot/internal/response"
)

func telegramHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Telegram-Bot-Api-Secret-Token") != os.Getenv("HTTP_TELEGRAM_TOKEN") {
		response.ErrUnauthorized(w)
		return
	}

	response.Ok(w, "Hi!")
}

func main() {
	http.HandleFunc("GET /internal/webhook", telegramHandler)
	log.Panic(http.ListenAndServe(":8080", nil))
}
