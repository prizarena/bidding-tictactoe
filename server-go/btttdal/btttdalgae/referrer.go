package btttdalgae

import (
	"context"
	"github.com/strongo/bidding-tictactoe-bot/btttdal"
	"github.com/strongo/bidding-tictactoe-bot/btttmodels"
)

type referrerDalGae struct {
}

func (referrerDalGae) GetReferrerByID(c context.Context, id string) (referrer btttmodels.Referrer, err error) {
	referrer.ID = id
	err = btttdal.DB.Get(c, &referrer)
	return
}
