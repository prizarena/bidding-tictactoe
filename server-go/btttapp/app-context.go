package btttapp

import (
	"context"
	"github.com/pkg/errors"
	"github.com/strongo/app"
	"github.com/strongo/bidding-tictactoe-bot/bttt-trans"
	"github.com/strongo/bidding-tictactoe-bot/btttmodels"
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/bots-framework/platforms/telegram"
	"reflect"
	"time"
)

type BiddingTicTacToeAppContext struct {
}

func (appCtx BiddingTicTacToeAppContext) AppUserEntityKind() string {
	return btttmodels.AppUserKind
}

func (appCtx BiddingTicTacToeAppContext) AppUserEntityType() reflect.Type {
	return reflect.TypeOf(&btttmodels.AppUserEntity{})
}

func (appCtx BiddingTicTacToeAppContext) NewBotAppUserEntity() bots.BotAppUser {
	return &btttmodels.AppUserEntity{
		DtCreated: time.Now(),
	}
}

func (appCtx BiddingTicTacToeAppContext) NewAppUserEntity() strongo.AppUser {
	return appCtx.NewBotAppUserEntity()
}

func (appCtx BiddingTicTacToeAppContext) GetTranslator(c context.Context) strongo.Translator {
	return strongo.NewMapTranslator(c, bttt_trans.TRANS)
}

type LocalesProvider struct {
}

func (LocalesProvider) GetLocaleByCode5(code5 string) (strongo.Locale, error) {
	return strongo.LocaleEnUS, nil
}

func (appCtx BiddingTicTacToeAppContext) SupportedLocales() strongo.LocalesProvider {
	return BtttLocalesProvider{}
}

type BtttLocalesProvider struct {
}

func (BtttLocalesProvider) GetLocaleByCode5(code5 string) (locale strongo.Locale, err error) {
	switch code5 {
	case strongo.LocaleCodeEnUS:
		return strongo.LocaleEnUS, nil
	case strongo.LocalCodeRuRu:
		return strongo.LocaleRuRu, nil
	default:
		return locale, errors.New("Unsupported locale: " + code5)
	}
}

var _ strongo.LocalesProvider = (*BtttLocalesProvider)(nil)

func (appCtx BiddingTicTacToeAppContext) GetBotChatEntityFactory(platform string) func() bots.BotChat {
	switch platform {
	case "telegram":
		return func() bots.BotChat {
			return &btttmodels.BtttTelegramChatEntity{
				TgChatEntityBase: *telegram.NewTelegramChatEntity(),
			}
		}
	default:
		panic("Unknown platform: " + platform)
	}
}

var _ bots.BotAppContext = (*BiddingTicTacToeAppContext)(nil)

var TheAppContext = BiddingTicTacToeAppContext{}
