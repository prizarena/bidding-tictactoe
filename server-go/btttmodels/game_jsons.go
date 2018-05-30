package btttmodels

//go:generate ffjson $GOFILE

import (
	"strconv"
	"time"
)

type gamePlayerTgJson struct {
	ChatID    string `json:",omitempty"`
	MessageID string `json:",omitempty"`
}

type GamePlayerJson struct {
	Name    string    `json:",omitempty"`
	Balance int       `json:",omitempty"`
	BidTime time.Time `json:",omitempty"`
	Tg gamePlayerTgJson // Use private struct as inline struct are not supported well by ffjson
}

func (p GamePlayerJson) TgChatID() int64 {
	if p.Tg.ChatID == "" {
		return 0
	}
	if v, err := strconv.ParseInt(p.Tg.ChatID, 10, 64); err != nil {
		panic(err)
	} else {
		return v
	}
}

func (p GamePlayerJson) TgMessageID() int {
	if p.Tg.MessageID == "" {
		return 0
	}
	if v, err := strconv.Atoi(p.Tg.MessageID); err != nil {
		panic(err)
	} else {
		return v
	}
}
