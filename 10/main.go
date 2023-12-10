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
	Data [][]byte
	Max  Coord
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

func (m *Matrix) AtCoord(c Coord) byte {
	return m.At(c.X, c.Y)
}

func (m *Matrix) At(x, y int) byte {
	return m.Data[y][x]
}

func (m *Matrix) SetAt(x, y int, b byte) {
	m.Data[y][x] = b
}

func (m *Matrix) SetAtCoord(c Coord, b byte) {
	m.SetAt(c.X, c.Y, b)
}

func (m *Matrix) In(x, y int) bool {
	return x >= 0 && x <= m.Max.X && y >= 0 && y <= m.Max.Y
}

func (m *Matrix) InCoord(c Coord) bool {
	return c.X >= 0 && c.X <= m.Max.X && c.Y >= 0 && c.Y <= m.Max.Y
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

// findNexts finds the possible moves from a position in the matrix,
// depending on the type of edge at coordinate 'c'.
func (m *Matrix) FindNexts(c Coord) []Coord {
	var nexts = make([]Coord, 0, 4)
	var nextCoord Coord

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
		if m.InCoord(nextCoord) && bytes.IndexByte(tc.ValidBytes, m.AtCoord(nextCoord)) != -1 {
			nexts = append(nexts, nextCoord)
		}
	}
	return nexts
}

// findType finds the type of a maze slot from neighbours values.
func (m *Matrix) findType(c Coord) byte {
	nexts := m.FindNexts(c)
	if len(nexts) != 2 {
		return 0
	}
	var (
		Up    = c.Up()
		Down  = c.Down()
		Left  = c.Left()
		Right = c.Right()
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

type PipeRunner struct {
	Matrix  *Matrix
	Current Coord
	Seen    map[Coord]struct{}
	Count   int
}

func FindLoop(m *Matrix, origin Coord) PipeRunner {
	var pr = PipeRunner{
		Matrix:  m,
		Current: origin,
		Seen:    make(map[Coord]struct{}, m.Max.X*m.Max.Y),
		Count:   0,
	}
	for {
		pr.Seen[pr.Current] = struct{}{}
		// Find Next
		var nexts []Coord
		switch pr.Matrix.AtCoord(pr.Current) {
		case '|':
			nexts = []Coord{pr.Current.Up(), pr.Current.Down()}
		case '-':
			nexts = []Coord{pr.Current.Left(), pr.Current.Right()}
		case 'J':
			nexts = []Coord{pr.Current.Up(), pr.Current.Left()}
		case '7':
			nexts = []Coord{pr.Current.Down(), pr.Current.Left()}
		case 'L':
			nexts = []Coord{pr.Current.Up(), pr.Current.Right()}
		case 'F':
			nexts = []Coord{pr.Current.Down(), pr.Current.Right()}
		case 'S':
			nexts = []Coord{}
			if bytes.IndexByte(ValidUps, pr.Matrix.AtCoord(pr.Current.Up())) != -1 {
				nexts = append(nexts, pr.Current.Up())
			}
			if bytes.IndexByte(ValidDowns, pr.Matrix.AtCoord(pr.Current.Down())) != -1 {
				nexts = append(nexts, pr.Current.Down())
			}
			if bytes.IndexByte(ValidLefts, pr.Matrix.AtCoord(pr.Current.Left())) != -1 {
				nexts = append(nexts, pr.Current.Left())
			}
			if bytes.IndexByte(ValidRights, pr.Matrix.AtCoord(pr.Current.Right())) != -1 {
				nexts = append(nexts, pr.Current.Right())
			}
		}

		var found bool
		for _, next := range nexts {
			// Skipping if we've already seen this coordinate.
			if _, ok := pr.Seen[next]; ok {
				// Check if we're not back to the start.
				if !(next == origin && pr.Count > 1) {
					continue
				}
			}
			// This is only aesthetic
			// m.SetAtCoord(pr.Current, 'X')
			pr.Count++
			pr.Current = next
			found = true
			break
		}
		// fmt.Println(m)
		if pr.Current == origin {
			return pr
		}
		if !found {
			panic("dead end")
		}
	}
}

func ParseInput(input io.Reader) (*Matrix, Coord) {
	matrix := &Matrix{}
	s := bufio.NewScanner(input)

	// Surround matrix with '.' to simplify the problem
	matrix.Data = append(matrix.Data, []byte{})
	for s.Scan() {
		line := bytes.NewBuffer(make([]byte, 0, len(s.Bytes())+2))
		line.WriteByte('.')
		line.Write(s.Bytes())
		line.WriteByte('.')
		matrix.Data = append(matrix.Data, line.Bytes())
	}
	matrix.Data = append(matrix.Data, []byte{})

	// Replace first and last placeholders with dotLines
	dotLine := bytes.Repeat([]byte{'.'}, len(matrix.Data[1]))
	matrix.Data[0] = dotLine
	matrix.Data[len(matrix.Data)-1] = dotLine

	// Find origin
	origin := matrix.findByte('S')
	if origin == nil {
		panic("no origin found")
	}
	matrix.Max = Coord{
		X: len(matrix.Data[0]) - 1,
		Y: len(matrix.Data) - 1,
	}
	return matrix, *origin
}

func Stage1(input io.Reader) (any, error) {
	m, origin := ParseInput(input)
	v := FindLoop(m, origin).Count / 2
	// fmt.Println(m)
	return v, nil
}

var (
	UpEdges   = []byte{'|', 'L', 'J'}
	DownEdges = []byte{'|', '7', 'F'}
)

func FindEnclosedTiles(p PipeRunner, origin Coord) int {
	// Replace Start with corresponding symbol
	originType := p.Matrix.findType(origin)
	if originType == 0 {
		panic("invalid origin type")
	}
	p.Matrix.SetAtCoord(origin, originType)

	// Find Edges and check innership
	var counter = 0
	for y, line := range p.Matrix.Data {
		// fmt.Println(string(line))
		var ue, de = 0, 0
		for x, b := range line {
			_, isEdge := p.Seen[Coord{x, y}]
			if isEdge && bytes.IndexByte(UpEdges, b) != -1 {
				ue += 1
			}
			if isEdge && bytes.IndexByte(DownEdges, b) != -1 {
				de += 1
			}
			if ue%2 != 0 && de%2 != 0 && !isEdge {
				counter++
				p.Matrix.SetAt(x, y, 'I')
			}
		}
	}
	return counter
}

func Stage2(input io.Reader) (any, error) {
	m, origin := ParseInput(input)
	r := FindLoop(m, origin)
	v := FindEnclosedTiles(r, origin)
	return v, nil
}

func main() {
	stage.RunCLI(input, Stage1, Stage2)
}
