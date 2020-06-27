package game

import (
	"errors"
	"math/rand"
)

const Rows = 6
const Cols = 7

type State struct {
	Grid           [Cols][Rows]int `json:"grid"`
	Outcome        int             `json:"int"`
	CurrentPlayer  int             `json:"currentplayer"`
	PreviousPlayer int             `json:"previousplayer"`
	Move           Coordinate      `json:"move"`
	NbMoves        int             `json:"nbmoves"`
}

func NewState() *State {
	return &State{CurrentPlayer: Player1}
}

//func (state *State) UnmarshalJSON(data []byte) error {
//	var results map[string]interface{}
//	if err := json.Unmarshal(data, &results); err != nil {
//		return err
//	}
//	for k, v := range results {
//		fmt.Println(k, reflect.TypeOf(v))
//
//	}
//
//	return nil
//}

func (state State) GetFreeRows() [7]int {
	var freeRows [Cols]int
	for col := 0; col < Cols; col++ {
		freeRows[col] = state.GetFreeRow(col)
	}
	return freeRows
}
func (state State) GetFreeRow(col int) int {
	var row = 0
	for row < Rows {
		if state.Grid[col][row] != 0 {
			break
		}
		row++
	}
	return row - 1
}
func (state *State) PlayMove(move Coordinate) {
	state.setCellState(move, state.CurrentPlayer)
	state.Move = move
	state.NbMoves += 1

	if state.NbMoves == Cols*Rows {
		state.Outcome = Draw
	}
	if state.HasConnectedFour(move) {
		state.Outcome = state.CurrentPlayer
	}

	state.changeTurn()

}
func (state *State) setCellState(coordinate Coordinate, player int) {
	state.Grid[coordinate.Col][coordinate.Row] = player
}
func (state State) HasConnectedFour(move Coordinate) bool {
	for direction := Up; direction <= DownRight; direction++ {
		var oppositeDirection = direction.GetOppositeDirection()
		var count = 1
		count += state.NbConnectedInDirection(move, direction)
		count += state.NbConnectedInDirection(move, oppositeDirection)
		if count >= 4 {
			return true
		}
	}

	return false
}
func (state State) NbConnectedInDirection(move Coordinate, direction Direction) int {
	var coordinate = move.GetCoordinateInDirection(direction)
	count := 0
	for coordinate.IsInBound(coordinate) {
		if state.GetPlayerFromCoordinate(coordinate) != state.CurrentPlayer {
			break
		}
		count++
		coordinate = coordinate.GetCoordinateInDirection(direction)
	}
	return count
}
func (state *State) changeTurn() {
	state.PreviousPlayer = state.CurrentPlayer

	switch state.CurrentPlayer {
	case Player1:
		state.CurrentPlayer = Player2
	case Player2:
		state.CurrentPlayer = Player1
	}
}
func (state State) Playout() int {
	//fmt.Println("Playout")
	for state.Outcome == None {
		freeRows := state.GetFreeRows()

		move, err := state.GetRandomFreeRow(freeRows)
		if err != nil {
			break
		}
		//fmt.Println(move)
		state.PlayMove(move)
	}
	//fmt.Println(state.int)
	//fmt.Println(state.Grid)
	return state.Outcome
}
func (state State) GetPlayerFromCoordinate(coordinate Coordinate) int {
	return state.Grid[coordinate.Col][coordinate.Row]
}
func (state *State) GetRandomFreeRow(freeRows [7]int) (Coordinate, error) {
	var col = rand.Intn(len(freeRows) - 1)
	var counter = 0
	for counter < Cols {
		var freeRow = freeRows[col]
		if freeRow != -1 {
			return Coordinate{
				Col: col,
				Row: freeRow,
			}, nil
		}
		col = (col + 1) % Cols
		counter++
	}
	return Coordinate{}, errors.New("no free row available")
}

func (state *State) GetLeftStateID() int {
	var total = 0
	var rowMultiplier = 1
	for row := Rows - 1; row >= 0; row-- {
		var rowTotal = 0
		var columnMultiplier = 1
		for col := 0; col < Cols; col++ {
			rowTotal += state.Grid[col][row] * columnMultiplier
			columnMultiplier *= 3
		}
		total += rowTotal * rowMultiplier
		rowMultiplier *= 158
	}
	return total
}
func (state *State) GetRightStateID() int {
	var total = 0
	var rowMultiplier = 1
	for row := Rows - 1; row >= 0; row-- {
		var rowTotal = 0
		var columnMultiplier = 1
		for col := 0; col < Cols; col++ {
			rowTotal += state.Grid[Cols-1-col][row] * columnMultiplier
			columnMultiplier *= 3
		}
		total += rowTotal * rowMultiplier
		rowMultiplier *= 158
	}
	return total
}

//func PrintMemUsage() {
//	var m runtime.MemStats
//	runtime.ReadMemStats(&m)
//	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
//	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
//	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
//	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
//	fmt.Printf("\tNumGC = %v\n", m.NumGC)
//}

//func bToMb(b uint64) uint64 {
//	return b / 1024 / 1024
//}
