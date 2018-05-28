package btttdalgae

import (
	"context"
	"github.com/strongo-games/bidding-tictactoe/server-go/btttdal"
	"github.com/strongo-games/bidding-tictactoe/server-go/btttmodels"
	"google.golang.org/appengine/datastore"
)

type gameDal struct {
}

var _ btttdal.GameDal = (*gameDal)(nil)

func (gameDal) InsertGame(c context.Context, gameEntity *btttmodels.GameEntity) (game btttmodels.Game, err error) {
	game = btttmodels.Game{
		GameEntity: gameEntity,
	}
	err = btttdal.DB.InsertWithRandomIntID(c, &game)
	return
}

func (gameDal) GetGameByID(c context.Context, id int64) (game btttmodels.Game, err error) {
	game.ID = id
	err = btttdal.DB.Get(c, &game)
	return
}

func (gameDal) SaveGame(c context.Context, game btttmodels.Game) error {
	return btttdal.DB.Update(c, &game)
}

func (gameDal) GetGameByIDbyInlineMessageID(c context.Context, id string) (gameID int64, err error) {
	q := datastore.NewQuery(btttmodels.GameKind).Limit(1).KeysOnly()
	var keys []*datastore.Key
	if keys, err = q.GetAll(c, nil); err != nil {
		return
	}
	gameID = keys[0].IntID()
	return
}
