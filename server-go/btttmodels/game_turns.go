package btttmodels

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type GameMove struct {
	Bid     int `json:"b,omitempty"`
	TargetX int `json:"x,omitempty"`
	TargetY int `json:"y,omitempty"`
}

func (move GameMove) HasTarget() bool {
	return move.TargetX > 0 && move.TargetY > 0
}

func (move GameMove) HasBid() bool {
	return move.Bid > 0
}

func (move GameMove) HasBidAndTarget() bool {
	return move.Bid > 0 && move.TargetX > 0 && move.TargetY > 0
}

type GameTurn struct {
	X      GameMove `json:"X,omitempty"`
	O      GameMove `json:"O,omitempty"`
	Winner Player   `json:"w,omitempty"`
}

func (t GameTurn) HasBothBidsAndTargets() bool {
	return t.X.HasBidAndTarget() && t.O.HasBidAndTarget()
}

func (t GameTurn) writeTurn(buf *bytes.Buffer) {
	// Example (left=X, right=O): 14@11>34@22 - < or > points to the winner
	writeMove := func(m GameMove) {
		if m.Bid != 0 || m.TargetY != 0 || m.TargetX != 0 {
			buf.WriteString(strconv.Itoa(m.Bid) + "@" + strconv.Itoa(m.TargetX) + strconv.Itoa(m.TargetY))
		}
	}
	writeMove(t.X)
	switch t.Winner {
	case PlayerX:
		buf.WriteString("<")
	case PlayerO:
		buf.WriteString(">")
	default:
		buf.WriteString("|")
	}
	writeMove(t.O)
}

const turnsSeparator = "\n"

type GameTurns string

func (gt GameTurns) Turns() (turns []GameTurn) {
	turnStrings := strings.Split(string(gt), turnsSeparator)
	turns = make([]GameTurn, len(turnStrings))
	var err error
	for i, turnStr := range turnStrings {
		if strings.HasSuffix(turnStr, "\r") { // Workaround for manual edit
			turnStr = turnStr[:len(turnStr)-1]
		}
		if turnStr == "" && i == len(turns)-1 {
			return
		} else if turns[i], err = parseTurn(turnStr); err != nil {
			panic(errors.WithMessage(err, fmt.Sprintf("invalid turn #%v: %v => %v", i+1, turnStr, []byte(turnStr))))
			return
		}
	}
	return
}

func parseTurn(s string) (turn GameTurn, err error) {
	var movesStrings []string
	if s == "|" {
		return
	} else if movesStrings = strings.Split(s, "<"); len(movesStrings) > 1 {
		turn.Winner = PlayerX
	} else if movesStrings = strings.Split(s, ">"); len(movesStrings) > 1 {
		turn.Winner = PlayerO
	} else if movesStrings = strings.Split(s, "|"); len(movesStrings) > 1 {
		// no winner yet, do nothing
	} else if s == "" {
		return
	} else {
		panic("turn record is missing players separator: " + s)
	}
	if len(movesStrings) != 2 {
		panic("unexpected number of moves: " + s)
	}
	if turn.X, err = parseMove(movesStrings[0]); err != nil {
		err = errors.WithMessage(err, "invalid move for X")
	}
	if turn.O, err = parseMove(movesStrings[1]); err != nil {
		err = errors.WithMessage(err, "invalid move for O")
	}
	return
}

func (gt GameTurns) CurrentTurn() (turn GameTurn) {
	s := string(gt)
	lastIndex := strings.LastIndex(s, turnsSeparator)
	if lastIndex > 0 {
		s = s[lastIndex+1:]
	}
	var err error
	if turn, err = parseTurn(s); err != nil {
		panic(err)
	}
	return
}

func (gt GameTurns) PreviousTurn() (turn GameTurn) {
	s := string(gt)
	lastIndex := strings.LastIndex(s, turnsSeparator)
	if lastIndex < 0 {
		return
	}
	return GameTurns(s[:lastIndex]).CurrentTurn()
}

func parseMove(s string) (move GameMove, err error) {
	if s == "" {
		return
	}
	vals := strings.Split(s, "@")
	if len(vals) != 2 {
		err = errors.New("invalid move string: " + s)
		return
	}
	if move.Bid, err = strconv.Atoi(vals[0]); err != nil {
		err = errors.New("invalid bid format: " + vals[0])
		return
	}
	switch len(vals[1]) {
	case 0:
		// Do nothing
	case 2:
		if move.TargetX, err = strconv.Atoi(vals[1][0:1]); err != nil {
			err = errors.New("invalid x format: " + vals[1][0:1])
			return
		}
		if move.TargetY, err = strconv.Atoi(vals[1][1:2]); err != nil {
			err = errors.New("invalid y format: " + vals[1][1:2])
			return
		}
	default:
		err = fmt.Errorf("invalid target (len=%v): [%v]", len(vals[1]), []byte(vals[1]))
		return
	}
	return
}

func (gt GameTurns) LogBid(player Player, bid int) GameTurns {
	turns := gt.Turns()
	lastTurnIndex := len(turns) - 1
	currentTurn := turns[lastTurnIndex]
	switch player {
	case PlayerX:
		currentTurn.X.Bid = bid
	case PlayerO:
		currentTurn.O.Bid = bid
	default:
		panic("unknown player: " + string([]byte{byte(player)}))
	}
	turns[lastTurnIndex] = currentTurn
	return turnsToGameTurns(turns)
}

func (gt GameTurns) LogTarget(player Player, x, y int) GameTurns {
	if x == 0 {
		panic("x == 0")
	}
	if y == 0 {
		panic("y == 0")
	}
	turns := gt.Turns()
	lastTurnIndex := len(turns) - 1
	currentTurn := turns[lastTurnIndex]
	switch player {
	case PlayerX:
		currentTurn.X.TargetX = x
		currentTurn.X.TargetY = y
	case PlayerO:
		currentTurn.O.TargetX = x
		currentTurn.O.TargetY = y
	default:
		panic("unknown player: " + string([]byte{byte(player)}))
	}
	turns[lastTurnIndex] = currentTurn
	return turnsToGameTurns(turns)
}

func (gt GameTurns) SetTurnWinner(player Player, startNewTurn bool) GameTurns {
	if player != PlayerX && player != PlayerO {
		panic("unknown player: " + string([]byte{byte(player)}))
	}
	turns := gt.Turns()
	lastTurnIndex := len(turns) - 1
	currentTurn := turns[lastTurnIndex]
	turns[lastTurnIndex] = currentTurn
	if startNewTurn {
		turns = append(turns, GameTurn{})
	}
	return turnsToGameTurns(turns)
}

func turnsToGameTurns(turns []GameTurn) GameTurns {
	buf := new(bytes.Buffer)
	for i := range turns {
		turns[i].writeTurn(buf)
		buf.WriteString(turnsSeparator)
	}
	b := buf.Bytes()
	return GameTurns(b[:len(b)-1])
}
