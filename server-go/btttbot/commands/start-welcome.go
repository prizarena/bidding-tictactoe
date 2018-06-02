package commands

import (
	"github.com/prizarena/bidding-tictactoe/server-go/bttt-trans"
	"github.com/strongo/bots-api-telegram"
	"github.com/strongo/bots-framework/core"
)

func welcomeUser(whc bots.WebhookContext) (m bots.MessageFromBot, err error) {
	//c := whc.Context()
	m = whc.NewMessage(whc.Translate(bttt_trans.MT_WELCOME) + "\n\n" + whc.Translate(bttt_trans.MT_HOW_TO_START_NEW_GAME))
	switchInlineQuery := "new_game"
	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text:              whc.Translate(bttt_trans.C_NEW_GAME),
				SwitchInlineQuery: &switchInlineQuery,
			},
		),
	)

	//if _, err = whc.Responder().SendMessage(c, m, bots.BotAPISendMessageOverHTTPS); err != nil {
	//	return
	//}
	//
	//mt = whc.Translate(bttt_trans.MT_TOURNAMENT_201710_SHORT) +
	//	"\n" + strings.Repeat("-", 12) +
	//	"\n" + whc.Translate(bttt_trans.MT_TOURNAMENT_201710_SPONSOR) +
	//	"\n" + whc.Translate(bttt_trans.MT_TOURNAMENT_201710_LEARN_MORE)
	//
	//m = whc.NewMessage(mt)
	//m.DisableWebPagePreview = true
	return
}
