package commands

import (
	"github.com/prizarena/bidding-tictactoe/server-go/bttt-trans"
	"github.com/strongo/bots-framework/core"
	"strconv"
	"strings"
)

var BidCommand = bots.Command{
	Code: "bid",
	Action: func(whc bots.WebhookContext) (m bots.MessageFromBot, err error) {
		text := strings.TrimSpace(whc.Input().(bots.WebhookTextMessage).Text())

		var (
			bid int
		)

		if bid, err = strconv.Atoi(text); err != nil {
			m = whc.NewMessageByCode(bttt_trans.MT_NOT_A_NUMBER)
			return
		}

		var gameID int64
		if gameID, err = strconv.ParseInt(whc.ChatEntity().GetWizardParam("g"), 10, 64); err != nil {
			return // TODO: Fail gracefully
		}

		return processBid(whc, gameID, int16(bid), 0, 0)
	},
}
