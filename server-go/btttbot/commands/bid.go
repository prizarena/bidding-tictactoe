package commands

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/qedus/nds"
	"github.com/strongo/bidding-tictactoe-bot/bttt-trans"
	"github.com/strongo/bidding-tictactoe-bot/btttbot/common"
	"github.com/strongo/bidding-tictactoe-bot/btttdal"
	"github.com/strongo/bidding-tictactoe-bot/btttdelays"
	"github.com/strongo/bidding-tictactoe-bot/btttmodels"
	"github.com/strongo/bots-api-telegram"
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/log"
	"time"
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

		switch player {
		case btttmodels.PlayerX:
			if currentTurn.X.Bid > game.XBalance {
				playerName = game.XUserName
				return errors.WithMessage(errBidTooBig, fmt.Sprintf("XBid:%d > game.XBalance:%v", currentTurn.X.Bid, game.XBalance))
			}
			game.XBidTime = time.Now()
			if !currentTurn.X.HasTarget() {
				return errors.New("player X has no target")
			}
		case btttmodels.PlayerO:
			if currentTurn.O.Bid > game.OBalance {
				playerName = game.OUserName
				return errors.WithMessage(errBidTooBig, fmt.Sprintf("game.OBid:%d > game.OBalance:%v", currentTurn.O.Bid, game.OBalance))
			}
			game.OBidTime = time.Now()
			if !currentTurn.O.HasTarget() {
				return errors.New("player O has no target")
			}

			if game.OUserName == "" {
				var botUser bots.BotAppUser
				if botUser, err = whc.GetAppUser(); err != nil {
					return err
				}
				user := botUser.(*btttmodels.AppUserEntity)
				game.OUserName = user.FullName()
			}
		default:
			panic(fmt.Sprintf("Unknown player value: %v", player))
		}
		if currentTurn.HasBothBidsAndTargets() ||
			(game.XBalance == 0 && currentTurn.X.HasTarget() && currentTurn.O.HasBidAndTarget()) ||
			(game.OBalance == 0 && currentTurn.O.HasTarget() && currentTurn.X.HasBidAndTarget()) {
			xBid, oBid := currentTurn.X.Bid, currentTurn.O.Bid
			log.Debugf(c, "Game has bids and targets from both players: XBid=%d, OBid=%d", xBid, oBid)
			if xBid > oBid {
				targetCell = game.TargetCell(btttmodels.PlayerX)
			} else if oBid > xBid {
				targetCell = game.TargetCell(btttmodels.PlayerO)
			} else if xBid == oBid {
				if game.XBalance > game.OBalance {
					targetCell = game.TargetCell(btttmodels.PlayerX)
				} else if game.OBalance > game.XBalance {
					targetCell = game.TargetCell(btttmodels.PlayerO)
				} else if game.XBidTime.Before(game.OBidTime) {
					targetCell = game.TargetCell(btttmodels.PlayerX)
				} else if game.OBidTime.Before(game.XBidTime) {
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
				if xBid > game.XBalance {
					return errors.Wrap(errBidTooBig, fmt.Sprintf("game.XBid:%d > game.XBalance:%v", xBid, game.XBalance))
				}
				game.XBalance -= xBid
				game.OBalance += xBid
			case btttmodels.PlayerO:
				if oBid > game.OBalance {
					return errors.Wrap(errBidTooBig, fmt.Sprintf("game.OBid:%d > game.OBalance:%v", oBid, game.OBalance))
				}
				game.OBalance -= oBid
				game.XBalance += oBid
			default:
				panic(fmt.Sprintf("Unknown player: %v", player))
			}
			game.XBidTime = time.Time{}
			game.OBidTime = time.Time{}
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
			tgMessageID := tgResponse.TelegramMessage.(tgbotapi.Message).MessageID // TODO: Temporary, should be abstracted?
			switch player {
			case btttmodels.PlayerX:
				game.XTgMessageID = tgMessageID
			case btttmodels.PlayerO:
				game.OTgMessageID = tgMessageID
			default:
				panic("Unknown player")
			}
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
			m = whc.NewMessageByCode(bttt_trans.MT_BID_TOO_BIG, playerName, game.UserBalance(user.ID))
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
