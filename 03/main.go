package main

import (
	_ "embed"
	"io"

	"github.com/nlm/adventofcode2023/internal/stage"
)

//go:embed data/input.txt
var input []byte

func CalulatePartsSum(parts []Part, symbols []Symbol) int {
	sum := 0
	for _, p := range parts {
		for _, s := range symbols {
			if p.IsAdjascent(s) {
				sum += p.Number
			}
		}
	}
	return sum
}

func CalulateGearRatio(parts []Part, symbols []Symbol) int {
	sum := 0
	for _, s := range symbols {
		if !s.IsGear {
			continue
		}
		var adjParts []Part
		for _, p := range parts {
			if p.IsAdjascent(s) {
				adjParts = append(adjParts, p)
			}
		}
		if len(adjParts) == 2 {
			sum += adjParts[0].Number * adjParts[1].Number
		}
	}
	return sum
}

func Stage1(input io.Reader) (any, error) {
	parts, symbols, err := Parse(input)
	if err != nil {
		return nil, err
	}
	return CalulatePartsSum(parts, symbols), nil
}

func Stage2(input io.Reader) (any, error) {
	parts, symbols, err := Parse(input)
	if err != nil {
		return nil, err
	}
	return CalulateGearRatio(parts, symbols), nil
}

func main() {
	stage.RunCLI(input, Stage1, Stage2)
}
