package btttmodels

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"
)

type TurnRequest struct {
	GameID int64
	X      int
	Y      int
}

func (tr TurnRequest) Valid() bool {
	return tr.GameID != 0 && tr.X > 0 && tr.Y > 0
}

func (tr TurnRequest) ToCell(player Player) Cell {
	return Cell{int8(tr.X), int8(tr.X), player}
}

func (tr TurnRequest) String() string {
	return tr.ToUrlQuery()
}

const (
	GAME_ID_ENCODING_BASE = 10
)

func (tr TurnRequest) ToUrlQuery() string {
	var buffer bytes.Buffer
	buffer.WriteString("g=" + strconv.FormatInt(tr.GameID, GAME_ID_ENCODING_BASE))
	buffer.WriteString(fmt.Sprintf("&c=%d%d", tr.X, tr.Y))
	return buffer.String()
}

func ParseQueryToTurnRequest(query string) (tr TurnRequest, err error) {
	values, err := url.ParseQuery(query)
	if err != nil {
		return tr, err
	}

	if v := values.Get("g"); v != "" {
		if tr.GameID, err = strconv.ParseInt(v, GAME_ID_ENCODING_BASE, 64); err != nil {
			return tr, err
		}
	}

	if c := values.Get("c"); c != "" {
		if tr.X, err = strconv.Atoi(c[:1]); err != nil {
			return tr, err
		}
		if tr.Y, err = strconv.Atoi(c[1:]); err != nil {
			return tr, err
		}
	}
	return tr, err
}
