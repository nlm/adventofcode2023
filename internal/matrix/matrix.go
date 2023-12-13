package matrix

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Matrix struct {
	Data [][]byte
	Len  Coord
}

func (m *Matrix) Clone() *Matrix {
	nm := &Matrix{
		Data: make([][]byte, len(m.Data)),
		Len:  m.Len,
	}
	for y := 0; y < m.Len.Y; y++ {
		nm.Data[y] = bytes.Clone(m.Data[y])
	}
	return nm
}

var ErrInconsistentGeometry = fmt.Errorf("inconsistent geometry")

func NewFromReader(input io.Reader) (*Matrix, error) {
	matrix := &Matrix{}
	s := bufio.NewScanner(input)
	cols := -1
	for s.Scan() {
		if cols != -1 && len(s.Bytes()) != cols {
			return nil, ErrInconsistentGeometry
		} else {
			cols = len(s.Bytes())
		}
		matrix.Data = append(matrix.Data, bytes.Clone(s.Bytes()))
	}
	if s.Err() != nil {
		return nil, s.Err()
	}
	matrix.Len.X = len(matrix.Data[0])
	matrix.Len.Y = len(matrix.Data)
	return matrix, nil
}

func (m Matrix) FindByte(b byte) (Coord, bool) {
	for y := 0; y < len(m.Data); y++ {
		for x := 0; x < len(m.Data[y]); x++ {
			if m.Data[y][x] == b {
				return Coord{x, y}, true
			}
		}
	}
	return Coord{}, false
}

func (m *Matrix) InsertLineBefore(y int, b byte) {
	yLen := m.Len.Y
	m.Data = append(m.Data, []byte{})
	for j := yLen; j > y; j-- {
		m.Data[j] = m.Data[j-1]
	}
	m.Data[y] = bytes.Repeat([]byte{b}, m.Len.X)
	m.Len.Y++
}

func (m *Matrix) InsertColumnBefore(x int, b byte) {
	xLen := m.Len.X
	for y := 0; y < m.Len.Y; y++ {
		m.Data[y] = append(m.Data[y], byte(0))
		for i := xLen; i > x; i-- {
			m.Data[y][i] = m.Data[y][i-1]
		}
		m.Data[y][x] = b
	}
	m.Len.X++
}

func (m Matrix) AtCoord(c Coord) byte {
	return m.At(c.X, c.Y)
}

func (m Matrix) At(x, y int) byte {
	return m.Data[y][x]
}

func (m *Matrix) SetAt(x, y int, b byte) {
	m.Data[y][x] = b
}

func (m *Matrix) SetAtCoord(c Coord, b byte) {
	m.SetAt(c.X, c.Y, b)
}

func (m Matrix) In(x, y int) bool {
	return x >= 0 && x <= m.Len.X-1 && y >= 0 && y <= m.Len.Y-1
}

func (m Matrix) InCoord(c Coord) bool {
	return c.X >= 0 && c.X <= m.Len.X-1 && c.Y >= 0 && c.Y <= m.Len.Y-1
}

func (m Matrix) String() string {
	sb := strings.Builder{}
	for y := 0; y < len(m.Data); y++ {
		sb.Write(m.Data[y])
		sb.WriteByte('\n')
	}
	return sb.String()
}
