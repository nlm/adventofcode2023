package main

import (
	_ "embed"
	"testing"

	"github.com/nlm/adventofcode2023/internal/stage"
)

//go:embed data/example.txt
var example []byte

func TestStage1(t *testing.T) {
	stage.Test(t, Stage1, []stage.TestCase{
		{
			Name:   "example",
			Input:  example,
			Result: 114,
		},
		{
			Name:   "input",
			Input:  input,
			Result: 1939607039,
		},
	})
}

func TestStage2(t *testing.T) {
	stage.Test(t, Stage2, []stage.TestCase{
		{
			Name:   "example",
			Input:  example,
			Result: 2,
		},
		{
			Name:   "input",
			Input:  input,
			Result: 1041,
		},
	})
}
