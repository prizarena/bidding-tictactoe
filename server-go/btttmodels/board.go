package btttmodels

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

const (
	CellX     = 'X'
	CellO     = 'O'
	CellEmpty = '_'
)

type Player byte

const (
	PlayerX     Player = CellX
	PlayerO     Player = CellO
	IsTie       Player = '='
	NotPlayer   Player = 0x00
	NoWinnerYet Player = NotPlayer
)

func IsValidPlayer(p Player) bool {
	return p == PlayerX || p == PlayerO
}

type Coordinates struct {
	X int8
	Y int8
}

type Strike []Coordinates

func (strike Strike) CanWin(grid Grid) bool {
	canWin := func(p Player) bool {
		for _, cell := range strike {
			v := grid[cell.X][cell.Y]
			if v != p && v != CellEmpty {
				return false
			}
		}
		return true
	}
	return canWin(PlayerX) || canWin(PlayerO)
}

var strikes3x3 = [8]Strike{
	{{0, 0}, {1, 0}, {2, 0}}, // top row
	{{0, 1}, {1, 1}, {2, 1}}, // middle row
	{{0, 2}, {1, 2}, {2, 2}}, // bottom row
	{{0, 0}, {0, 1}, {0, 2}}, // left col
	{{1, 0}, {1, 1}, {1, 2}}, // middle col
	{{2, 0}, {2, 1}, {2, 2}}, // right row
	{{0, 0}, {1, 1}, {2, 2}}, // left top => right bottom
	{{2, 0}, {1, 1}, {0, 2}}, // right top => left bottom
}

const (
	MIN_BOARD_SIZE = 3
	MAX_BOARD_SIZE = 5
)

type Board string

func (board Board) Size() int {
	switch len(board) {
	case 3 * 3:
		return 3
	case 4 * 4:
		return 4
	case 5 * 5:
		return 5
	default:
		panic(fmt.Sprintf("Inalid board length: %d", len(board)))
	}
}

const (
	NoBoard Board = ""
)

func EmptyBoard(size int) Board {
	if size < MIN_BOARD_SIZE || size > MAX_BOARD_SIZE {
		panic(fmt.Sprintf("Invalid board size: %d", size))
	}
	return Board(strings.Repeat(string([]byte{CellEmpty}), size*size))
}

type Cell struct {
	X int8
	Y int8
	V Player
}

func ParseBoard(s string) (Board, error) {
	board := Board(s)
	_ = board.Size()
	for i, b := range []byte(s) {
		if !(b == byte(CellEmpty) || b == byte(CellX) || b == byte(CellO)) {
			return NoBoard, fmt.Errorf("Invalid board cell at position %d: '%v'", i, string([]byte{b}))
		}
	}
	return board, nil
}

func DrawPlayerToCell(p Player) string {
	switch p {
	case PlayerX:
		return "⚔️"
	case PlayerO:
		return "🍩"
	default:
		return "\u2B1C"
	}
}

func (board Board) Turn(cell Cell) (Board, error) {
	if cell.V != PlayerX && cell.V != PlayerO {
		return board, errors.New("cell.V != PlayerX || cell.V != PlayerO")
	}
	if cell.X <= 0 || cell.Y <= 0 {
		panic("Invalid cell")
	}
	b := []byte(board)
	x, y := cell.X-1, cell.Y-1
	cellIndex := y*3 + x
	if b[cellIndex] != CellEmpty {
		return board, errors.New(fmt.Sprintf("Cell already occupied: %v", string([]byte{b[cellIndex]})))
	}
	b[cellIndex] = byte(cell.V)
	return Board(b), nil
}

func (board Board) Draw() string {
	var buffer bytes.Buffer
	boardSizeMinus1 := int8(board.Size() - 1)
	for _, cell := range board.Cells() {
		buffer.WriteString(DrawPlayerToCell(cell.V))
		if cell.Y == boardSizeMinus1 {
			if cell.X < boardSizeMinus1 {
				buffer.WriteString("\n")
			}
		}
	}
	return buffer.String()
}

func (board Board) Cells() (cells []Cell) {
	var col, row int8
	boardLastIndex := int8(board.Size() - 1)
	cells = make([]Cell, len(board))
	for i, v := range []Player(board) {
		cells[i] = Cell{X: col, Y: row, V: v}
		if row == boardLastIndex {
			if col < boardLastIndex {
				row = 0
				col += 1
			}
		} else {
			row += 1
		}
	}
	return
}

func (board Board) CellValue(x, y int) (player Player) {
	return Player(board[(y-1)*board.Size()+(x-1)])
}

type Grid [][]Player

func (board Board) Grid() (grid Grid) {
	boardSize := board.Size()
	grid = make(Grid, 0, boardSize)
	var row []Player
	for _, cell := range board.Cells() {
		if cell.Y == 0 {
			row = make([]Player, boardSize)
			grid = append(grid, row)
		}
		row[cell.Y] = cell.V
	}
	return
}

func (board Board) IsEmpty() bool {
	return board == EmptyBoard(board.Size())
}

func (board Board) IsWinner(player Player) bool {
	grid := board.Grid()
	for _, strike := range strikes3x3 {
		isStrike := true
		for _, cell := range strike {
			if grid[cell.X][cell.Y] != player {
				isStrike = false
				break
			}
		}
		if isStrike {
			return true
		}
	}
	return false
}

func (board Board) Winner() Player {
	grid := board.Grid()
	var canWin bool
	for _, strike := range strikes3x3 {
		canWin = canWin || strike.CanWin(grid)
	}
	if !canWin {
		return IsTie
	}
	if board.IsWinner(PlayerX) {
		return PlayerX
	}
	if board.IsWinner(PlayerO) {
		return PlayerO
	}
	return NoWinnerYet
}
