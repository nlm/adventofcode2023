package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/nlm/adventofcode2023/internal/matrix"
	"github.com/nlm/adventofcode2023/internal/utils"
)

func FindYReflection(m *matrix.Matrix[byte], smudges int) int {
	var maxY = 0
	ref := m
	for y := 0; y < m.Len.Y; y++ {
		// clone for writing into the matrix for debugging
		m = ref.Clone()
		var asymetric bool
		var smudge int
		var ok bool
	outer:
		for diff := 0; y-diff >= 0 && y+1+diff < m.Len.Y; diff++ {
			ok = true
			for x := 0; x < m.Len.X; x++ {
				if m.At(x, y-diff) != m.At(x, y+1+diff) {
					smudge++
					if smudge > smudges {
						asymetric = true
						break outer
					}
				}
				m.SetAt(x, y-diff, 'U')
				m.SetAt(x, y+1+diff, 'D')
			}
		}
		if ok && !asymetric && smudge == smudges {
			maxY = y + 1
			fmt.Println(m, "Y:", maxY)
		}
	}
	return maxY
}

func FindXReflection(m *matrix.Matrix[byte], smudges int) int {
	var maxX = 0
	ref := m
	for x := 0; x < m.Len.X; x++ {
		// clone for writing into the matrix for debugging
		m = ref.Clone()
		var asymetric bool
		var smudge int
		var ok bool
	outer:
		for diff := 0; x-diff >= 0 && x+1+diff < m.Len.X; diff++ {
			ok = true
			for y := 0; y < m.Len.Y; y++ {
				if m.At(x-diff, y) != m.At(x+1+diff, y) {
					smudge++
					if smudge > smudges {
						asymetric = true
						break outer
					}
				}
				m.SetAt(x-diff, y, 'L')
				m.SetAt(x+1+diff, y, 'R')
			}
		}
		if ok && !asymetric && smudge == smudges {
			maxX = x + 1
			fmt.Println(m, "X:", maxX)
		}
	}
	return maxX
}

func Stage(input io.Reader, smudges int) (any, error) {
	matrixes := bytes.Split(utils.Must(io.ReadAll(input)), []byte("\n\n"))
	var (
		xr int
		yr int
	)
	for i := 0; i < len(matrixes); i++ {
		m, err := matrix.NewFromReader(bytes.NewReader(matrixes[i]))
		if err != nil {
			return nil, err
		}
		fmt.Println(m)
		xr += FindXReflection(m, smudges)
		yr += FindYReflection(m, smudges)
		fmt.Printf("\nSUM X:%d Y:%d\n\n", xr, yr)
	}
	return yr*100 + xr, nil
}

func Stage1(input io.Reader) (any, error) {
	return Stage(input, 0)
}

func Stage2(input io.Reader) (any, error) {
	return Stage(input, 1)
}
