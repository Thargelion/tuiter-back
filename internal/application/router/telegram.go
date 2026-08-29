package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type telegramHandler interface {
	Relay(w http.ResponseWriter, r *http.Request)
}

func NewTelegramRouter(telegramHandler telegramHandler) *TelegramRouter {
	return &TelegramRouter{telegramHandler: telegramHandler}
}

func (tr *TelegramRouter) Route(r chi.Router) {
	r.Post("/", tr.telegramHandler.Relay)
}

type TelegramRouter struct {
	telegramHandler telegramHandler
}
