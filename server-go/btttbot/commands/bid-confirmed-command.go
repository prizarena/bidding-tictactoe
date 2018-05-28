package commands

import (
	"github.com/pkg/errors"
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/log"
)

var errBidTooBig = errors.New("Bid too big")

func bidConfirmed(whc bots.WebhookContext, query string) (m bots.MessageFromBot, err error) {
	c := whc.Context()
	log.Debugf(c, "bidConfirmed()")
	matches := ReInlineBid.FindStringSubmatch(query)
	var (
		gameID int64
		x, y   int8
		bid    int16
	)
	if gameID, x, y, bid, err = inlineQueryMatchesToBidInfo(c, matches); err != nil {
		return
	}

	//inlineMessageID := whc.Input().(telegram.TelegramWebhookChosenInlineResult).GetInlineMessageID()

	return processBid(whc, gameID, bid, x, y)
}
