package common

import (
	"context"
	"fmt"
	"github.com/strongo/app"
	"github.com/strongo-games/bidding-tictactoe/server-go/bttt-trans"
	"github.com/strongo-games/bidding-tictactoe/server-go/btttmodels"
	"github.com/strongo/bots-api-telegram"
	"github.com/strongo/log"
	"html"
	"strconv"
)

func BoardToInlineKeyboard(c context.Context, translator strongo.SingleLocaleTranslator, mode Mode, game btttmodels.Game, winner btttmodels.Player, currentUserID int64, botID string) *tgbotapi.InlineKeyboardMarkup {
	log.Debugf(c, "BoardToInlineKeyboard() => mode: %v, currentUserID: %d, bot: %v", mode, currentUserID, botID)

	var keyboard [][]tgbotapi.InlineKeyboardButton

	if winner == btttmodels.NoWinnerYet {
		keyboard = ongoingGameKeyboardForTelegram(game, mode, currentUserID, botID)
	} else {
		newGame := "new_game"
		var newGameThisChatButton tgbotapi.InlineKeyboardButton
		switch mode {
		case MODE_INLINE:
			newGameThisChatButton = tgbotapi.NewInlineKeyboardButtonSwitchInlineQueryCurrentChat(
				translator.Translate(bttt_trans.C_NEW_GAME_THIS_CHAT),
				newGame,
			)
		case MODE_INBOT_NEW, MODE_INBOT_EDIT:
			var opponentName string
			switch currentUserID {
			case game.XUserID:
				opponentName = game.OUserName
			case game.OUserID:
				opponentName = game.XUserName
			default:
				panic("Unknown user")
			}
			newGameThisChatButton = tgbotapi.InlineKeyboardButton{
				Text:         translator.Translate(bttt_trans.C_NEW_GAME_WITH, html.EscapeString(opponentName)),
				CallbackData: newGame + fmt.Sprintf("?g=%d", game.ID),
			}
		default:

		}
		keyboard = [][]tgbotapi.InlineKeyboardButton{
			{
				newGameThisChatButton,
			},
			{
				{
					Text:              translator.Translate(bttt_trans.C_NEW_GAME_OTHER_CHAT),
					SwitchInlineQuery: &newGame,
				},
			},
		}
		if winner == btttmodels.IsTie {
		} else if winner == btttmodels.PlayerX || winner == btttmodels.PlayerO {
		} else {
			panic("Invalid programming")
		}
	}

	return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboard}
}

func ongoingGameKeyboardForTelegram(game btttmodels.Game, mode Mode, currentUserID int64, botID string) [][]tgbotapi.InlineKeyboardButton {
	boardSize := game.Board.Size()
	var keyboard [][]tgbotapi.InlineKeyboardButton
	keyboard = make([][]tgbotapi.InlineKeyboardButton, boardSize)

	callbackStr := "hit?c=%d%d&g=" + strconv.FormatInt(game.ID, 10)
	urlPrefix := fmt.Sprintf("https://t.me/%v?start=hit-%d-", botID, game.ID)

	btnAction := func(x, y int8, btn tgbotapi.InlineKeyboardButton) tgbotapi.InlineKeyboardButton {
		if mode == MODE_INLINE {
			btn.URL = urlPrefix + fmt.Sprintf("%d%d", x, y)
		} else {
			btn.CallbackData = fmt.Sprintf(callbackStr, x, y)
		}
		return btn
	}

	player := game.Player(currentUserID)

	var targetCell btttmodels.Cell
	if player != btttmodels.NotPlayer {
		targetCell = game.TargetCell(player)
	}

	for row, gridRow := range game.Board.Grid() {
		kbRow := make([]tgbotapi.InlineKeyboardButton, boardSize)
		for col, cellV := range gridRow {
			x, y := int8(col+1), int8(row+1)
			if cellV == btttmodels.CellEmpty {
				if mode == MODE_INBOT_EDIT || mode == MODE_INBOT_NEW {
					var text string
					if x == targetCell.X && y == targetCell.Y {
						if game.HasBid(player) {
							text = "💰"
						} else {
							text = "📍"
						}
					} else {
						text = " "
					}
					kbRow[col] = btnAction(x, y, tgbotapi.InlineKeyboardButton{
						Text: text,
					})
				} else {
					kbRow[col] = btnAction(x, y, tgbotapi.InlineKeyboardButton{
						Text: " ",
					})
				}
			} else {
				kbRow[col] = tgbotapi.InlineKeyboardButton{
					Text:         btttmodels.DrawPlayerToCell(cellV),
					CallbackData: fmt.Sprintf(callbackStr, x, y),
				}
			}
		}
		keyboard[row] = kbRow
	}
	return keyboard
}
