package btttmodels

import (
	"log"
	"testing"
)

func TestGameTurns_LogBid(t *testing.T) {
	var gameTurns GameTurns

	gameTurns = gameTurns.LogBid(PlayerX, 10)
	if string(gameTurns) != "10@00|" {
		log.Fatalf("Unexpected gameTurns: %v", gameTurns)
	}
	gameTurns = gameTurns.LogTarget(PlayerX, 2, 3)
	if string(gameTurns) != "10@23|" {
		log.Fatalf("Unexpected gameTurns: %v", gameTurns)
	}
}

func TestGameTurns_LogTarget(t *testing.T) {
	var gameTurns GameTurns

	gameTurns = gameTurns.LogTarget(PlayerX, 2, 3)
	if string(gameTurns) != "0@23|" {
		log.Fatalf("Unexpected gameTurns: %v", gameTurns)
	}
	gameTurns = gameTurns.LogBid(PlayerX, 10)
	if string(gameTurns) != "10@23|" {
		log.Fatalf("Unexpected gameTurns: %v", gameTurns)
	}
}

func TestGameTurns_Turns(t *testing.T) {
	turns := GameTurns("10@23|12@23").Turns()
	if len(turns) != 1 {
		t.Fatalf("len(turns) != 1: %v", len(turns))
	}
	if turns[0].X.Bid != 10 {
		t.Fatalf("turns[0].X.Bid != 10: %v", turns[0].X.Bid)
	}
	if turns[0].X.TargetX != 2 {
		t.Fatalf("turns[0].X.TargetX != 2: %v", turns[0].X.TargetX)
	}
	if turns[0].X.TargetY != 3 {
		t.Fatalf("turns[0].X.TargetY != 3: %v", turns[0].X.TargetY)
	}
}

func TestParseMove(t *testing.T) {
	move, err := parseMove("10@23")
	if err != nil {
		log.Fatal(err)
	}
	t.Log(move)
}
