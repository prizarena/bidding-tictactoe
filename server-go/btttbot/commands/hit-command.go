package commands

import (
	"github.com/strongo-games/bidding-tictactoe/server-go/btttmodels"
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/log"
	"net/url"
)

var HitCommand = bots.Command{
	Code: "hit",
	CallbackAction: func(whc bots.WebhookContext, callbackURL *url.URL) (m bots.MessageFromBot, err error) {
		c := whc.Context()
		log.Debugf(c, "HitCommand.CallbackAction()")
		var turnRequest btttmodels.TurnRequest
		if turnRequest, err = btttmodels.ParseQueryToTurnRequest(callbackURL.RawQuery); err != nil {
			return
		}
		return processHitForAlreadyJoinedPlayer(whc, turnRequest)
	},
}
