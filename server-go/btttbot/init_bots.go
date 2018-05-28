package btttbot

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/strongo/app"
	"github.com/strongo/app/gaestandard"
	"github.com/strongo/bidding-tictactoe-bot/btttbot-secrets"
	"github.com/strongo/bidding-tictactoe-bot/btttbot/platforms/tgbots"
	"github.com/strongo/bidding-tictactoe-bot/btttbot/routing"
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/bots-framework/platforms/telegram"
	"net/http"
)

func InitBot(httpRouter *httprouter.Router, botHost bots.BotHost, appContext bots.BotAppContext) error {
	gaSettings := bots.AnalyticsSettings{
		GaTrackingID: btttbot_secrets.GaTrackingID,
		Enabled: func(r *http.Request) bool {
			return false
		},
	}

	driver := bots.NewBotDriver(gaSettings, appContext, botHost,
		"Please report any issues to @",
	)
	//routing.WebhooksRouter

	if driver == nil {
		panic("driver == nil")
	}

	newTranslator := func(c context.Context) strongo.Translator {
		return strongo.NewMapTranslator(c, nil)
	}

	driver.RegisterWebhookHandlers(httpRouter, "/bot",
		telegram.NewTelegramWebhookHandler(
			func(c context.Context) bots.SettingsBy {
				return tgbots.Bots(c, gaestandard.GetEnvironment(c), routing.WebhooksRouter)
			}, // Maps of bots by code, language, token, etc...
			newTranslator, // Creates newTranslator that gets a context.Context (for logging purpose)
		),
	)

	return nil
}
