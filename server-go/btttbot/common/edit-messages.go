package common

import (
	"bytes"
	"context"
	"fmt"
	"github.com/strongo/app"
	"github.com/strongo-games/bidding-tictactoe/server-go/bttt-trans"
	"github.com/strongo-games/bidding-tictactoe/server-go/btttmodels"
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/log"
	"strings"
)

func User2player(userID int64, gameEntity *btttmodels.GameEntity) btttmodels.Player {
	if userID == 0 {
		panic("userID == 0")
	}
	if gameEntity == nil {
		panic("gameEntity == nil")
	}
	switch userID {
	case gameEntity.XUserID:
		return btttmodels.PlayerX
	case gameEntity.OUserID:
		return btttmodels.PlayerO
	default:
		switch int64(0) {
		case gameEntity.XUserID:
			gameEntity.XUserID = userID
			return btttmodels.PlayerX
		case gameEntity.OUserID:
			gameEntity.OUserID = userID
			return btttmodels.PlayerO
		default:
			return btttmodels.NotPlayer
		}
	}
}

// translator := strongo.NewSingleMapTranslator(game.GetLocale(), strongo.NewMapTranslator(c, bttt_trans.TRANS))

func GameNewMessageToBot(whc bots.WebhookContext, mode Mode, game btttmodels.Game, winner btttmodels.Player, currentUser btttmodels.AppUser) (m bots.MessageFromBot, err error) {
	c := whc.Context()
	m = whc.NewMessage("")
	if m.Text, err = GameMessageText(c, whc, mode, game, winner, currentUser); err != nil {
		return
	}
	m.Keyboard = BoardToInlineKeyboard(c, whc, mode, game, winner, currentUser.ID, whc.GetBotCode())
	m.DisableWebPagePreview = true
	return
}

func GameMessageText(c context.Context, translator strongo.SingleLocaleTranslator, mode Mode, game btttmodels.Game, winner btttmodels.Player, currentUser btttmodels.AppUser) (messageText string, err error) {
	log.Debugf(c, "GameMessageText() => mode=%v, game.ID: %d, currentUser.ID: %d", mode, game.ID, currentUser.ID)
	if mode != MODE_INBOT_NEW && mode != MODE_INBOT_EDIT && mode != MODE_INLINE {
		panic(fmt.Sprintf("Unknown mode: %d", mode))
	}
	if mode != MODE_INLINE && currentUser.ID == 0 {
		panic("mode != MODE_INLINE && currentUser.ID == 0")
	}
	boardSize := game.Board.Size()

	var hostAndPath = "biddingtictactoe.com/"
	if translator.Locale().Code5 == "ru-RU" {
		hostAndPath += "ru"
	}
	mt := new(bytes.Buffer)
	fmt.Fprintf(mt, `<b>%v</b> %dx%d — <a href="https://%v#rules">%v</a>`, translator.Translate(bttt_trans.MT_GAME_NAME), boardSize, boardSize, hostAndPath, translator.Translate(bttt_trans.MT_RULES))
	if mode == MODE_INLINE && game.Board.IsEmpty() {
		mt.WriteString("\n<i>" + translator.Translate(bttt_trans.MT_HOW_TO_INLINE) + "</i>")
		mt.WriteString("\n──────────────")
	}
	const padding = "\n  "

	writePlayer := func(p string, gameUserID int64, gamePlayerJson btttmodels.GamePlayerJson, bid, previousBid int, hasTarget bool) {
		mt.WriteString("\n" + translator.Translate(bttt_trans.MT_PLAYER, p))
		if gamePlayerJson.Name == "" {
			mt.WriteString(" " + translator.Translate(bttt_trans.MT_AWAITING_PLAYER))
		} else {
			mt.WriteString(" " + gamePlayerJson.Name)
		}
		mt.WriteString(padding + translator.Translate(bttt_trans.MT_PLAYER_BALANCE, gamePlayerJson.Balance))
		hasBid := bid > 0
		if gameUserID == currentUser.ID {
			if hasBid {
				mt.WriteString("; " + translator.Translate(bttt_trans.MT_YOUR_BID, bid))
			}
		} else {
			if hasTarget && hasBid {
				mt.WriteString("; " + translator.Translate(bttt_trans.MT_RIVAL_HAS_TARGET_AND_BID))
			} else if hasTarget {
				mt.WriteString("; " + translator.Translate(bttt_trans.MT_RIVAL_HAS_TARGET))
			} else if hasBid {
				mt.WriteString("; " + translator.Translate(bttt_trans.MT_RIVAL_HAS_BID))
			}
		}
		if previousBid > 0 {
			if bid > 0 {
				if gameUserID != currentUser.ID {
					mt.WriteString(padding + translator.Translate(bttt_trans.MT_PREVIOUS_BID, previousBid))
				}
			} else {
				mt.WriteString(padding + translator.Translate(bttt_trans.MT_LAST_BID, previousBid))
			}
		}
	}

	currentTurn := game.Logbook.CurrentTurn()
	previousTurn := game.Logbook.PreviousTurn()
	writePlayer("X", game.XUserID, game.PlayerX(), currentTurn.X.Bid, previousTurn.X.Bid, game.HasTarget(btttmodels.PlayerX))
	writePlayer("O", game.OUserID, game.PlayerO(), currentTurn.O.Bid, previousTurn.O.Bid, game.HasTarget(btttmodels.PlayerO))

	if winner == btttmodels.NoWinnerYet {
		if currentUser.ID > 0 {
			player := game.Player(currentUser.ID)
			hasTarget := game.HasTarget(player)
			hasBid := game.HasBid(player)
			if !hasTarget {
				mt.WriteString("\n<b>" + translator.Translate(bttt_trans.MT_PLEASE_CHOOSE_YOUR_TARGET) + "</b>")
			} else if !hasBid {
				mt.WriteString("\n<b>" + translator.Translate(bttt_trans.MT_PLEASE_MAKE_A_BID) + "</b>")
			} else {
				rival := game.Rival(player)
				if !game.HasTarget(rival) || !game.HasBid(rival) {
					mt.WriteString("\n" + translator.Translate(bttt_trans.MT_AWAITING_RIVAL_TURN, fmt.Sprintf("<b>%v</b>", btttmodels.DrawPlayerToCell(rival))))
				}
			}
		}
	} else {
		log.Debugf(c, "winner: %v", string([]byte{byte(winner)}))
		mt.WriteString("\n\n")
		if err = writeTgFooterForWinner(c, mt, currentUser, translator, game); err != nil {
			return
		}
	}

	if mode == MODE_INLINE || winner != btttmodels.NoWinnerYet {
		sponsoredBy := "DebtsTrackerBot"
		if translator.Locale().Code5 == strongo.LocalCodeRuRu {
			sponsoredBy = "DebtsTrackerRuBot"
		}
		mt.WriteString("\n" + strings.Repeat("─", 10) + "\n" + translator.Translate(bttt_trans.MT_FREE_GAME_SPONSORED_BY, fmt.Sprintf(`<a href="https://t.me/%v?start=ref-BiddingTicTacToeBot">@%v</a>`, sponsoredBy, sponsoredBy)))
	}

	return mt.String(), err
}
