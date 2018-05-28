package commands

import (
	"context"
	"github.com/pkg/errors"
	"github.com/qedus/nds"
	"github.com/strongo-games/bidding-tictactoe/server-go/btttdal"
	"github.com/strongo-games/bidding-tictactoe/server-go/btttmodels"
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/log"
	"net/url"
	"strconv"
)

var NewGameCommand = bots.Command{
	Code: "new_game",
	CallbackAction: func(whc bots.WebhookContext, callbackURL *url.URL) (m bots.MessageFromBot, err error) {
		c := whc.Context()
		log.Debugf(c, "NewGameCommand.CallbackAction()")

		var oldGame, newGame btttmodels.Game
		if oldGame.ID, err = strconv.ParseInt(callbackURL.Query().Get("g"), 10, 64); err != nil {
			return
		}

		if oldGame, err = btttdal.Game.GetGameByID(c, oldGame.ID); err != nil {
			return
		}

		user := btttmodels.NewAppUser(whc.AppUserIntID())

		var botUser bots.BotAppUser
		if botUser, err = whc.GetAppUser(); err != nil {
			return
		}
		user.AppUserEntity = botUser.(*btttmodels.AppUserEntity)

		if err = nds.RunInTransaction(c, func(c context.Context) error {
			if newGame, err = btttdal.NewGameInTelegramChat(c, oldGame.Board.Size(), oldGame.Locale, user.ID, user.FullName(), oldGame.TgChatInstance, "", oldGame); err != nil {
				return errors.Wrap(err, "Failed to create new game in datastore")
			}
			rivalUserID := newGame.RivalUserUD(user.ID)

			if err = updateOtherMessages(c, whc.GetBotCode(), newGame, user.ID, false); err != nil {
				return err
			}
			if err = updateOtherMessages(c, whc.GetBotCode(), newGame, rivalUserID, false); err != nil {
				return err
			}
			if err = updateOtherMessages(c, whc.GetBotCode(), oldGame, user.ID, false); err != nil {
				return err
			}
			if err = updateOtherMessages(c, whc.GetBotCode(), oldGame, rivalUserID, false); err != nil {
				return err
			}
			return err
		}, nil); err != nil {
			return
		}

		return
	},
}
