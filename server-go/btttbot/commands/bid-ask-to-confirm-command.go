package commands

import (
	"context"
	"github.com/pkg/errors"
	"github.com/strongo/bidding-tictactoe-bot/bttt-trans"
	"github.com/strongo/bidding-tictactoe-bot/btttcommon"
	"github.com/strongo/bidding-tictactoe-bot/btttdal"
	"github.com/strongo/bidding-tictactoe-bot/btttmodels"
	"github.com/strongo/bots-api-telegram"
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/bots-framework/platforms/telegram"
	"github.com/strongo/log"
	"html"
	"regexp"
	"strconv"
)

var ReInlineBid = regexp.MustCompile(`^#(\w+)-([1-5])([1-5]) => (Your bid|Ваша ставка):\s*(\d*)$`)

func AskToConfirmBid(whc bots.WebhookContext, matches []string) (m bots.MessageFromBot, err error) {
	c := whc.Context()
	log.Debugf(c, "InlineEmptyQuery()")
	var (
		bid  int16
		game btttmodels.Game
	)
	if game.ID, _, _, bid, err = inlineQueryMatchesToBidInfo(c, matches); err != nil || bid == 0 {
		return
	}

	var botUser bots.BotAppUser
	if botUser, err = whc.GetAppUser(); err != nil {
		return
	}
	user := botUser.(*btttmodels.AppUserEntity)

	if game, err = btttdal.Game.GetGameByID(c, game.ID); err != nil {
		log.Errorf(c, errors.Wrap(err, "Failed to get game by ID").Error())
		err = nil
		return
	}

	//whc.ChatEntity().SetPreferredLanguage(game.Locale)
	whc.SetLocale(game.Locale)

	inlineQueryID := whc.Input().(telegram.TgWebhookInlineQuery).GetInlineQueryID()
	m.BotMessage = telegram.InlineBotMessage(tgbotapi.InlineConfig{
		InlineQueryID: inlineQueryID,
		Results: []interface{}{
			tgbotapi.InlineQueryResultArticle{
				ID:          "confirm-bid",
				Type:        "article",
				Title:       whc.Translate(bttt_trans.INLINE_BID_TITLE, bid),
				Description: whc.Translate(bttt_trans.INLINE_BID_DESC),
				InputMessageContent: tgbotapi.InputTextMessageContent{
					Text:                  whc.Translate(bttt_trans.MT_BID_BY, html.EscapeString(user.FullName())),
					ParseMode:             "HTML",
					DisableWebPagePreview: true,
				},
			},
		},
	})

	return m, err
}

func inlineQueryMatchesToBidInfo(c context.Context, matches []string) (gameID int64, x, y int8, bid int16, err error) {
	if gameID, err = btttcommon.DecodeID(matches[1]); err != nil {
		return
	}
	if gameID == 0 {
		log.Debugf(c, "gameID == 0, gameCode: %v", matches[1])
		return
	}
	var (
		xV, yV, bidV int
	)
	if xV, err = strconv.Atoi(matches[2]); err != nil {
		err = nil
		return
	} else {
		x = int8(xV)
	}
	if yV, err = strconv.Atoi(matches[3]); err != nil {
		err = nil
		return
	} else {
		y = int8(yV)
	}
	if len(matches) < 6 || matches[5] == "" {
		bid = 0
	} else if bidV, err = strconv.Atoi(matches[5]); err != nil {
		err = nil
		return
	} else {
		bid = int16(bidV)
	}
	return
}
