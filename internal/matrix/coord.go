package matrix

import (
	"fmt"
)

type Coord struct {
	X int
	Y int
}

func (c Coord) Left() Coord {
	return Coord{c.X - 1, c.Y}
}

func (c Coord) Right() Coord {
	return Coord{c.X + 1, c.Y}
}

func (c Coord) Up() Coord {
	return Coord{c.X, c.Y - 1}
}

func (c Coord) Down() Coord {
	return Coord{c.X, c.Y + 1}
}

func (c Coord) String() string {
	return fmt.Sprintf("{X: %d, Y: %d}", c.X, c.Y)
}
