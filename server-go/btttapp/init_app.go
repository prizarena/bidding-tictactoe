package btttapp

import (
	"github.com/julienschmidt/httprouter"
	"github.com/strongo-games/bidding-tictactoe/server-go/btttbot"
	"github.com/strongo-games/bidding-tictactoe/server-go/btttdal/btttdalgae"
	"github.com/strongo/bots-framework/core"
	"net/http"
)

func InitApp(httpRouter *httprouter.Router, botHost bots.BotHost) {

	btttdalgae.RegisterGaeDal()

	http.Handle("/", httpRouter)

	btttbot.InitBot(httpRouter, botHost, TheAppContext)
}
