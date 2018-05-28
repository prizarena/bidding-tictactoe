package btttmodels

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/strongo/app"
	"github.com/strongo/db"
	"github.com/strongo/db/gaedb"
	"google.golang.org/appengine/datastore"
	"time"
)

const GameKind = "G"

type Game struct {
	db.IntegerID
	*GameEntity
}

var _ db.EntityHolder = (*Game)(nil)

func (Game) Kind() string {
	return GameKind
}

func (game Game) Entity() interface{} {
	return game.GameEntity
}

func (Game) NewEntity() interface{} {
	return new(GameEntity)
}

func (game *Game) SetEntity(v interface{}) {
	if v == nil {
		game.GameEntity = nil
	} else {
		game.GameEntity = v.(*GameEntity)
	}
}

func (game GameEntity) GetLocale() strongo.Locale {
	switch game.Locale {
	case "ru-RU":
		return strongo.LocaleRuRu
	}
	return strongo.LocaleEnUS
}

type GameEntity struct {
	DtLastTurn        time.Time
	Status            string `datastore:",noindex,omitempty"`
	Locale            string `datastore:",noindex,omitempty"`
	TgChatInstance    string `datastore:",omitempty"`
	TgInlineMessageID string `datastore:",noindex,omitempty"`
	PrevGameID        int64  `datastore:",noindex,omitempty"`

	CountOfMoves int `datastore:",noindex,omitempty"`
	CountOfTurns int `datastore:",noindex,omitempty"`

	XToken string `datastore:",noindex,omitempty"`
	OToken string `datastore:",noindex,omitempty"`

	//XBid int16 `datastore:",noindex,omitempty"`
	//OBid int16 `datastore:",noindex,omitempty"`

	XBidTime time.Time `datastore:",noindex,omitempty"`
	OBidTime time.Time `datastore:",noindex,omitempty"`

	//XPreviousBid int16 `datastore:",noindex,omitempty"`
	//OPreviousBid int16 `datastore:",noindex,omitempty"`
	XBalance int `datastore:",noindex,omitempty"`
	OBalance int `datastore:",noindex,omitempty"`

	//XTargetX int8 `datastore:",noindex,omitempty"`
	//XTargetY int8 `datastore:",noindex,omitempty"`
	//OTargetX int8 `datastore:",noindex,omitempty"`
	//OTargetY int8 `datastore:",noindex,omitempty"`

	XUserID int64 `datastore:",noindex,omitempty"`
	OUserID int64 `datastore:",noindex,omitempty"`

	XUserName string `datastore:",noindex,omitempty"`
	OUserName string `datastore:",noindex,omitempty"`

	XTgChatID int64 `datastore:",noindex,omitempty"`
	OTgChatID int64 `datastore:",noindex,omitempty"`

	XTgMessageID int `datastore:",noindex,omitempty"`
	OTgMessageID int `datastore:",noindex,omitempty"`

	Board   Board     `datastore:",noindex,omitempty"`
	Logbook GameTurns `datastore:",noindex,omitempty"`
}

func (game GameEntity) Rival(player Player) Player {
	switch player {
	case PlayerX:
		return PlayerO
	case PlayerO:
		return PlayerX
	default:
		panic(fmt.Sprintf("Unknown player: %v", player))
	}
}

func (game GameEntity) HasBid(player Player) bool {
	currentTurn := game.Logbook.CurrentTurn()
	switch player {
	case PlayerX:
		return currentTurn.X.HasBid()
	case PlayerO:
		return currentTurn.O.HasBid()
	default:
		panic(fmt.Sprintf("Unknown player: %v", player))
	}
}

func (game GameEntity) HasTarget(player Player) bool {
	currentTurn := game.Logbook.CurrentTurn()
	switch player {
	case PlayerX:
		return currentTurn.X.HasTarget()
	case PlayerO:
		return currentTurn.O.HasTarget()
	default:
		panic(fmt.Sprintf("Unknown player: %v", player))
	}
}

func (game GameEntity) HasBidAndTarget(player Player) bool {
	currentTurn := game.Logbook.CurrentTurn()
	switch player {
	case PlayerX:
		return currentTurn.X.HasBidAndTarget() || (game.XBalance == 0 && currentTurn.X.HasTarget())
	case PlayerO:
		return currentTurn.O.HasBidAndTarget() || (game.OBalance == 0 && currentTurn.X.HasTarget())
	default:
		panic(fmt.Sprintf("Unknown player: %v", player))
	}
}

func (game GameEntity) TargetCell(player Player) Cell {
	currentTurn := game.Logbook.CurrentTurn()
	switch player {
	case PlayerX:
		return Cell{X: int8(currentTurn.X.TargetX), Y: int8(currentTurn.X.TargetY), V: player}
	case PlayerO:
		return Cell{X: int8(currentTurn.O.TargetX), Y: int8(currentTurn.O.TargetY), V: player}
	default:
		panic(fmt.Sprintf("Unknown player: %v", player))
	}
}

func (game GameEntity) Player(userID int64) Player {
	switch userID {
	case game.XUserID:
		return PlayerX
	case game.OUserID:
		return PlayerO
	default:
		return NotPlayer
	}
}

func (game GameEntity) UserBalance(userID int64) int {
	switch userID {
	case game.XUserID:
		return game.XBalance
	case game.OUserID:
		return game.OBalance
	default:
		panic(fmt.Sprintf("user does not belong to the game: %v", userID))
	}
}

func (game GameEntity) RivalUserUD(userID int64) int64 {
	switch userID {
	case game.XUserID:
		return game.OUserID
	case game.OUserID:
		return game.XUserID
	default:
		return 0
	}
}

func (game GameEntity) HasFreeSlots() bool {
	return game.OUserID == 0 || game.XUserID == 0
}

func (game GameEntity) UserTelegramData(userID int64) (charID int64, messageID int) {
	switch userID {
	case game.XUserID:
		return game.XTgChatID, game.XTgMessageID
	case game.OUserID:
		return game.OTgChatID, game.OTgMessageID
	default:
		return 0, 0
	}
}

func (game *GameEntity) Load(ps []datastore.Property) error {
	return datastore.LoadStruct(game, ps)
}

func (game *GameEntity) Save() (properties []datastore.Property, err error) {
	if properties, err = datastore.SaveStruct(game); err != nil {
		return
	}
	_ = game.Board.Size() // Check size
	if game.XBalance < 0 {
		return properties, errors.New("XBalance < 0")
	}
	if game.OBalance < 0 {
		return properties, errors.New("OBalance < 0")
	}
	currentTurn := game.Logbook.CurrentTurn()
	if currentTurn.X.Bid > game.XBalance {
		return properties, errors.New("game.XBid > game.XBalance")
	}
	if currentTurn.O.Bid > game.OBalance {
		return properties, errors.New("game.OBid > game.OBalance")
	}
	if (game.XBalance + game.OBalance) != 200 {
		return properties, errors.New(fmt.Sprintf("(game.XBalance + game.OBalance): %d != 100", game.XBalance+game.OBalance))
	}
	if game.XUserID != 0 && game.XUserName == "" {
		return properties, errors.New("game.XUserID != 0 && game.XUserName is empty string")
	}
	if game.OUserID != 0 && game.OUserName == "" {
		return properties, errors.New("game.OUserID != 0 && game.OUserName is empty string")
	}

	if properties, err = gaedb.CleanProperties(properties, map[string]gaedb.IsOkToRemove{
		"DtLastTurn": gaedb.IsZeroTime,
	}); err != nil {
		return
	}
	return
}
