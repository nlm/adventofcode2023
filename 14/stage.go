package main

import (
	"fmt"
	"hash/maphash"
	"io"

	"github.com/nlm/adventofcode2023/internal/matrix"
	"github.com/nlm/adventofcode2023/internal/utils"
)

const (
	Empty      = '.'
	RoundRock  = 'O'
	SquareRock = '#'
)

func TiltMatrixUp(m *matrix.Matrix[byte]) {
	var moved bool
	for x := 0; x < m.Len.X; x++ {
		for {
			moved = false
			for y := 0; y < m.Len.Y-1; y++ {
				cur := matrix.Coord{X: x, Y: y}
				next := cur.Down()
				// fmt.Println("CUR", string(m.AtCoord(cur)), "NEXT", string(m.AtCoord(next)))
				if !m.InCoord(next) {
					fmt.Println("out of bounds", next)
					continue
				}
				if m.AtCoord(cur) == '.' && m.AtCoord(next) == 'O' {
					moved = true
					m.SetAtCoord(cur, 'O')
					m.SetAtCoord(next, '.')
				}
			}
			if !moved {
				break
			}
		}
	}
}

func TiltMatrixDown(m *matrix.Matrix[byte]) {
	var moved bool
	for x := 0; x < m.Len.X; x++ {
		for {
			moved = false
			for y := m.Len.Y - 1; y > 0; y-- {
				cur := matrix.Coord{X: x, Y: y}
				next := cur.Up()
				// fmt.Println("CUR", string(m.AtCoord(cur)), "NEXT", string(m.AtCoord(next)))
				if !m.InCoord(next) {
					fmt.Println("out of bounds", next)
					continue
				}
				if m.AtCoord(cur) == '.' && m.AtCoord(next) == 'O' {
					moved = true
					m.SetAtCoord(cur, 'O')
					m.SetAtCoord(next, '.')
				}
			}
			if !moved {
				break
			}
		}
	}
}

func TiltMatrixLeft(m *matrix.Matrix[byte]) {
	var moved bool
	for y := 0; y < m.Len.Y; y++ {
		for {
			moved = false
			for x := 0; x < m.Len.X-1; x++ {
				cur := matrix.Coord{X: x, Y: y}
				next := cur.Right()
				// fmt.Println("CUR", string(m.AtCoord(cur)), "NEXT", string(m.AtCoord(next)))
				if !m.InCoord(next) {
					fmt.Println("out of bounds", next)
					continue
				}
				if m.AtCoord(cur) == '.' && m.AtCoord(next) == 'O' {
					moved = true
					m.SetAtCoord(cur, 'O')
					m.SetAtCoord(next, '.')
				}
			}
			if !moved {
				break
			}
		}
	}
}

func TiltMatrixRight(m *matrix.Matrix[byte]) {
	var moved bool
	for y := 0; y < m.Len.Y; y++ {
		for {
			moved = false
			for x := m.Len.X - 1; x > 0; x-- {
				cur := matrix.Coord{X: x, Y: y}
				next := cur.Left()
				// fmt.Println("CUR", string(m.AtCoord(cur)), "NEXT", string(m.AtCoord(next)))
				if !m.InCoord(next) {
					fmt.Println("out of bounds", next)
					continue
				}
				if m.AtCoord(cur) == '.' && m.AtCoord(next) == 'O' {
					moved = true
					m.SetAtCoord(cur, 'O')
					m.SetAtCoord(next, '.')
				}
			}
			if !moved {
				break
			}
		}
	}
}

func CalculateLoad(m *matrix.Matrix[byte]) int {
	var totalWeight int
	for y := 0; y < m.Len.Y; y++ {
		weight := m.Len.Y - y
		for x := 0; x < m.Len.X; x++ {
			if m.At(x, y) == 'O' {
				totalWeight += weight
			}
		}
	}
	return totalWeight
}

func Stage1(input io.Reader) (any, error) {
	m := utils.Must(matrix.NewFromReader(input))
	fmt.Println(matrix.SMatrix(m))
	TiltMatrixUp(m)
	fmt.Println(matrix.SMatrix(m))
	return CalculateLoad(m), nil
}

type Cycle struct {
	data   []uint64
	data2  []int
	values map[uint64]int
}

func NewCycle() *Cycle {
	return &Cycle{
		data:   []uint64{},
		data2:  []int{},
		values: map[uint64]int{},
	}
}

func (c *Cycle) Contains(k uint64) bool {
	for i := 0; i < len(c.data); i++ {
		if c.data[i] == k {
			return true
		}
	}
	return false
}

func (c *Cycle) Append(k uint64, v int) {
	c.data = append(c.data, k)
	c.data2 = append(c.data2, v)
	c.values[k] = v
}

func (c *Cycle) Value(k uint64) int {
	return c.values[k]
}

func SliceEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// [1 2 3 4 3 4]
// [1 2 3 4 1 2 3 4]
// [1 2 3 4 2 3 4]
func FindCycle(slice []int, minCycleSize, minRepeats int) []int {
	if minCycleSize <= 0 || minRepeats <= 0 {
		return nil
	}
	// fmt.Println("SIZE", len(slice))
	minSliceSize := minCycleSize * minRepeats
	if len(slice) < minSliceSize {
		return nil
	}
	for cycleSize := minCycleSize; cycleSize < len(slice)/minRepeats; cycleSize++ {
		// same := false
		lowerBound := len(slice) - cycleSize
		upperBound := len(slice)
		var lastSlice []int
		var matching bool
		for i := 0; i < minRepeats; i++ {
			matching = true
			subSlice := slice[lowerBound:upperBound]
			if lastSlice != nil && !SliceEqual(lastSlice, subSlice) {
				matching = false
				break
			}
			lastSlice = subSlice
			lowerBound -= cycleSize
			upperBound -= cycleSize
			// fmt.Println(slice, len(slice), (i+2)*cycleSize+1)
			// fmt.Println("L", cycleSize)
			// (i+2)*l +1 <= len(slice)
		}
		if matching {
			// fmt.Println(lastSlice)
			return lastSlice
		}
		// fmt.Println(slice)
	}
	return nil
}

func Stage2(input io.Reader) (any, error) {
	seed := maphash.MakeSeed()
	m := utils.Must(matrix.NewFromReader(input))
	fmt.Println(matrix.SMatrix(m))
	cycle := NewCycle()
	nCycles := 1000000000
	for i := 0; i < nCycles; i++ {
		TiltMatrixUp(m)
		TiltMatrixLeft(m)
		TiltMatrixDown(m)
		TiltMatrixRight(m)
		// fmt.Println()
		// fmt.Println("AFTER", i+1)
		// fmt.Println(matrix.SMatrix(m))
		h := maphash.Bytes(seed, m.Data)
		// fmt.Println(h, "->", CalculateLoad(m))
		cycle.Append(h, CalculateLoad(m))
		// 0 0 1 2 3 4 1 2 3 4
		theCycle := FindCycle(cycle.data2, 10, 10000)
		// 102065 too low
		// 106404 too high
		if theCycle != nil {
			fmt.Println(theCycle)
			offset := len(cycle.data2) % len(theCycle)
			return cycle.data2[nCycles%len(theCycle)+offset], nil
		}
	}
	return -1, nil
}
