package btttdalgae

import (
	"context"
	"github.com/prizarena/bidding-tictactoe/server-go/btttdal"
	"github.com/prizarena/bidding-tictactoe/server-go/btttmodels"
)

type userDalGae struct {
}

func (userDalGae) GetUserByID(c context.Context, userID int64) (user btttmodels.AppUser, err error) {
	user.ID = userID
	err = btttdal.DB.Get(c, &user)
	return
}

func (userDalGae) SaveUserByID(c context.Context, user btttmodels.AppUser) (err error) {
	return btttdal.DB.Update(c, &user)
}
