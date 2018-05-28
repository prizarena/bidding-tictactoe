package btttdal

import (
	"context"
	"github.com/strongo-games/bidding-tictactoe/server-go/btttmodels"
	"github.com/strongo/log"
	"time"
)

func NewGameInTelegramChat(c context.Context, boardSize int, localeCode5 string, userID int64, userName, tgChatInstance, tgInlineMessageID string, oldGame btttmodels.Game) (game btttmodels.Game, err error) {
	log.Debugf(c, "boardSize=%d, userID=%d, userName=%v, oldGame.ID=%d", boardSize, userID, userName, oldGame.ID)
	game.GameEntity = &btttmodels.GameEntity{
		Status:            "new",
		Locale:            localeCode5,
		TgChatInstance:    tgChatInstance,
		TgInlineMessageID: tgInlineMessageID,
		DtLastTurn:        time.Now(),
		Board:             btttmodels.EmptyBoard(boardSize),
		XBalance:          100,
		OBalance:          100,
	}
	if oldGame.ID == 0 {
		game.XUserID = userID
		game.XUserName = userName
	} else {
		game.PrevGameID = oldGame.ID
		game.XUserID = oldGame.OUserID
		game.XUserName = oldGame.OUserName
		game.XTgChatID = oldGame.OTgChatID
		game.OUserID = oldGame.XUserID
		game.OUserName = oldGame.XUserName
		game.OTgChatID = oldGame.XTgChatID
	}
	if game, err = Game.InsertGame(c, game.GameEntity); err != nil {
		return
	}
	return
}
