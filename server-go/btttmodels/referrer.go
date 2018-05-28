package btttmodels

import (
	"github.com/strongo/db"
	"time"
)

const ReferrerKind = "R"

type Referrer struct {
	db.StringID
	*ReferrerEntity
}

var _ db.EntityHolder = (*Referrer)(nil)

func (Referrer) Kind() string {
	return ReferrerKind
}

func (referrer Referrer) Entity() interface{} {
	return referrer.ReferrerEntity
}

func (referrer *Referrer) SetEntity(v interface{}) {
	if v == nil {
		referrer.ReferrerEntity = nil
	} else {
		referrer.ReferrerEntity = v.(*ReferrerEntity)
	}
}

func (Referrer) NewEntity() interface{} {
	return new(ReferrerEntity)
}

type ReferrerEntity struct {
	Created       time.Time
	CreatorUserID int64
}
