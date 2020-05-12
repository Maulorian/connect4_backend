package game

type Coordinate struct {
	Col int `json:"col"`
	Row int `json:"row"`
}

func (coordinate Coordinate) IsInBound(move Coordinate) bool {
	if move.Col < 0 || move.Col > Cols-1 {
		return false
	}
	if move.Row < 0 || move.Row > Rows-1 {
		return false
	}
	return true
}

func (coordinate *Coordinate) GetCoordinateInDirection(direction Direction) Coordinate {
	switch direction {
	case Up:
		return Coordinate{Col: coordinate.Col, Row: coordinate.Row + 1}
	case UpRight:
		return Coordinate{Col: coordinate.Col + 1, Row: coordinate.Row + 1}
	case Right:
		return Coordinate{Col: coordinate.Col + 1, Row: coordinate.Row}
	case DownRight:
		return Coordinate{Col: coordinate.Col + 1, Row: coordinate.Row - 1}
	case Down:
		return Coordinate{Col: coordinate.Col, Row: coordinate.Row - 1}
	case DownLeft:
		return Coordinate{Col: coordinate.Col - 1, Row: coordinate.Row - 1}
	case Left:
		return Coordinate{Col: coordinate.Col - 1, Row: coordinate.Row}
	default:
		return Coordinate{Col: coordinate.Col - 1, Row: coordinate.Row + 1}
	}
}
