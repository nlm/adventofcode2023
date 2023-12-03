package main

import (
	"bufio"
	_ "embed"
	"io"

	"github.com/nlm/adventofcode2023/internal/stage"
)

//go:embed data/input.txt
var input []byte

func Stage1(input io.Reader) (any, error) {
	ref := Draw{
		Red:   12,
		Green: 13,
		Blue:  14,
	}
	s := bufio.NewScanner(input)
	sum := 0
	for s.Scan() {
		id, draws, err := ParseLine(s.Bytes())
		if err != nil {
			return nil, err
		}
		possible := true
		for _, d := range draws {
			if ref.Red < d.Red || ref.Green < d.Green || ref.Blue < d.Blue {
				possible = false
				break
			}
		}
		if possible {
			sum += id
		}
	}
	return sum, nil
}

func Stage2(input io.Reader) (any, error) {
	s := bufio.NewScanner(input)
	sum := 0
	for s.Scan() {
		_, draws, err := ParseLine(s.Bytes())
		if err != nil {
			return nil, err
		}
		ref := Draw{}
		for _, d := range draws {
			if d.Red > ref.Red {
				ref.Red = d.Red
			}
			if d.Green > ref.Green {
				ref.Green = d.Green
			}
			if d.Blue > ref.Blue {
				ref.Blue = d.Blue
			}
		}
		sum += ref.Power()
	}
	return sum, nil
}

func main() {
	stage.RunCLI(input, Stage1, Stage2)
}
