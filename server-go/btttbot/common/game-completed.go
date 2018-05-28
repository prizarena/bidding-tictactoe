package common

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/strongo/app"
	"github.com/strongo/bidding-tictactoe-bot/bttt-trans"
	"github.com/strongo/bidding-tictactoe-bot/btttdal"
	"github.com/strongo/bidding-tictactoe-bot/btttmodels"
	"html"
)

func writeTgFooterForWinner(c context.Context, buf *bytes.Buffer, currentUser btttmodels.AppUser, translator strongo.SingleLocaleTranslator, game btttmodels.Game) (err error) {
	var winnerUser btttmodels.AppUser
	if currentUser.ID != 0 {
		currentUserPlayer := game.Player(currentUser.ID)
		if game.Board.IsWinner(currentUserPlayer) {
			winnerUser = currentUser
		}
	}
	if winnerUser.ID == 0 {
		if game.Board.IsWinner(btttmodels.PlayerX) {
			winnerUser.ID = game.XUserID
		} else if game.Board.IsWinner(btttmodels.PlayerO) {
			winnerUser.ID = game.OUserID
		} else {
			err = errors.New("Unknown winner")
			return
		}
	}

	if winnerUser.AppUserEntity == nil {
		if winnerUser, err = btttdal.User.GetUserByID(c, winnerUser.ID); err != nil {
			return
		}
	}

	buf.WriteString(translator.Translate(bttt_trans.MT_USER_WON_GAME, html.EscapeString(winnerUser.FullName())))
	fmt.Fprintf(buf, "\n<pre>%v</pre>", game.Board.Draw())

	//footer += "\n" + translator.Translate(bttt_trans.MT_TOURNAMENT_201710_SHORT) +
	//	"\n" + translator.Translate(bttt_trans.MT_TOURNAMENT_201710_SPONSOR) +
	//	"\n" + translator.Translate(bttt_trans.MT_TOURNAMENT_201710_LEARN_MORE)

	buf.WriteString("\n\n" + translator.Translate(bttt_trans.MT_ASK_TO_RATE))

	buf.WriteString("\n\n" + translator.Translate(bttt_trans.OUR_TWITTER, `<a href="https://twitter.com/TicTacToeBid">@TicTacToeBid</a>`))
	buf.WriteString("\n" + translator.Translate(bttt_trans.OUR_FB_PAGE, `<a href="https://fm.me/BiddingTicTacToe">@BiddingTicTacToe</a>`))
	buf.WriteString("\n" + translator.Translate(bttt_trans.OUR_WEBSITE, `<a href="https://biddingtictactoe.com/">BiddingTicTacToe.com</a>`))

	return
}
