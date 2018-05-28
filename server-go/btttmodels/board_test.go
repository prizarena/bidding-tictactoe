package btttmodels

import "testing"

func TestBoard_Draw(t *testing.T) {
	board, err := ParseBoard("__xooox__")
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v:\n%v", board, board.Draw())
	t.Logf("Is X a winner? %v", board.IsWinner(PlayerX))
	t.Logf("Is O a winner? %v", board.IsWinner(PlayerO))
}

func TestBoard_Cells(t *testing.T) {
	board := Board("__xooox__")
	cells := board.Cells()
	if len(cells) != 9 {
		t.Errorf("len(cells) != 9")
	}
}

func TestBoard_Grid(t *testing.T) {
	board := Board("__xooox__")
	grid := board.Grid()
	if len(grid) != 3 {
		t.Errorf("len(grid):%d != 3\n\tgrid: %v\n\tcells: %v", len(grid), grid, board.Cells())
	}
}
