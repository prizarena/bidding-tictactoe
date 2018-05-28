package commands

import (
	"github.com/strongo/bidding-tictactoe-bot/bttt-trans"
	"github.com/strongo/bots-api-telegram"
	"github.com/strongo/bots-framework/core"
)

func joinTournamentAction(whc bots.WebhookContext) (m bots.MessageFromBot, err error) {
	m = whc.NewMessageByCode(bttt_trans.MT_ENROLLED_TO_TOURNAMENT)
	switchInlineQiery := "new_game"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text:              whc.Translate(""),
				SwitchInlineQuery: &switchInlineQiery,
			},
		),
	)
	m.Keyboard = keyboard
	return
}
