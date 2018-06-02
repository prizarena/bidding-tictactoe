package btttfacade

import (
	"context"
	"github.com/prizarena/bidding-tictactoe/server-go/btttmodels"
	"github.com/prizarena/arena/arena-go"
	"time"
	"github.com/strongo/log"
	"github.com/strongo/db"
	"github.com/prizarena/bidding-tictactoe/server-go/btttdal"
	"strconv"
)

type BidOutput struct {
	RivalKey            arena.BattleID
	Game                btttmodels.Game
	User                btttmodels.AppUser
	RivalUser           btttmodels.AppUser
	UserContestant      arena.Contestant
	RivalUserContestant arena.Contestant
}


type btttStrangerFacade struct {
}

var BtttStrangerFacade = btttStrangerFacade{}

func (btttStrangerFacade) PlaceBidAgainstStranger(c context.Context, userID int64, tournamentID string,  cell btttmodels.Cell, bid int) (err error) {
	now := time.Now()

	onRivalFound := func(rivalUserID string) (err error) {
		log.Debugf(c, "strangerFacade.PlaceBidAgainstStranger() => will link 2 strangers")
		// bidOutput, err = GreedGameFacade.PlaceBidAgainstRival(c, now, userID, tournamentID, rivalUserID, true, bid)
		return
	}

	onStranger := func(contestant *arena.Contestant) error {
		// err = sf.registerNewStranger(c, now, bid, &bidOutput, tournamentID, userID, contestant)
		return err
	}

	user := btttmodels.AppUser{IntegerID: db.NewIntID(userID)}

	if err = arena.MakeMoveAgainstStranger(c, now, tournamentID, &user, onRivalFound, onStranger); err != nil {
		return
	}

	return
}

func (btttStrangerFacade) registerNewStranger(c context.Context, now time.Time, bid int, bidOutput *BidOutput, userID int64, tournamentID string, contestant *arena.Contestant) (err error) {
	var (
		user        btttmodels.AppUser
		// userBattles []models.Battle
	)

	updateUser := func(tc context.Context, strangerRivalKey arena.BattleID) (userEntityHolder db.EntityHolder, err error) {
		if user, err = btttdal.User.GetUserByID(c, userID); err != nil {
			return
		}
		userEntityHolder = &user
		bidOutput.User = user
		// if _, userBattles, err = user.RecordBid(strangerRivalKey, bid, now); err != nil {
		// 	return
		// }
		// user.SetBattles(userBattles)
		return
	}

	userStrID := strconv.FormatInt(userID, 10)
	return arena.RegisterStranger(c, now, tournamentID, userStrID, contestant, updateUser)
}
