package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/nlm/adventofcode2023/internal/matrix"
	"github.com/nlm/adventofcode2023/internal/utils"
)

var (
	Left  = matrix.Vec{X: -1, Y: 0}
	Right = matrix.Vec{X: +1, Y: 0}
	Up    = matrix.Vec{X: 0, Y: -1}
	Down  = matrix.Vec{X: 0, Y: +1}
)

// var wg = sync.WaitGroup{}

func Run(m *matrix.Matrix[byte], coord matrix.Coord, direction matrix.Vec, visits *matrix.Matrix[Visit]) {
	// // stage.Println("RUN", coord, direction)
	// stage.Println(SVisits(m, visits, coord))

out:
	for {
		if m.InCoord(coord) {
			// visited
			visited := visits.AtCoord(coord)
			var dirVisited bool
			switch direction {
			case Up:
				dirVisited = visited.Up
				visited.Up = true
			case Down:
				dirVisited = visited.Down
				visited.Down = true
			case Left:
				dirVisited = visited.Left
				visited.Left = true
			case Right:
				dirVisited = visited.Right
				visited.Right = true
			}
			if dirVisited {
				// stage.Println(SVisits(m, visits))
				// stage.Println("-> visited", coord)
				break out
			}
			visits.SetAtCoord(coord, visited)
		}
		// fullVisits.SetAtCoord(coord, visited)

		// next
		nextCoord := coord.Add(direction)
		if !m.InCoord(nextCoord) {
			// stage.Println(SVisits(m, visits))
			// stage.Println("-> out of bounds")
			break out
		}
		next := m.AtCoord(nextCoord)

		switch next {
		case '\\':
			switch direction {
			case Up:
				// if visits.AtCoord(nextCoord).Right {
				// 	break out
				// }
				direction = Left
			case Down:
				// if visits.AtCoord(nextCoord).Left {
				// 	break out
				// }
				direction = Right
			case Left:
				// if visits.AtCoord(nextCoord).Down {
				// 	break out
				// }
				direction = Up
			case Right:
				// if visits.AtCoord(nextCoord).Up {
				// 	break out
				// }
				direction = Down
			}
		case '/':
			switch direction {
			case Up:
				// if visits.AtCoord(nextCoord).Left {
				// 	break out
				// }
				direction = Right
			case Down:
				// if visits.AtCoord(nextCoord).Right {
				// 	break out
				// }
				direction = Left
			case Left:
				direction = Down
			case Right:
				direction = Up
			}
		case '|':
			switch direction {
			case Left, Right:
				// wg.Add(2)
				// stage.Println("spawn up")
				Run(m, nextCoord, Up, visits)
				// stage.Println("spawn down")
				Run(m, nextCoord, Down, visits)
				// stage.Println("vsplit done")
				break out
			}
		case '-':
			switch direction {
			case Up, Down:
				// wg.Add(2)
				// stage.Println("spawn left")
				Run(m, nextCoord, Left, visits)
				// stage.Println("spawn right")
				Run(m, nextCoord, Right, visits)
				// stage.Println("hsplit done")
				break out
			}
		}
		// stage.Println(SVisits(m, visits, coord))
		coord = nextCoord
		// stage.Println(matrix.SMatrix(m))
		// stage.Println(matrix.SMatrix(visits))
		// if next != '.' {
		// stage.Println("\033[H\033[2J")
		// }
		// time.Sleep(1 * time.Millisecond)
	}
	// stage.Println(matrix.SMatrix(visits))
	// wg.Done()
}

type Visit struct {
	Up    bool
	Down  bool
	Left  bool
	Right bool
}

func (v Visit) Visited() bool {
	return v.Up || v.Down || v.Left || v.Right
}

func (v Visit) String() string {
	sum := 0
	utils.Filter[bool]([]bool{v.Up, v.Down, v.Left, v.Right}, func(b bool) bool {
		if b {
			sum++
		}
		return sum > 0
	})
	if sum > 1 {
		return "X"
	}
	if v.Up {
		return "^"
	}
	if v.Down {
		return "v"
	}
	if v.Left {
		return "<"
	}
	if v.Right {
		return ">"
	}
	return "?"
}

func SVisits(m *matrix.Matrix[byte], v *matrix.Matrix[Visit], coord matrix.Coord) string {
	sb := strings.Builder{}
	fmt.Fprintf(&sb, "   0123456789\n")
	for y := 0; y < v.Len.Y; y++ {
		fmt.Fprintf(&sb, "%2d ", y)
		for x := 0; x < v.Len.X; x++ {
			// if coord.X == x && coord.Y == y {
			// 	fmt.Fprint(&sb, "@")
			// 	continue
			// }
			vs := v.At(x, y)
			if m.At(x, y) != '.' {
				fmt.Fprint(&sb, string(m.At(x, y)))
				continue
			}
			if vs.Visited() {
				// fmt.Fprint(&sb, vs)
				fmt.Fprint(&sb, "*")
				continue
			}
			// fmt.Fprint(&sb, string(m.At(x, y)))
			fmt.Fprint(&sb, ".")
		}
		fmt.Fprintln(&sb)
	}
	return sb.String()
}

func CountEnergized(m *matrix.Matrix[Visit]) int {
	var count int
	for y := 0; y < m.Len.Y; y++ {
		for x := 0; x < m.Len.X; x++ {
			if m.At(x, y).Visited() {
				count++
			}
		}
	}
	return count
}

func Stage1(input io.Reader) (any, error) {
	m := utils.Must(matrix.NewFromReader(input))
	coord := matrix.Coord{X: -1, Y: 0}
	direction := Right
	visits := matrix.New[Visit](m.Len.X, m.Len.Y)
	// go func() {
	// 	t := time.NewTicker(100 * time.Millisecond)
	// 	for ; true; <-t.C {
	// 		fmt.Println(SVisits(m, visits, coord))
	// 	}
	// }()
	// fullVisits := visits.Clone()
	// wg.Add(1)
	Run(m, coord, direction, visits)
	// wg.Wait()
	return CountEnergized(visits), nil
}

func Stage2(input io.Reader) (any, error) {
	m := utils.Must(matrix.NewFromReader(input))

	var (
		coord     matrix.Coord
		direction matrix.Vec
		visits    *matrix.Matrix[Visit]
		maxEnergy int
	)

	// Up
	for x := 0; x < m.Len.X; x++ {
		coord = matrix.Coord{X: x, Y: -1}
		direction = Down
		visits = matrix.New[Visit](m.Len.X, m.Len.Y)
		Run(m, coord, direction, visits)
		energy := CountEnergized(visits)
		if energy > maxEnergy {
			maxEnergy = energy
		}
		// fmt.Println(SVisits(m, visits, coord))
	}

	// Down
	for x := 0; x < m.Len.X; x++ {
		coord = matrix.Coord{X: x, Y: m.Len.Y}
		direction = Up
		visits = matrix.New[Visit](m.Len.X, m.Len.Y)
		Run(m, coord, direction, visits)
		energy := CountEnergized(visits)
		if energy > maxEnergy {
			maxEnergy = energy
		}
		// fmt.Println(SVisits(m, visits, coord))
	}

	// Left
	for y := 0; y < m.Len.Y; y++ {
		coord = matrix.Coord{X: -1, Y: y}
		direction = Up
		visits = matrix.New[Visit](m.Len.X, m.Len.Y)
		Run(m, coord, direction, visits)
		energy := CountEnergized(visits)
		if energy > maxEnergy {
			maxEnergy = energy
		}
		// fmt.Println(SVisits(m, visits, coord))
	}

	//Right
	for y := 0; y < m.Len.Y; y++ {
		coord = matrix.Coord{X: m.Len.Y, Y: y}
		direction = Up
		visits = matrix.New[Visit](m.Len.X, m.Len.Y)
		Run(m, coord, direction, visits)
		energy := CountEnergized(visits)
		if energy > maxEnergy {
			maxEnergy = energy
		}
		// fmt.Println(SVisits(m, visits, coord))
	}

	return maxEnergy, nil
}
