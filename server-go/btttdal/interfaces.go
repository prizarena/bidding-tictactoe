package btttdal

import (
	"context"
	"github.com/strongo/bidding-tictactoe-bot/btttmodels"
)

type UserDal interface {
	GetUserByID(c context.Context, userID int64) (user btttmodels.AppUser, err error)
	SaveUserByID(c context.Context, user btttmodels.AppUser) (err error)
}

type GameDal interface {
	InsertGame(c context.Context, gameEntity *btttmodels.GameEntity) (game btttmodels.Game, err error)
	GetGameByID(c context.Context, id int64) (game btttmodels.Game, err error)
	GetGameByIDbyInlineMessageID(c context.Context, id string) (gameID int64, err error)
	SaveGame(c context.Context, game btttmodels.Game) error
}

type ReferrerDal interface {
	GetReferrerByID(c context.Context, id string) (referrer btttmodels.Referrer, err error)
}
