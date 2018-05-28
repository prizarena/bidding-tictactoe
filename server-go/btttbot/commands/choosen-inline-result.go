package commands

import (
	"github.com/strongo/bots-framework/core"
	//"net/url"
	//"bitbucket.com/debtstracker/gae_app/debtstracker/bot/cmd/dtb_transfer"
	"github.com/strongo/log"
	"strings"
)

var ChosenInlineResultCommand = bots.Command{
	Code:       "inline-create-invite",
	InputTypes: []bots.WebhookInputType{bots.WebhookInputChosenInlineResult},
	Action: func(whc bots.WebhookContext) (m bots.MessageFromBot, err error) {
		c := whc.Context()
		chosenResult := whc.Input().(bots.WebhookChosenInlineResult)
		query := chosenResult.GetQuery()
		log.Infof(c, "ChosenInlineResultCommand.Action(): query: %v", query)
		if strings.HasPrefix(query, "#") {
			m, err = bidConfirmed(whc, query)
		}
		return
	},
}
