package game

type Direction int

const (
	Up        Direction = 0
	UpRight   Direction = 1
	Right     Direction = 2
	DownRight Direction = 3
	Down      Direction = 4
	DownLeft  Direction = 5
	Left      Direction = 6
	UpLeft    Direction = 7
)

func (direction Direction) GetOppositeDirection() Direction {
	switch direction {
	case Up:
		return Down
	case UpRight:
		return DownLeft
	case Right:
		return Left
	default :
		return UpLeft
	}
}



