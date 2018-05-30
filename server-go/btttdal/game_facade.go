package btttdal

import (
	"context"
	"github.com/strongo-games/bidding-tictactoe/server-go/btttmodels"
	"github.com/strongo/log"
	"time"
	)

func NewGameInTelegramChat(c context.Context, boardSize int, localeCode5 string, userID int64, userName, tgChatInstance, tgInlineMessageID string, oldGame btttmodels.Game) (game btttmodels.Game, err error) {
	log.Debugf(c, "boardSize=%d, userID=%d, userName=%v, oldGame.ID=%d", boardSize, userID, userName, oldGame.ID)
	const initialBalance = 100;

	game.GameEntity = &btttmodels.GameEntity{
		Status:            "new",
		Locale:            localeCode5,
		TgChatInstance:    tgChatInstance,
		TgInlineMessageID: tgInlineMessageID,
		DtLastTurn:        time.Now(),
		Board:             btttmodels.EmptyBoard(boardSize),
	}
	if oldGame.ID == 0 {

		game.XUserID = userID
		game.SetPlayerJson(userID, btttmodels.GamePlayerJson{
			Name: userName,
			Balance: initialBalance,
		})
	} else {
		game.PrevGameID = oldGame.ID
		game.XUserID = oldGame.OUserID
		game.OUserID = oldGame.XUserID

		oldUserX, oldUserO := oldGame.GetUsersXO()

		game.SetPlayers(
			btttmodels.GamePlayerJson{Balance: initialBalance, Name: oldUserO.Name, Tg: oldUserO.Tg},
			btttmodels.GamePlayerJson{Balance: initialBalance, Name: oldUserX.Name, Tg: oldUserX.Tg},
		)
	}
	if game, err = Game.InsertGame(c, game.GameEntity); err != nil {
		return
	}
	return
}
