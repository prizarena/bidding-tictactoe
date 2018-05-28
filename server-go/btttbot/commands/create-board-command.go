package commands

import (
	"context"
	"github.com/pkg/errors"
	"github.com/qedus/nds"
	"github.com/strongo-games/bidding-tictactoe/server-go/btttbot/common"
	"github.com/strongo-games/bidding-tictactoe/server-go/btttdal"
	"github.com/strongo-games/bidding-tictactoe/server-go/btttmodels"
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/bots-framework/platforms/telegram"
	"net/url"
	"strconv"
)

var CreateBoardCommand = bots.Command{
	Code: "board",
	CallbackAction: func(whc bots.WebhookContext, callbackURL *url.URL) (m bots.MessageFromBot, err error) {
		c := whc.Context()
		var boardSize int
		if boardSize, err = strconv.Atoi(callbackURL.RawQuery[:1]); err != nil {
			return
		}
		localeCode := callbackURL.RawQuery[1:]
		switch localeCode {
		case "en":
			localeCode = "en-US"
		case "ru":
			localeCode = "ru-RU"
		default:
			panic("Unknown locale: " + localeCode)
		}
		if err = whc.SetLocale(localeCode); err != nil {
			return
		}
		var game btttmodels.Game
		user := btttmodels.NewAppUser(whc.AppUserIntID())

		var botUser bots.BotAppUser
		if botUser, err = whc.GetAppUser(); err != nil {
			return
		}
		user.AppUserEntity = botUser.(*btttmodels.AppUserEntity)

		tgUpdate := whc.Input().(telegram.TgWebhookCallbackQuery).TgUpdate()
		callbackQuery := tgUpdate.CallbackQuery
		inlineMessageID := callbackQuery.InlineMessageID

		if err = nds.RunInTransaction(c, func(c context.Context) error {
			if game, err = btttdal.NewGameInTelegramChat(c, boardSize, localeCode, user.ID, user.FullName(), callbackQuery.ChatInstance, inlineMessageID, game); err != nil {
				return errors.Wrap(err, "Failed to create new game in datastore")
			}
			return err
		}, nil); err != nil {
			return
		}

		mode := common.MODE_INLINE

		if whc.Input().Chat() != nil {
			mode = common.MODE_INBOT_NEW
		}

		m, err = common.GameNewMessageToBot(whc, mode, game, btttmodels.NoWinnerYet, user)

		return
	},
}
