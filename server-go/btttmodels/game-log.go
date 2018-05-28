package btttmodels

import (
	"encoding/binary"
	"github.com/strongo/decimal"
)

const turnRecordLength = 5

type GameLog []byte

func (gameLog GameLog) TurnsCount() int {
	l := len(gameLog)
	if l%2 != 0 {
		//noinspection ALL
		panic("len(gameLog) % 2 != 0")
	}
	return l / turnRecordLength
}

func (gameLog GameLog) TurnWinner(i int) Move {
	i = i * turnRecordLength
	turn := Turn(gameLog[i : i+5])
	p1 := turn.PlayerMove(1)
	p2 := turn.PlayerMove(2)
	if p1.bid > p2.bid {
		return p1
	} else if p2.Bid() > p1.Bid() {
		return p2
	}
	return p1
}

// 5 bytes:
// 0 - 1st Player bid
// 1 - 1st Player bid
// 2 - 2nd Player bid
// 3 - 2nd Player bid
// 4 & 0f - 1st Player target
// 4 & f0 - 2nd Player target

type Turn []byte

func (turn Turn) PlayerMove(player uint8) Move {
	i := (player - 1) * 2
	return NewMove(player, turn[4], []byte(turn)[i:i+2])
}

type Move struct {
	player uint8
	bid    uint16
	target byte // TODO: Document what it is.
}

func NewMove(player uint8, target byte, bid []byte) Move {
	move := Move{
		player: player,
		bid:    binary.BigEndian.Uint16(bid),
		target: target,
	}
	return move
}

func (move Move) Target() uint8 {
	switch move.player {
	case 1:
		return move.target & 0x0F
	case 2:
		return move.target & 0xF0
	default:
		panic("Player out of range")
	}
}

func (move Move) Bid() decimal.Decimal64p2 {
	return decimal.Decimal64p2(move.bid)
}
