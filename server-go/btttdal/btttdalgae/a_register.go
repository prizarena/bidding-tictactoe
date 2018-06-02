package btttdalgae

import (
	"github.com/prizarena/bidding-tictactoe/server-go/btttdal"
	"github.com/strongo/db/gaedb"
)

func RegisterGaeDal() {
	btttdal.DB = gaedb.NewDatabase()
	btttdal.User = userDalGae{}
	btttdal.Game = gameDal{}
	btttdal.Referrer = referrerDalGae{}
}
