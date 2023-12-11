package main

import (
	_ "embed"
	"fmt"
	"io"
)

func FindEmpty(m *Matrix, emptyb byte) ([]int, []int) {
	var (
		emptyX []int
		emptyY []int
	)
	// find empty lines
	for y := m.Len.Y - 1; y >= 0; y-- {
		var notempty bool
		for x := 0; x < m.Len.X; x++ {
			if m.At(x, y) != emptyb {
				notempty = true
				break
			}
		}
		if !notempty {
			emptyY = append(emptyY, y)
		}
	}
	// find empty columns
	for x := m.Len.X - 1; x >= 0; x-- {
		var notempty bool
		for y := 0; y < m.Len.Y; y++ {
			if m.At(x, y) != emptyb {
				notempty = true
				break
			}
		}
		if !notempty {
			emptyX = append(emptyX, x)
		}
	}
	return emptyX, emptyY
}

type Galaxy struct {
	Coord
	Id int
}

func (g Galaxy) String() string {
	return fmt.Sprintf("{G%d (%d, %d)}", g.Id, g.X, g.Y)
}

func FindGalaxyCoordinates(m Matrix) []Galaxy {
	id := 1
	var galaxies []Galaxy
	for y := 0; y < m.Len.Y; y++ {
		for x := 0; x < m.Len.X; x++ {
			if m.At(x, y) == '#' {
				galaxies = append(galaxies, Galaxy{
					Coord: Coord{x, y},
					Id:    id,
				})
				id++
			}
		}
	}
	return galaxies
}

func lowHigh(a, b int) (int, int) {
	if a < b {
		return a, b
	}
	return b, a
}

func PathLength(g1, g2 Galaxy, emptyX, emptyY []int, emptyVal int) int {
	lowX, highX := lowHigh(g1.X, g2.X)
	lowY, highY := lowHigh(g1.Y, g2.Y)
	diff := highX - lowX + highY - lowY
	for _, x := range emptyX {
		if x > lowX && x < highX {
			diff += emptyVal
		}
	}
	for _, y := range emptyY {
		if y > lowY && y < highY {
			diff += emptyVal
		}
	}
	return diff
}

func Stage(input io.Reader, expansionSize int) (any, error) {
	m := NewFromReader(input)
	emptyX, emptyY := FindEmpty(m, '.')
	gx := FindGalaxyCoordinates(*m)
	sum := 0
	for i := 0; i < len(gx); i++ {
		for j := i; j < len(gx); j++ {
			if i == j {
				continue
			}
			sum += PathLength(gx[i], gx[j], emptyX, emptyY, expansionSize)
		}
	}
	return sum, nil
}
func Stage1(input io.Reader) (any, error) {
	return Stage(input, 1)
}

func Stage2(input io.Reader) (any, error) {
	return Stage(input, 999999)
}
