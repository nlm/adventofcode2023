package main

import (
	_ "embed"
	"io"

	"github.com/nlm/adventofcode2023/internal/stage"
)

//go:embed data/input.txt
var input []byte

func Stage1(input io.Reader) (any, error) {
	return nil, stage.ErrUnimplemented
}

func Stage2(input io.Reader) (any, error) {
	return nil, stage.ErrUnimplemented
}

func main() {
	stage.RunCLI(input, Stage1, Stage2)
}
