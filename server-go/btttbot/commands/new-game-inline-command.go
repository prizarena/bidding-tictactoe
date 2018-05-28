package commands

import (
	"fmt"
	"github.com/strongo/app"
	"github.com/strongo-games/bidding-tictactoe/server-go/bttt-trans"
	"github.com/strongo/bots-api-telegram"
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/bots-framework/platforms/telegram"
	"github.com/strongo/log"
)

var (
	newGameEnResult, newGameRuResult tgbotapi.InlineQueryResultArticle
	_newGameResults                  []interface{}
)

func StartNewGame(whc bots.WebhookContext) (m bots.MessageFromBot, err error) {
	log.Debugf(whc.Context(), "StartNewGame()")
	inlineQuery := whc.Input().(bots.WebhookInlineQuery)

	m.BotMessage = telegram.InlineBotMessage(tgbotapi.InlineConfig{
		InlineQueryID: inlineQuery.GetInlineQueryID(),
		CacheTime:     60,
		Results:       newGameResults(whc),
	})
	return m, err
}

const howToPlayLink = `<a href="https://biddingtictactoe.com/#how-to-play">How to play</a>`

func newGameResults(whc bots.WebhookContext) []interface{} {
	c := whc.Context()

	translator := strongo.NewSingleMapTranslator(strongo.LocaleEnUS, strongo.NewMapTranslator(c, bttt_trans.TRANS))
	if newGameEnResult.ID == "" {
		textFormat := "<b>%v</b>\n%v"
		enText := fmt.Sprintf(textFormat, translator.Translate(bttt_trans.MT_GAME_NAME), translator.Translate(bttt_trans.MT_NEW_GAME_WELCOME))
		newGameEnResult = tgbotapi.InlineQueryResultArticle{
			ID:          "new-game-en",
			Type:        "article",
			Title:       "🇺🇸 New game",
			Description: "Create a new board",
			InputMessageContent: tgbotapi.InputTextMessageContent{
				Text:                  enText,
				ParseMode:             "HTML",
				DisableWebPagePreview: true,
			},
			ReplyMarkup: &tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
					{{Text: "Board 3x3", CallbackData: "board?3en"}},
					//{{Text: "Board 4x4", CallbackData: "board?4en"}},
					//{{Text: "Board 5x5", CallbackData: "board?5en"}},
				},
			},
		}

		translator = strongo.NewSingleMapTranslator(strongo.LocaleRuRu, strongo.NewMapTranslator(c, bttt_trans.TRANS))
		ruText := fmt.Sprintf(textFormat, translator.Translate(bttt_trans.MT_GAME_NAME), translator.Translate(bttt_trans.MT_NEW_GAME_WELCOME))
		newGameRuResult = tgbotapi.InlineQueryResultArticle{
			ID:          "new-game-ru",
			Type:        "article",
			Title:       "🇷🇺 Новая игра",
			Description: "Создать новую доску",
			InputMessageContent: tgbotapi.InputTextMessageContent{
				Text:                  ruText,
				ParseMode:             "HTML",
				DisableWebPagePreview: true,
			},
			ReplyMarkup: &tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
					{{Text: "Доска 3x3", CallbackData: "board?3ru"}},
					//{{Text: "Доска 4x4", CallbackData: "board?4ru"}},
					//{{Text: "Доска 5x5", CallbackData: "board?5ru"}},
				},
			},
		}
		_newGameResults = []interface{}{newGameEnResult, newGameRuResult}
	}
	return _newGameResults
}
