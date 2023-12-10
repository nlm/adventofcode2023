package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"strings"

	"github.com/nlm/adventofcode2023/internal/stage"
)

//go:embed data/input.txt
var input []byte

type Matrix struct {
	Data  [][]byte
	Start Coord
	Max   Coord
}

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

func (m *Matrix) findByte(b byte) *Coord {
	for y := 0; y < len(m.Data); y++ {
		for x := 0; x < len(m.Data[y]); x++ {
			if m.Data[y][x] == b {
				return &Coord{x, y}
			}
		}
	}
	return nil
}

func (m *Matrix) AtCoord(c Coord) *byte {
	return m.At(c.X, c.Y)
}

func (m *Matrix) At(x, y int) *byte {
	if !m.In(x, y) {
		return nil
	}
	return &m.Data[y][x]
}

func (m *Matrix) SetAt(x, y int, b byte) {
	if !m.In(x, y) {
		panic("out of bounds")
	}
	m.Data[y][x] = b
}

func (m *Matrix) SetAtCoord(c Coord, b byte) {
	m.SetAt(c.X, c.Y, b)
}

func (m *Matrix) In(x, y int) bool {
	return x >= 0 && x <= m.Max.X && y >= 0 && y <= m.Max.Y
}

func (m *Matrix) String() string {
	sb := strings.Builder{}
	for y := 0; y < len(m.Data); y++ {
		sb.Write(m.Data[y])
		sb.WriteByte('\n')
	}
	return sb.String()
}

// .....
// ..U..
// .LSR.
// ..D..
// .....
//
// | is a vertical pipe connecting north and south.
// - is a horizontal pipe connecting east and west.
// L is a 90-degree bend connecting north and east.
// J is a 90-degree bend connecting north and west.
// 7 is a 90-degree bend connecting south and west.
// F is a 90-degree bend connecting south and east.
// . is ground; there is no pipe in this tile.
// S is the starting position of the animal; there is a pipe on this tile, but your sketch doesn't show what shape the pipe has.
var (
	ValidUps    = []byte{'|', '7', 'F', 'S'}
	ValidDowns  = []byte{'|', 'L', 'J', 'S'}
	ValidLefts  = []byte{'-', 'L', 'F', 'S'}
	ValidRights = []byte{'-', '7', 'J', 'S'}
)

func (m *Matrix) FindNexts(c Coord) []Coord {
	var nexts = make([]Coord, 0, 4)
	var nextCoord Coord
	var b *byte

	for _, tc := range []struct {
		Coord      Coord
		ValidBytes []byte
	}{
		{c.Up(), ValidUps},
		{c.Down(), ValidDowns},
		{c.Left(), ValidLefts},
		{c.Right(), ValidRights},
	} {
		nextCoord = tc.Coord
		b = m.AtCoord(nextCoord)
		if b != nil {
			i := bytes.IndexByte(tc.ValidBytes, *b)
			if i != -1 {
				nexts = append(nexts, nextCoord)
			}
		}
	}
	return nexts
}

type PipeRunner struct {
	Matrix *Matrix
	Coord  Coord
	Seen   map[Coord]struct{}
	Done   bool
	Steps  int
	Value  int
}

// func (p PipeRunner) Step() ([]PipeRunner, bool) {
// 	var runners []PipeRunner
// 	// var doneRunners []PipeRunner
// 	var found bool
// 	p.Seen[p.Coord] = struct{}{}
// 	for _, next := range p.Matrix.FindNexts(p.Coord) {
// 		if _, ok := p.Seen[next]; ok {
// 			if *p.Matrix.AtCoord(next) == 'S' {
// 				found = true
// 			}
// 			continue
// 		}
// 		runners = append(runners, PipeRunner{
// 			Matrix: p.Matrix,
// 			Coord:  next,
// 			Seen:   CopyMap(p.Seen),
// 			// Value:  p.Value + 1,
// 		})
// 	}
// 	// p.Matrix.SetAtCoord(p.Coord, 'x')
// 	return runners, found
// }

func (p PipeRunner) Step2() ([]PipeRunner, []PipeRunner) {
	var runners []PipeRunner
	var doneRunners []PipeRunner
	p.Seen[p.Coord] = struct{}{}
	for i, next := range p.Matrix.FindNexts(p.Coord) {
		if _, ok := p.Seen[next]; ok {
			if *p.Matrix.AtCoord(next) == 'S' {
				doneRunners = append(doneRunners, PipeRunner{
					Matrix: p.Matrix,
					Coord:  next,
					Seen:   CopyMap(p.Seen),
					Value:  p.Value + 1,
					Done:   true,
				})
			}
			continue
		}
		seen := p.Seen
		if i > 0 {
			seen = CopyMap(p.Seen)
		}
		runners = append(runners, PipeRunner{
			Matrix: p.Matrix,
			Coord:  next,
			// Seen:   CopyMap(p.Seen),
			Seen:  seen,
			Value: p.Value + 1,
		})
	}
	// p.Matrix.SetAtCoord(p.Coord, 'x')
	return runners, doneRunners
}

func CopyMap(m map[Coord]struct{}) map[Coord]struct{} {
	nm := make(map[Coord]struct{}, len(m)+1)
	for k := range m {
		nm[k] = struct{}{}
	}
	return nm
}

func FindLargestLoop(m *Matrix) PipeRunner {
	var runners = []PipeRunner{{
		Matrix: m,
		Coord:  m.Start,
		Seen:   make(map[Coord]struct{}),
		// Value:  0,
	}}
	var maxRunner PipeRunner
	// var counter = 0
	// go func() {
	// 	lastC := 0
	// 	for range time.NewTicker(time.Second).C {
	// 		fmt.Println("R:", len(runners), "S:", counter, "M:", maxCounter, "I:", counter-lastC)
	// 		lastC = counter
	// 	}
	// }()
	// newRunners := make([]PipeRunner, 0, len(runners))
	for {
		// newRunners = newRunners[:0]
		newRunners := make([]PipeRunner, 0, len(runners))
		for _, r := range runners {
			rns, drns := r.Step2()
			newRunners = append(newRunners, rns...)
			for _, drn := range drns {
				if drn.Value > maxRunner.Value {
					maxRunner = drn
				}
			}
		}
		if len(newRunners) == 0 {
			break
		}
		runners = newRunners
	}
	return maxRunner
}

func ParseInput(input io.Reader) *Matrix {
	matrix := &Matrix{}
	s := bufio.NewScanner(input)
	for s.Scan() {
		line := make([]byte, 0, len(s.Bytes())+2)
		line = append(line, '.')
		line = append(line, bytes.Clone(s.Bytes())...)
		line = append(line, '.')
		// matrix.Data = append(matrix.Data, bytes.Clone(s.Bytes()))
		matrix.Data = append(matrix.Data, line)
	}
	var mData [][]byte
	fakeLine := bytes.Repeat([]byte{'.'}, len(matrix.Data[0]))
	mData = append(mData, fakeLine)
	mData = append(mData, matrix.Data...)
	mData = append(mData, fakeLine)
	matrix.Data = mData
	start := matrix.findByte('S')
	if start == nil {
		panic("no start")
	}
	matrix.Start = *start
	matrix.Max = Coord{
		X: len(matrix.Data[0]) - 1,
		Y: len(matrix.Data) - 1,
	}
	return matrix
}

func (m *Matrix) findType(c Coord) byte {
	nexts := m.FindNexts(c)
	if len(nexts) != 2 {
		return 0
	}
	var (
		Up    = Coord{c.X, c.Y - 1}
		Down  = Coord{c.X, c.Y + 1}
		Left  = Coord{c.X - 1, c.Y}
		Right = Coord{c.X + 1, c.Y}
	)
	switch [2]Coord(nexts) {
	// Up Down
	case [2]Coord{Up, Down}:
		return '|'
	// Left Right
	case [2]Coord{Left, Right}:
		return '-'
	// Up Left
	case [2]Coord{Up, Left}:
		return 'J'
	// Up Right
	case [2]Coord{Up, Right}:
		return 'L'
	// Down Left
	case [2]Coord{Down, Left}:
		return '7'
	// Down Right
	case [2]Coord{Down, Right}:
		return 'F'
	}
	return 0
}

func Stage1(input io.Reader) (any, error) {
	m := ParseInput(input)
	// fmt.Println(m)
	// fmt.Println(m.Start, "->", m.FindNexts(m.Start))
	v := FindLargestLoop(m).Value / 2
	// fmt.Println(m)
	return v, nil
}

var (
	UpEdges   = []byte{'|', 'L', 'J'}
	DownEdges = []byte{'|', '7', 'F'}
)

func FindEnclosedTiles(p PipeRunner) int {
	// Replace Start with corresponding symbol
	startType := p.Matrix.findType(p.Matrix.Start)
	if startType == 0 {
		panic("invalid start type")
	}
	p.Matrix.SetAtCoord(p.Matrix.Start, startType)

	// Find Edges and check innership
	var counter = 0
	for y, line := range p.Matrix.Data {
		// fmt.Println(string(line))
		var ue, de = 0, 0
		for x, b := range line {
			_, isEdge := p.Seen[Coord{x, y}]
			// fmt.Print("(", string(b), ")")
			if isEdge && bytes.IndexByte(UpEdges, b) != -1 {
				ue += 1
			}
			if isEdge && bytes.IndexByte(DownEdges, b) != -1 {
				de += 1
			}
			// if inLoop && !isEdge {
			// 	counter++
			// }
			// fmt.Print(tp%2 != 0, " ")
			// fmt.Print(bp%2 != 0, " ")
			// fmt.Print(!isEdge, " ")
			if ue%2 != 0 && de%2 != 0 && !isEdge {
				// fmt.Print("I")
				counter++
				p.Matrix.SetAt(x, y, 'I')
			}
			// if isEdge {
			// 	p.Matrix.SetAt(x, y, 'E')
			// }
			// fmt.Print("\n")
			// if b == '-' {
			// 	continue
			// }
			// if isEdge && bytes.IndexByte(LeftEdges, b) != -1 {
			// 	p.Matrix.SetAt(x, y, '*')
			// 	inLoop = !inLoop
			// }
			// if isEdge && bytes.IndexByte(LeftEdges, b) != -1 {
			// 	p.Matrix.SetAt(x, y, 'i')
			// 	inLoop++
			// 	continue
			// }
			// if inLoop > 0 && isEdge && bytes.IndexByte(RightEdges, b) != -1 {
			// 	p.Matrix.SetAt(x, y, 'o')
			// 	inLoop--
			// 	continue
			// }
			// if inLoop && !isEdge {
			// 	counter++
			// 	p.Matrix.SetAt(x, y, 'I')
			// }
		}
		// inLoop = false
	}
	return counter
}

func Stage2(input io.Reader) (any, error) {
	m := ParseInput(input)
	// fmt.Println(m)
	// fmt.Println(m.Start, "->", m.FindNexts(m.Start))
	r := FindLargestLoop(m)
	v := FindEnclosedTiles(r)
	// fmt.Println(m)
	return v, nil
}

//go:embed data/example2.txt
var example []byte

func main() {
	stage.RunCLI(input, Stage1, Stage2)
	// stage.RunCLI(example, Stage1, Stage2)
}
