package btttmodels

import (
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/db"
	"github.com/strongo/db/gaedb"
	"google.golang.org/appengine/datastore"
	"time"
	"github.com/strongo/app"
)

const AppUserKind = "User"

type AppUser struct {
	db.IntegerID
	*AppUserEntity
}

func NewAppUser(id int64) AppUser {
	return AppUser{IntegerID: db.IntegerID{ID: id}}
}

func (AppUser) Kind() string {
	return AppUserKind
}

func (user AppUser) Entity() interface{} {
	return user.AppUserEntity
}

func (user *AppUser) SetEntity(v interface{}) {
	if v == nil {
		user.AppUserEntity = nil
	} else {
		user.AppUserEntity = v.(*AppUserEntity)
	}
}

func (AppUser) NewEntity() interface{} {
	return new(AppUserEntity)
}

type AppUserEntity struct {
	DtCreated  time.Time
	ReferrerID string
	strongo.AppUserBase
}

var _ bots.BotAppUser = (*AppUserEntity)(nil)

func (entity *AppUserEntity) SetBotUserID(platform, botID, botUserId string) {

}

func (entity *AppUserEntity) SetPreferredLocale(code5 string) error {
	entity.Locale = code5
	return nil
}

func (entity *AppUserEntity) GetPreferredLocale() string {
	return entity.Locale
}

func (entity *AppUserEntity) SetNames(first, last, user string) {
	entity.FirstName = first
	entity.LastName = last
	entity.UserName = user
}

func (entity *AppUserEntity) GetCurrencies() (result []string) { // TODO: Temporary to satisfy obsolete member of interface
	return
}

func (entity *AppUserEntity) Load(ps []datastore.Property) error {
	return datastore.LoadStruct(entity, ps)
}

func (entity *AppUserEntity) Save() (properties []datastore.Property, err error) {
	if properties, err = datastore.SaveStruct(entity); err != nil {
		return
	}

	if properties, err = gaedb.CleanProperties(properties, map[string]gaedb.IsOkToRemove{
		"ReferrerID": gaedb.IsEmptyString,
		"FirstName":  gaedb.IsEmptyString,
		"LastName":   gaedb.IsEmptyString,
		"UserName":   gaedb.IsEmptyString,
		"Locale":     gaedb.IsEmptyString,
	}); err != nil {
		return
	}
	return
}
