package bttt_app_gae_standard

import (
	"testing"
	"github.com/prizarena/bidding-tictactoe/server-go/btttdal"
	"github.com/strongo/log"
)

func TestMain1(t *testing.T) {
	if btttdal.DB == nil {
		t.Error("btttdal.DB == nil")
	}
	if log.NumberOfLoggers() == 0 {
		t.Error("NumberOfLoggers() == 0")
	}
}
