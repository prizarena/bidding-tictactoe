package commands

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/qedus/nds"
	"github.com/prizarena/bidding-tictactoe/server-go/bttt-trans"
	"github.com/prizarena/bidding-tictactoe/server-go/btttbot/common"
	"github.com/prizarena/bidding-tictactoe/server-go/btttdal"
	"github.com/prizarena/bidding-tictactoe/server-go/btttdelays"
	"github.com/prizarena/bidding-tictactoe/server-go/btttmodels"
	"github.com/strongo/bots-api-telegram"
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/log"
	"time"
	"strconv"
)

func processBid(whc bots.WebhookContext, gameID int64, bid int16, x, y int8) (m bots.MessageFromBot, err error) {
	c := whc.Context()
	log.Debugf(c, "processBid() => gameID: %d, bid: %d, target: %dx%d", gameID, bid, x, y)
	var (
		game       btttmodels.Game
		player     btttmodels.Player
		targetCell btttmodels.Cell
		playerName string
		user       btttmodels.AppUser
	)
	if user, err = btttdal.User.GetUserByID(c, whc.AppUserIntID()); err != nil {
		return
	}

	botID := whc.GetBotCode()

	if err = nds.RunInTransaction(c, func(tc context.Context) error {
		if game, err = btttdal.Game.GetGameByID(tc, gameID); err != nil {
			return err
		}
		//whc.ChatEntity().SetPreferredLanguage(game.Locale)
		if err = whc.SetLocale(game.Locale); err != nil {
			return err
		}

		if player = common.User2player(user.ID, game.GameEntity); player == btttmodels.NotPlayer {
			log.Warningf(c, "User is not a player for the game")
			return nil
		}

		//cell = btttmodels.Cell{X: , Y: int8(y), V: player}

		game.Logbook = game.Logbook.LogBid(player, int(bid))

		if x > 0 && y > 0 {
			game.Logbook = game.Logbook.LogTarget(player, int(x), int(y))
		}

		currentTurn := game.Logbook.CurrentTurn()

		processPlayer := func(userXO btttmodels.GamePlayerJson, userTurn btttmodels.GameMove) (btttmodels.GamePlayerJson, error) {
			if userTurn.Bid > userXO.Balance {
				playerName = userXO.Name
				return userXO, errors.WithMessage(errBidTooBig, fmt.Sprintf("%v: bid:%d > balance:%v", player, userTurn.Bid, userXO.Balance))
			}
			userXO.BidTime = time.Now()
			if !currentTurn.X.HasTarget() {
				return userXO, errors.New("player X has no target")
			}
			return userXO, nil
		}

		userX, userO := game.GetUsersXO()

		switch player {
		case btttmodels.PlayerX:
			if userX, err = processPlayer(userX, currentTurn.X); err != nil {
				return err
			}
		case btttmodels.PlayerO:
			if userO, err = processPlayer(userO, currentTurn.O); err != nil {
				return err
			}
			if userO.Name == "" { // TODO: Move inside processPlayer()?
				var botUser bots.BotAppUser
				if botUser, err = whc.GetAppUser(); err != nil {
					return err
				}
				user := botUser.(*btttmodels.AppUserEntity)
				userO.Name = user.GetFullName()
			}
		default:
			panic(fmt.Sprintf("Unknown player value: %v", player))
		}
		if currentTurn.HasBothBidsAndTargets() ||
			(userX.Balance == 0 && currentTurn.X.HasTarget() && currentTurn.O.HasBidAndTarget()) ||
			(userO.Balance == 0 && currentTurn.O.HasTarget() && currentTurn.X.HasBidAndTarget()) {
			xBid, oBid := currentTurn.X.Bid, currentTurn.O.Bid
			log.Debugf(c, "Game has bids and targets from both players: XBid=%d, OBid=%d", xBid, oBid)

			if xBid > oBid {
				targetCell = game.TargetCell(btttmodels.PlayerX)
			} else if oBid > xBid {
				targetCell = game.TargetCell(btttmodels.PlayerO)
			} else if xBid == oBid {
				if userX.Balance > userO.Balance {
					targetCell = game.TargetCell(btttmodels.PlayerX)
				} else if userO.Balance > userX.Balance {
					targetCell = game.TargetCell(btttmodels.PlayerO)
				} else if userX.BidTime.Before(userO.BidTime) {
					targetCell = game.TargetCell(btttmodels.PlayerX)
				} else if userO.BidTime.Before(userX.BidTime) {
					targetCell = game.TargetCell(btttmodels.PlayerO)
				} else {
					panic("Program logic error: XBid == OBid && OBalance == XBalance && XBidTime == OBidTime")
				}
			} else {
				panic("Program logic error")
			}

			if game.Board, err = game.Board.Turn(targetCell); err != nil {
				return errors.Wrap(err, "Failed to make a turn")
			}
			switch targetCell.V {
			case btttmodels.PlayerX:
				if xBid > userX.Balance {
					return errors.Wrap(errBidTooBig, fmt.Sprintf("game.XBid:%d > game.XBalance:%v", xBid, userX.Balance))
				}
				userX.Balance -= xBid
				userO.Balance += xBid
			case btttmodels.PlayerO:
				if oBid > userO.Balance {
					return errors.Wrap(errBidTooBig, fmt.Sprintf("game.OBid:%d > game.OBalance:%v", oBid, userO.Balance))
				}
				userO.Balance -= oBid
				userX.Balance += oBid

			default:
				panic(fmt.Sprintf("Unknown player: %v", player))
			}
			userX.BidTime = time.Time{}
			userO.BidTime = time.Time{}
			game.SetPlayers(userX, userO)
			game.CountOfTurns += 1
			game.Logbook = game.Logbook.SetTurnWinner(targetCell.V, game.Board.Winner() == btttmodels.NoWinnerYet)
			if game.TgInlineMessageID != "" {
				if err = btttdelays.DelayUpdateInlineMessage(tc, botID, gameID); err != nil {
					log.Errorf(c, errors.Wrap(err, "Not critical: Failed to queue update of inline message").Error())
					err = nil
				}
			}
		} else {
			log.Debugf(c, "Game is not ready to make the turn")
		}

		oldChatID, oldMessageID := game.UserTelegramData(user.ID)
		m = whc.NewMessage("")
		winner := game.Board.Winner()
		if m.Text, err = common.GameMessageText(c, whc, common.MODE_INBOT_NEW, game, winner, user); err != nil {
			return err
		}
		m.Keyboard = common.BoardToInlineKeyboard(c, whc, common.MODE_INBOT_NEW, game, winner, user.ID, botID)
		m.DisableWebPagePreview = true
		log.Debugf(c, "Sending new game message to Telegram over HTTPS request...")
		if tgResponse, err := whc.Responder().SendMessage(c, m, bots.BotAPISendMessageOverHTTPS); err != nil {
			err = errors.Wrap(err, "Failed to send game message to Telegram")
			log.Errorf(c, err.Error())
			return err
		} else {
			tgMessageID := strconv.Itoa(tgResponse.TelegramMessage.(tgbotapi.Message).MessageID) // TODO: Temporary, should be abstracted?
			playerJson := game.GetPlayerJsonByUserID(user.ID)
			playerJson.Tg.MessageID = tgMessageID
			game.SetPlayerJson(user.ID, playerJson)
		}

		if oldChatID != 0 && oldMessageID != 0 {
			if err = btttdelays.DelayDeleteTelegramMessage(tc, botID, oldChatID, oldMessageID); err != nil {
				log.Debugf(c, errors.Wrap(err, "Not critical: Failed to queue Telegram message for deletion").Error())
				err = nil
			}
		} else {
			log.Debugf(c, "oldChatID: %d, oldMessageID: %d", oldChatID, oldMessageID)
		}

		if err = btttdal.Game.SaveGame(tc, game); err != nil {
			return err
		}
		updateOtherMessages(tc, botID, game, user.ID, true)
		return err
	}, nil); err != nil {
		if errors.Cause(err) == errBidTooBig {
			m = whc.NewMessageByCode(bttt_trans.MT_BID_TOO_BIG, playerName, game.GetPlayerJsonByUserID(user.ID).Balance)
			err = nil
		}
		return m, err
	}

	if player == btttmodels.NotPlayer {
		return
	}

	m = bots.MessageFromBot{}
	return
}
