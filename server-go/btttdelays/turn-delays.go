package btttdelays

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/qedus/nds"
	"github.com/strongo/app"
	"github.com/strongo/app/gae"
	"github.com/strongo/bidding-tictactoe-bot/bttt-trans"
	"github.com/strongo/bidding-tictactoe-bot/btttbot/common"
	"github.com/strongo/bidding-tictactoe-bot/btttbot/platforms/tgbots"
	"github.com/strongo/bidding-tictactoe-bot/btttdal"
	"github.com/strongo/bidding-tictactoe-bot/btttmodels"
	"github.com/strongo/bots-api-telegram"
	"github.com/strongo/log"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/delay"
	"google.golang.org/appengine/urlfetch"
	"net/http"
	"strings"
)

func DelayUpdateInBotMessage(c context.Context, botID string, gameID, userID int64) error {
	return gae.CallDelayFunc(c, QUEUE_TURN, "update-inbot-message", delayedUpdateInBotMessage, botID, gameID, userID)
}

var delayedUpdateInBotMessage = delay.Func("UpdateInBotMessage", func(c context.Context, botID string, gameID, userID int64) (err error) {
	log.Debugf(c, "delayedUpdateInBotMessage() => gameID=%d, userID=%v, botID=%v", gameID, userID, botID)
	if botID == "" {
		log.Criticalf(c, "botID is empty string")
		return
	}
	if gameID == 0 {
		log.Criticalf(c, "gameID == 0")
		return
	}
	if userID == 0 {
		log.Criticalf(c, "userID == 0")
		return
	}
	var (
		game     btttmodels.Game
		baseEdit tgbotapi.BaseEdit
	)
	if game, err = btttdal.Game.GetGameByID(c, gameID); err != nil {
		if errors.Cause(err) == datastore.ErrNoSuchEntity {
			log.Errorf(c, errors.Wrap(err, "Game not found by ID").Error())
			err = nil
		}
		return
	}
	switch userID {
	case game.XUserID:
		baseEdit.ChatID = game.XTgChatID
		baseEdit.MessageID = game.XTgMessageID
	case game.OUserID:
		baseEdit.ChatID = game.OTgChatID
		baseEdit.MessageID = game.OTgMessageID
	default:
		log.Errorf(c, fmt.Sprintf("User %d does not belong to the game %d", userID, gameID))
		return
	}

	if baseEdit.ChatID != 0 {
		if baseEdit.MessageID == 0 {
			if err = newGameTelegramMessage(c, botID, baseEdit.ChatID, game, userID); err != nil {
				return
			}
		} else {
			if err = updateGameTelegramMessage(c, botID, baseEdit, common.MODE_INBOT_EDIT, game, userID); err != nil {
				return
			}
		}
	} else if baseEdit.MessageID == 0 {
		log.Debugf(c, "Nothing to update as user has no inbot message yet.")
	} else if baseEdit.MessageID != 0 {
		log.Criticalf(c, "Data integrity issue: baseEdit.ChatID == 0 && baseEdit.MessageID != 0")
	} else {
		panic("Program logic error")
	}
	return
})

func DelayUpdateInlineMessage(c context.Context, botID string, gameID int64) error {
	return gae.CallDelayFunc(c, QUEUE_TURN, "update-inline-message", delayedUpdateInlineMessage, botID, gameID)
}

var delayedUpdateInlineMessage = delay.Func("UpdateInlineMessage", func(c context.Context, botID string, gameID int64) (err error) {
	log.Debugf(c, "delayedUpdateInlineMessage() => gameID=%d, botID=%v", gameID, botID)
	var game btttmodels.Game
	if game, err = btttdal.Game.GetGameByID(c, gameID); err != nil {
		return
	}
	if game.TgInlineMessageID == "" {
		log.Warningf(c, "game.TgInlineMessageID is empty string")
		return
	}
	if err = updateGameTelegramMessage(c, botID, tgbotapi.BaseEdit{InlineMessageID: game.TgInlineMessageID}, common.MODE_INLINE, game, 0); err != nil {
		if err.Error() == "Bad Request: message is not modified" {
			log.Errorf(c, errors.Wrap(err, "Failed to update inline message").Error())
			err = nil
		}
		return
	}
	return
})

func updateGameTelegramMessage(c context.Context, botID string, baseEdit tgbotapi.BaseEdit, mode common.Mode, game btttmodels.Game, currentUserID int64) (err error) {
	if baseEdit.InlineMessageID == "" && currentUserID == 0 {
		panic("baseEdit.InlineMessageID is empty string && currentUserID == 0")
	}
	if baseEdit.InlineMessageID != "" && currentUserID != 0 {
		panic("baseEdit.InlineMessageID is NOT empty string && currentUserID != 0")
	}
	translator := strongo.NewSingleMapTranslator(game.GetLocale(), strongo.NewMapTranslator(c, bttt_trans.TRANS))

	var mt string
	winner := game.Board.Winner()
	if mt, err = common.GameMessageText(c, translator, mode, game, winner, btttmodels.NewAppUser(currentUserID)); err != nil {
		return
	}

	if mode == common.MODE_INLINE && baseEdit.InlineMessageID == "" {
		if game.TgInlineMessageID == "" {
			return errors.New("Unknown Telegram inline message ID for the game")
		}
		baseEdit.InlineMessageID = game.TgInlineMessageID
	}

	editMessageTextConfig := &tgbotapi.EditMessageTextConfig{
		BaseEdit:              baseEdit,
		Text:                  mt,
		ParseMode:             "HTML",
		DisableWebPagePreview: true,
	}
	editMessageTextConfig.ReplyMarkup = common.BoardToInlineKeyboard(c, translator, mode, game, winner, currentUserID, botID)

	if botSettings, ok := tgbots.BotsBy(c).ByCode[botID]; !ok {
		logBotIsUnknown(c, botID, "updateGameTelegramMessage")
		return
	} else {
		tgBotApi := tgbotapi.NewBotAPIWithClient(botSettings.Token, &http.Client{Transport: &urlfetch.Transport{Context: c}})
		tgBotApi.EnableDebug(c)
		if _, err = tgBotApi.Send(editMessageTextConfig); err != nil {
			errMsg := err.Error()
			if strings.Contains(errMsg, "Bad Request: message is not modified") {
				log.Debugf(c, errors.WithMessage(err, "Failed to update Telegram message").Error())
				err = nil
				return
			} else if strings.Contains(errMsg, "Bad Request:") {
				log.Errorf(c, errors.WithMessage(err, "Failed to update Telegram message").Error())
				return
			}
			return
		}
	}

	return
}

func newGameTelegramMessage(c context.Context, botID string, chatID int64, game btttmodels.Game, currentUserID int64) (err error) {
	if botID == "" {
		panic("botID is empty string")
	}
	if chatID == 0 {
		panic("chatID == 0")
	}
	if currentUserID == 0 {
		panic("currentUserID == 0")
	}

	translator := strongo.NewSingleMapTranslator(game.GetLocale(), strongo.NewMapTranslator(c, bttt_trans.TRANS))

	var mt string
	winner := game.Board.Winner()
	if mt, err = common.GameMessageText(c, translator, common.MODE_INBOT_NEW, game, winner, btttmodels.NewAppUser(currentUserID)); err != nil {
		return
	}

	messageConfig := &tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      chatID,
			ReplyMarkup: common.BoardToInlineKeyboard(c, translator, common.MODE_INBOT_NEW, game, winner, currentUserID, botID),
		},
		Text:                  mt,
		ParseMode:             "HTML",
		DisableWebPagePreview: true,
	}

	if botSettings, ok := tgbots.BotsBy(c).ByCode[botID]; !ok {
		logBotIsUnknown(c, botID, "newGameTelegramMessage")
		return
	} else {
		tgBotApi := tgbotapi.NewBotAPIWithClient(botSettings.Token, &http.Client{Transport: &urlfetch.Transport{Context: c}})
		tgBotApi.EnableDebug(c)
		var response tgbotapi.Message
		if response, err = tgBotApi.Send(messageConfig); err != nil {
			errMsg := err.Error()
			if strings.Contains(errMsg, "Bad Request: message is not modified") {
				log.Warningf(c, errors.Wrap(err, "Failed to update Telegram message").Error())
				err = nil
				return
			} else if strings.Contains(errMsg, "Bad Request:") {
				log.Errorf(c, errors.Wrap(err, "Failed to update Telegram message").Error())
				return
			}
			return
		}
		if err = nds.RunInTransaction(c, func(c context.Context) error {
			if game, err = btttdal.Game.GetGameByID(c, game.ID); err != nil {
				return err
			}
			changed := false
			switch currentUserID {
			case game.XUserID:
				if game.XTgMessageID == 0 {
					game.XTgMessageID = response.MessageID
					changed = true
				}
			case game.OUserID:
				if game.OTgMessageID == 0 {
					game.OTgMessageID = response.MessageID
					changed = true
				}
			default:
				panic(fmt.Sprintf("Unknown currentUserID: %d", currentUserID))
			}
			if changed {
				if err = btttdal.Game.SaveGame(c, game); err != nil {
					return err
				}
			}
			return nil
		}, nil); err != nil {
			return
		}
	}

	return
}

func DelayDeleteTelegramMessage(c context.Context, botID string, chatID int64, messageID int) error {
	return gae.CallDelayFunc(c, QUEUE_TURN, "delete-tg-message", delayedDeleteTelegramMessage, botID, chatID, messageID)
}

var delayedDeleteTelegramMessage = delay.Func("DeleteTelegramMessage", func(c context.Context, botID string, chatID int64, messageID int) (err error) {
	log.Debugf(c, "delayedDeleteTelegramMessage() => bot: %v, chat: %d, message: %d", botID, chatID, messageID)
	if chatID == 0 {
		log.Errorf(c, "chatID == 0")
		return
	}
	if messageID == 0 {
		log.Errorf(c, "chatID == 0")
		return
	}
	if botSettings, ok := tgbots.BotsBy(c).ByCode[botID]; !ok {
		logBotIsUnknown(c, botID, "delayedDeleteTelegramMessage")
	} else {
		tgBotApi := tgbotapi.NewBotAPIWithClient(botSettings.Token, &http.Client{Transport: &urlfetch.Transport{Context: c}})
		if _, err = tgBotApi.Send(tgbotapi.DeleteMessage{ChatID: chatID, MessageID: messageID}); err != nil {
			log.Errorf(c, errors.WithMessage(err, "failed to delete Telegram message").Error())
			switch err.(type) {
			case tgbotapi.APIResponse:
				tgError := err.(tgbotapi.APIResponse)
				switch tgError.ErrorCode {
				case 400:
					if strings.Contains(tgError.Description, "Bad Request") {
						err = nil
					}
				}
			}
			return
		}
	}
	return
})


func logBotIsUnknown(c context.Context, botID, source string) {
	botsBy := tgbots.BotsBy(c)
	knownBots := make([]string, 0, len(botsBy.ByCode))
	for code, _ := range botsBy.ByCode {
		knownBots = append(knownBots, code)
	}
	log.Errorf(c, "%v: unknown bot ID: %v, known (%v): %v", source, botID, len(knownBots), knownBots)
}