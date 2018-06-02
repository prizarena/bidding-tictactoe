package commands

import (
	"github.com/strongo/app"
	"github.com/prizarena/bidding-tictactoe/server-go/bttt-trans"
	"github.com/strongo/bots-api-telegram"
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/log"
	"strings"
)

var StartCommand = bots.Command{
	Code:     "start",
	Commands: []string{"/start"},
	Title:    "/start",
	Action: func(whc bots.WebhookContext) (m bots.MessageFromBot, err error) {
		c := whc.Context()
		log.Debugf(c, "StartCommand.Action()")
		textToMatch := whc.Input().(bots.WebhookTextMessage).Text()
		log.Debugf(c, "text: %v", textToMatch)
		if textToMatch == "/start tournament-2017-10-01" {
			return joinTournamentAction(whc)
		} else if strings.HasPrefix(textToMatch, "/start hit-") {
			return startHit(whc, textToMatch)
		} else if strings.HasPrefix(textToMatch, "/start ref-") {
			if err = checkReferrer(whc, textToMatch); err != nil {
				return
			}
		}

		m = whc.NewMessageByCode(bttt_trans.MT_START_SELECT_LANG)
		m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🇬🇧 English", "start-locale?code5="+strongo.LocaleCodeEnUS),
				tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", "start-locale?code5="+strongo.LocalCodeRuRu),
			),
		)

		return
	},
}
