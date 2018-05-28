package tgbots

import (
	"github.com/strongo/app"
	"github.com/strongo-games/bidding-tictactoe/server-go/btttbot-secrets"
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/bots-framework/platforms/telegram"
	"context"
	"github.com/strongo/log"
)

var botsBy bots.SettingsBy

const BotProfile = "BiddingTicTacToe"

func Bots(c context.Context, env strongo.Environment, router bots.WebhooksRouter) bots.SettingsBy {
	if len(botsBy.ByCode) == 0 {
		routerByProfile := func(profile string) bots.WebhooksRouter {
			return router // We have single profile for now
		}

		switch env {
		case strongo.EnvProduction:
			botsBy = bots.NewBotSettingsBy(routerByProfile,
				telegram.NewTelegramBot(strongo.EnvProduction, BotProfile,
					btttbot_secrets.TelegramProdBot, btttbot_secrets.TelegramProdToken,
					"", "", btttbot_secrets.GaTrackingID, strongo.LocaleEnUS),
			)
		case strongo.EnvLocal:
			botsBy = bots.NewBotSettingsBy(routerByProfile,
				telegram.NewTelegramBot(strongo.EnvLocal, BotProfile,
					btttbot_secrets.TelegramLocalBot, btttbot_secrets.TelegramLocalToken,
					"", "", "", strongo.LocaleEnUS),
			)
		default:
			log.Errorf(c, "Unknown environment: %v=%v", env, strongo.EnvironmentNames[env])
		}
	}
	return botsBy
}

func BotsBy(_ context.Context) bots.SettingsBy {
	return botsBy
}