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
	"github.com/strongo/bots-framework/platforms/telegram"
	"github.com/strongo/log"
	"strconv"
	"strings"
	"time"
)

var errNoFreeSlots = errors.New("No free slots")

func startHit(whc bots.WebhookContext, mt string) (m bots.MessageFromBot, err error) {
	c := whc.Context()
	log.Debugf(c, "startHit() => mt: %v", mt)
	vals := strings.Split(mt, "-")
	var (
		turnRequest btttmodels.TurnRequest
		targetCell  btttmodels.Cell
	)
	if turnRequest.GameID, err = strconv.ParseInt(vals[1], 10, 64); err != nil {
		return
	}

	if turnRequest.X, err = strconv.Atoi(vals[2][0:1]); err != nil {
		return
	}
	if turnRequest.Y, err = strconv.Atoi(vals[2][1:2]); err != nil {
		return
	}

	log.Debugf(c, "turnRequest: %v", turnRequest)

	var game btttmodels.Game
	if game, err = btttdal.Game.GetGameByID(c, turnRequest.GameID); err != nil {
		return
	}

	log.Debugf(c, "Game loaded from DB")

	userID := whc.AppUserIntID()

	player := game.Player(userID)

	if player == btttmodels.NotPlayer {
		targetCell = btttmodels.Cell{X: int8(turnRequest.X), Y: int8(turnRequest.Y), V: player}
		sayNoFreeSlots := func() {
			log.Debugf(c, "User is not a player in the game and not free slots available")
			m = whc.NewMessage("Sorry, this game already have 2 players. Please start a new game.")
		}
		if game.HasFreeSlots() {
			log.Debugf(c, "User is not a player in the game with a free slot available")
			user := btttmodels.NewAppUser(whc.AppUserIntID())
			var botUser bots.BotAppUser
			if botUser, err = whc.GetAppUser(); err != nil {
				return
			}
			user.AppUserEntity = botUser.(*btttmodels.AppUserEntity)

			if err = nds.RunInTransaction(c, func(c context.Context) error {
				if game, err = btttdal.Game.GetGameByID(c, game.ID); err != nil {
					return err
				}
				initPlayerJson := func() (userXO btttmodels.GamePlayerJson) {
					userXO = game.GetPlayerJsonByUserID(userID)
					userXO.Balance = 100
					userXO.Name = user.FullName()
					return
				}
				switch int64(0) {
				case game.XUserID:
					player = btttmodels.PlayerX
					game.XUserID = userID
					game.SetPlayerX(initPlayerJson())
				case game.OUserID:
					player = btttmodels.PlayerO
					game.OUserID = userID
					game.SetPlayerO(initPlayerJson())
				default:
					return errNoFreeSlots
				}
				targetCell.V = player
				if err = updateGameWithTargetChatMessageID(whc, game, targetCell); err != nil {
					return err
				}
				return btttdal.Game.SaveGame(c, game)
			}, nil); err != nil {
				if err == errNoFreeSlots {
					sayNoFreeSlots()
					err = nil
				}
				return
			}
			m, err = askForBidAndUpdateOthers(whc, game)
		} else {
			sayNoFreeSlots()
			return
		}
	} else {
		log.Debugf(c, "User already player in the game")
		m, err = processHitForAlreadyJoinedPlayer(whc, turnRequest)
	}

	return
}

var errCellOccupied = errors.New("Cell is already occupied")

func processHitForAlreadyJoinedPlayer(whc bots.WebhookContext, turnRequest btttmodels.TurnRequest) (m bots.MessageFromBot, err error) {
	c := whc.Context()
	log.Debugf(c, "processHitForAlreadyJoinedPlayer() => turnRequest: %v", turnRequest)
	if !turnRequest.Valid() {
		err = errors.New("Can't process bid as turnRequest is not valid")
		return
	}
	userID := whc.AppUserIntID()
	var game btttmodels.Game
	if err = nds.RunInTransaction(c, func(c context.Context) error {
		if game, err = btttdal.Game.GetGameByID(c, turnRequest.GameID); err != nil {
			return err
		}
		cellValue := game.Board.CellValue(turnRequest.X, turnRequest.Y)
		switch cellValue {
		case btttmodels.CellEmpty:
			// Is empty as expected
		case btttmodels.CellX, btttmodels.CellO:
			return errCellOccupied
		default:
			panic(fmt.Sprintf("Unknown cell value: %v", cellValue))
		}
		if chatID, tgMessageID := game.UserTelegramData(userID); chatID != 0 && tgMessageID != 0 {
			if err = btttdelays.DelayDeleteTelegramMessage(c, whc.GetBotCode(), chatID, tgMessageID); err != nil {
				log.Errorf(c, "Not critical: Failed to queue Telegram message for deletion")
			}
		}
		player := game.Player(userID)
		if err = updateGameWithTargetChatMessageID(whc, game, btttmodels.Cell{X: int8(turnRequest.X), Y: int8(turnRequest.Y), V: player}); err != nil {
			return err
		}
		return btttdal.Game.SaveGame(c, game)
	}, nil); err != nil {
		if err == errCellOccupied {
			if _, isCallbackQuery := whc.Input().(bots.WebhookCallbackQuery); isCallbackQuery {
				m.BotMessage = telegram.CallbackAnswer(tgbotapi.AnswerCallbackQueryConfig{
					Text:      "🚫 " + whc.Translate(bttt_trans.MT_CELL_OCCUPIED),
					ShowAlert: true,
					CacheTime: 60,
				})
			} else {
				m = whc.NewMessageByCode(bttt_trans.MT_CELL_OCCUPIED)
				m.Text = "🚫 " + m.Text
				log.Errorf(c, "TODO: MT_CELL_OCCUPIED - Display board")
				return
			}
			err = nil
		}
		return
	}
	m, err = askForBidAndUpdateOthers(whc, game)
	log.Debugf(c, "askForBidAndUpdateOthers() => m.Text: %v", m.Text)
	return
}

func askForBidAndUpdateOthers(whc bots.WebhookContext, game btttmodels.Game) (m bots.MessageFromBot, err error) {
	c := whc.Context()
	log.Debugf(c, "askForBidAndUpdateOthers()")
	whc.ChatEntity().SetAwaitingReplyTo("bid?g=" + strconv.FormatInt(game.ID, 10))

	if err = updateOtherMessages(whc.Context(), whc.GetBotCode(), game, whc.AppUserIntID(), true); err != nil {
		log.Errorf(c, errors.Wrap(err, "Not critical: Failed to update other messages").Error())
		err = nil
	}

	m = whc.NewMessageByCode(bttt_trans.MT_ASK_BID)
	m.Keyboard = tgbotapi.ForceReply{ForceReply: true, Selective: true}
	return
}

func updateOtherMessages(c context.Context, botID string, game btttmodels.Game, currentUserID int64, updateInline bool) (err error) {
	log.Debugf(c, "updateOtherMessages(game.ID=%d, currentUserID=%d, updateInline=%v)", game.ID, currentUserID, updateInline)
	var userID int64
	switch currentUserID {
	case game.XUserID:
		userID = game.OUserID
	case game.OUserID:
		userID = game.XUserID
	default:
		panic("User ID does not belong to the game")
	}

	if userID != 0 {
		userXO := game.GetPlayerJsonByUserID(userID)
		if userXO.Tg.ChatID != "" {
			if err = btttdelays.DelayUpdateInBotMessage(c, botID, game.ID, userID); err != nil {
				return err
			}
		}
	}

	if updateInline && game.TgInlineMessageID != "" {
		if err = btttdelays.DelayUpdateInlineMessage(c, botID, game.ID); err != nil {
			return err
		}
	}

	return nil
}

func updateGameWithTargetChatMessageID(whc bots.WebhookContext, game btttmodels.Game, targetCell btttmodels.Cell) (err error) {
	c := whc.Context()
	log.Debugf(c, "updateGameWithTargetChatMessageID(), targetCell=%+v", targetCell)
	var playerUserID int64
	switch targetCell.V {
	case btttmodels.PlayerX:
		playerUserID = game.XUserID
	case btttmodels.PlayerO:
		playerUserID = game.OUserID
	default:
		return fmt.Errorf("updateGameWithTargetChatMessageID: target cell references unknown player: %v", targetCell.V)
	}
	if playerUserID == 0 {
		return fmt.Errorf("game %v has no user ID for player %v", game.ID, targetCell.V)
	}

	if targetCell.X > 0 && targetCell.Y > 0 {
		game.Logbook = game.Logbook.LogTarget(targetCell.V, int(targetCell.X), int(targetCell.Y))
		playerJson := game.GetPlayerJsonByUserID(playerUserID)
		playerJson.BidTime = time.Time{}
		game.SetPlayerJson(playerUserID, playerJson)
	}

	var m bots.MessageFromBot
	var botUser bots.BotAppUser
	if botUser, err = whc.GetAppUser(); err != nil {
		return
	}
	user := btttmodels.NewAppUser(whc.AppUserIntID())
	user.AppUserEntity = botUser.(*btttmodels.AppUserEntity)

	m, err = common.GameNewMessageToBot(whc, common.MODE_INBOT_NEW, game, game.Board.Winner(), user)

	var response bots.OnMessageSentResponse
	if response, err = whc.Responder().SendMessage(whc.Context(), m, bots.BotAPISendMessageOverHTTPS); err != nil {
		return err
	}
	tgMessage := response.TelegramMessage.(tgbotapi.Message)
	{ // setTelegramInfo
		p := game.GetPlayerJsonByUserID(playerUserID)
		p.Tg.ChatID = strconv.FormatInt(tgMessage.Chat.ID, 10)
		p.Tg.MessageID = strconv.Itoa(tgMessage.MessageID)
		game.SetPlayerJson(playerUserID, p)
	}
	return nil
}
