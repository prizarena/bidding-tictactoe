package commands

import (
	"context"
	"github.com/pkg/errors"
	"github.com/prizarena/bidding-tictactoe/server-go/btttdal"
	"github.com/prizarena/bidding-tictactoe/server-go/btttmodels"
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/db"
	"github.com/strongo/log"
	"google.golang.org/appengine/datastore"
)

func checkReferrer(whc bots.WebhookContext, textToMatch string) (err error) {
	referrerID := textToMatch[4:]
	if err = btttdal.DB.RunInTransaction(whc.Context(), func(c context.Context) error {
		if _, err = btttdal.Referrer.GetReferrerByID(c, referrerID); err == datastore.ErrNoSuchEntity {
			log.Infof(c, "Referrer not found by id: %v", referrerID)
		} else if err != nil {
			log.Errorf(c, errors.Wrap(err, "Failed to get referrer by ID").Error())
		} else {
			var botUser bots.BotAppUser
			if botUser, err = whc.GetAppUser(); err != nil {
				return err
			}
			user := btttmodels.NewAppUser(whc.AppUserIntID())
			user.AppUserEntity = botUser.(*btttmodels.AppUserEntity)
			user.ReferrerID = referrerID

			if err = btttdal.User.SaveUserByID(c, user); err != nil {
				return err
			}
		}
		return err
	}, db.CrossGroupTransaction); err != nil {
		return
	}
	return
}
