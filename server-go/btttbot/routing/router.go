package routing

import (
	"github.com/strongo-games/bidding-tictactoe/server-go/btttbot/commands"
	"github.com/strongo/bots-framework/core"
)

var WebhooksRouter = bots.NewWebhookRouter(
	map[bots.WebhookInputType][]bots.Command{
		bots.WebhookInputText: {
			commands.StartCommand,
			commands.BidCommand,
		},
		bots.WebhookInputInlineQuery: {
			InlineQueryCommand,
		},
		bots.WebhookInputChosenInlineResult: {
			commands.ChosenInlineResultCommand,
		},
		bots.WebhookInputCallbackQuery: {
			commands.NewGameCommand,
			commands.HitCommand,
			commands.CreateBoardCommand,
			commands.StartLanguageSelectedCommand,
		},
	},
	func() string { return "Oops..." },
)

var InlineQueryCommand = bots.Command{
	Code:       "inline-query",
	InputTypes: []bots.WebhookInputType{bots.WebhookInputInlineQuery},
	Action: func(whc bots.WebhookContext) (m bots.MessageFromBot, err error) {
		//inlineQuery := whc.Input().(bots.WebhookInlineQuery)
		queryText := whc.Input().(bots.WebhookInlineQuery).GetQuery()
		matches := commands.ReInlineBid.FindStringSubmatch(queryText)
		if matches != nil && len(matches) > 0 {
			m, err = commands.AskToConfirmBid(whc, matches)
		} else {
			m, err = commands.StartNewGame(whc)
		}
		return
	},
}
