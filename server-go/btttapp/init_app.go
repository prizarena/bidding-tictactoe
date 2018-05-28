package btttapp

import (
	"github.com/julienschmidt/httprouter"
	"github.com/strongo/bidding-tictactoe-bot/btttbot"
	"github.com/strongo/bidding-tictactoe-bot/btttdal/btttdalgae"
	"github.com/strongo/bots-framework/core"
	"net/http"
)

func InitApp(httpRouter *httprouter.Router, botHost bots.BotHost) {

	btttdalgae.RegisterGaeDal()

	http.Handle("/", httpRouter)

	btttbot.InitBot(httpRouter, botHost, TheAppContext)
}
