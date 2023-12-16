package main

import (
	"fmt"
	"io"

	"github.com/nlm/adventofcode2023/internal/matrix"
	"github.com/nlm/adventofcode2023/internal/stage"
	"github.com/nlm/adventofcode2023/internal/utils"
)

func TiltMatrixUp(m *matrix.Matrix[byte]) {
	var moved bool
	for x := 0; x < m.Len.X; x++ {
		for {
			moved = false
			for y := 0; y < m.Len.Y-1; y++ {
				cur := matrix.Coord{X: x, Y: y}
				next := cur.Down()
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
				if m.At(x, y) == '.' && m.At(x-1, y) == 'O' {
					moved = true
					m.SetAt(x, y, 'O')
					m.SetAt(x-1, y, '.')
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
	if stage.Verbose() {
		fmt.Println(matrix.SMatrix(m))
	}
	TiltMatrixUp(m)
	if stage.Verbose() {
		fmt.Println(matrix.SMatrix(m))
	}
	return CalculateLoad(m), nil
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

// return size, offset
func FindCycle(slice []int, minCycleSize, minRepeats int) (int, int) {
	if minCycleSize <= 0 || minRepeats <= 0 {
		return 0, 0
	}
	minSliceSize := minCycleSize * minRepeats
	if len(slice) < minSliceSize {
		return 0, 0
	}
	for cycleSize := minCycleSize; cycleSize < len(slice)/minRepeats; cycleSize++ {
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
		}
		if matching {
			return cycleSize, len(slice) - cycleSize*minRepeats
		}
	}
	return 0, 0
}

func Stage2(input io.Reader) (any, error) {
	m := utils.Must(matrix.NewFromReader(input))
	if stage.Verbose() {
		fmt.Println(matrix.SMatrix(m))
	}
	cycle := []int{}
	nCycles := 1000000000
	for i := 0; i < nCycles; i++ {
		TiltMatrixUp(m)
		TiltMatrixLeft(m)
		TiltMatrixDown(m)
		TiltMatrixRight(m)
		cycle = append(cycle, CalculateLoad(m))
		size, offset := FindCycle(cycle, 10, 2)
		if size > 0 {
			return ValueAt(cycle, offset, size, nCycles), nil
		}
	}
	return -1, nil
}

func ValueAt(data []int, offset, size, idx int) int {
	if idx <= offset {
		return data[idx-1]
	}
	return data[((idx-offset-1)%size)+offset]
}
