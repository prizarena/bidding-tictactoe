package btttdal

import "github.com/strongo/db"

var (
	DB       db.Database
	User     UserDal
	Game     GameDal
	Referrer ReferrerDal
)
