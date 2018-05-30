package btttmodels

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/strongo/app"
	"github.com/strongo/db"
	"github.com/strongo/db/gaedb"
	"google.golang.org/appengine/datastore"
	"time"
	"github.com/pquerna/ffjson/ffjson"
	"strconv"
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

	// XToken string `datastore:",noindex,omitempty"`
	// OToken string `datastore:",noindex,omitempty"`

	// XBid int16 `datastore:",noindex,omitempty"`
	// OBid int16 `datastore:",noindex,omitempty"`

	// XBidTime time.Time `datastore:",noindex,omitempty"`
	// OBidTime time.Time `datastore:",noindex,omitempty"`

	// XPreviousBid int16 `datastore:",noindex,omitempty"`
	// OPreviousBid int16 `datastore:",noindex,omitempty"`
	// XBalance int `datastore:",noindex,omitempty"`
	// OBalance int `datastore:",noindex,omitempty"`

	// XTargetX int8 `datastore:",noindex,omitempty"`
	// XTargetY int8 `datastore:",noindex,omitempty"`
	// OTargetX int8 `datastore:",noindex,omitempty"`
	// OTargetY int8 `datastore:",noindex,omitempty"`

	XPlayer string `datastore:",noindex,omitempty"` // Holds GamePlayerJson
	OPlayer string `datastore:",noindex,omitempty"` // Holds GamePlayerJson

	// userX GameUserJson
	// userO GameUserJson

	XUserID int64 `datastore:",noindex,omitempty"`
	OUserID int64 `datastore:",noindex,omitempty"`

	// XUserName string `datastore:",noindex,omitempty"`
	// OUserName string `datastore:",noindex,omitempty"`

	// XTgChatID int64 `datastore:",noindex,omitempty"`
	// OTgChatID int64 `datastore:",noindex,omitempty"`

	// XTgMessageID int `datastore:",noindex,omitempty"`
	// OTgMessageID int `datastore:",noindex,omitempty"`

	Board   Board     `datastore:",noindex,omitempty"`
	Logbook GameTurns `datastore:",noindex,omitempty"`
}

func (game GameEntity) GetUsersXO() (userX, userO GamePlayerJson) {
	return game.PlayerX(), game.PlayerO()
}

func (game GameEntity) getUserJson(player Player, s string) (gameUserJson GamePlayerJson) {
	if s == "" {
		return
	}
	if err := ffjson.UnmarshalFast([]byte(s), &gameUserJson); err != nil {
		panic("Failed to unmarshal game.User" + string(player))
	}
	return
}

func (game GameEntity) PlayerX() (gameUserJson GamePlayerJson) {
	return game.getUserJson(PlayerX, game.XPlayer)
}

func (game GameEntity) PlayerO() (gameUserJson GamePlayerJson) {
	return game.getUserJson(PlayerO, game.OPlayer)
}

func (game GameEntity) GetPlayerJsonByUserID(userID int64) GamePlayerJson {
	switch userID {
	case game.XUserID:
		return game.PlayerX()
	case game.OUserID:
		return game.PlayerO()
	case 0:
		return GamePlayerJson{}
	default:
		panic(fmt.Sprintf("user does not belong to the game: userID=%v, XUserID=%v, OUserID=%v",
			userID, game.XUserID, game.OUserID))
	}
}

func (game GameEntity) GetPlayerJson(p Player) GamePlayerJson {
	switch p {
	case PlayerX:
		return game.PlayerX()
	case PlayerO:
		return game.PlayerO()
	default:
		panic(fmt.Sprintf("unknown player: %v", p))
	}
}

func (game *GameEntity) marshalUser(user GamePlayerJson) string {
	data, err := ffjson.MarshalFast(&user)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func (game *GameEntity) SetPlayerX(userX GamePlayerJson) (changed bool) {
	return game.SetPlayerJson(game.XUserID, userX)
}

func (game *GameEntity) SetPlayerO(userO GamePlayerJson) (changed bool) {
	return game.SetPlayerJson(game.OUserID, userO)
}

func (game *GameEntity) SetPlayerJson(userID int64, playerJson GamePlayerJson) (changed bool) {
	s := game.marshalUser(playerJson)
	if s == "{}" || s == "{Tg:{}}" {
		s = ""
	}
	switch userID {
	case game.XUserID:
		if s != game.XPlayer {
			game.XPlayer = s
			return true
		}
	case game.OUserID:
		if s != game.OPlayer {
			game.OPlayer = s
			return true
		}
	default:
		panic(fmt.Sprintf("not a player userID=%v", userID))
	}
	return false
}

func (game *GameEntity) SetPlayers(playerX, playerO GamePlayerJson) {
	game.SetPlayerJson(game.XUserID, playerX)
	game.SetPlayerJson(game.OUserID, playerO)
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
		return currentTurn.X.HasBidAndTarget() || (game.PlayerX().Balance == 0 && currentTurn.X.HasTarget())
	case PlayerO:
		return currentTurn.O.HasBidAndTarget() || (game.PlayerO().Balance == 0 && currentTurn.X.HasTarget())
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

func (game GameEntity) UserTelegramData(userID int64) (chatID int64, messageID int) {
	player := game.GetPlayerJsonByUserID(userID)
	// switch userID {
	// case game.XUserID:
	// 	player = game.PlayerX()
	// case game.OUserID:
	// 	return game.OTgChatID, game.OTgMessageID
	// default:
	// 	return 0, 0
	// }
	chatID, _  = strconv.ParseInt(player.Tg.ChatID, 10, 64)
	messageID, _ = strconv.Atoi(player.Tg.MessageID)
	return
}

func (game *GameEntity) Load(ps []datastore.Property) error {
	return datastore.LoadStruct(game, ps)
}

func (game *GameEntity) Save() (properties []datastore.Property, err error) {
	if properties, err = datastore.SaveStruct(game); err != nil {
		return
	}
	_ = game.Board.Size() // Check size

	verifyPlayerJson := func (userID int64, bid int) (GamePlayerJson, error) {
		player := game.GetPlayerJsonByUserID(userID)
		if player.Balance < 0 {
			return player, fmt.Errorf("balance < 0 for userID=%v", userID)
		}
		if userID != 0 && player.Name == "" {
			return player, fmt.Errorf("user name is empty for player with UserID=%v", userID)
		}
		if bid > player.Balance {
			return player, fmt.Errorf("bid > balance (%v > %v)", bid, player.Balance)
		}
		return player, nil
	}
	var userX, userO GamePlayerJson
	currentTurn := game.Logbook.CurrentTurn()
	if userX, err = verifyPlayerJson(game.XUserID, currentTurn.X.Bid); err != nil {
		return
	}
	if userO, err = verifyPlayerJson(game.OUserID, currentTurn.O.Bid); err != nil {
		return
	}
	if game.XUserID != 0 && game.OUserID != 0 {
		if (userX.Balance + userO.Balance) != 200 {
			return properties, errors.New(fmt.Sprintf("(game.XBalance + game.OBalance): %d != 200", userX.Balance+userO.Balance))
		}
	}
	if properties, err = gaedb.CleanProperties(properties, map[string]gaedb.IsOkToRemove{
		"DtLastTurn": gaedb.IsZeroTime,
	}); err != nil {
		return
	}
	return
}
