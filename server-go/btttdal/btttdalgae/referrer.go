package btttdalgae

import (
	"context"
	"github.com/strongo-games/bidding-tictactoe/server-go/btttdal"
	"github.com/strongo-games/bidding-tictactoe/server-go/btttmodels"
)

type referrerDalGae struct {
}

func (referrerDalGae) GetReferrerByID(c context.Context, id string) (referrer btttmodels.Referrer, err error) {
	referrer.ID = id
	err = btttdal.DB.Get(c, &referrer)
	return
}
