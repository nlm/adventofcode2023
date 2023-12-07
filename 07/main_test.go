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
			Result: 6440,
		},
		{
			Name:   "example",
			Input:  input,
			Result: 247961593,
		},
	})
}

func TestStage2(t *testing.T) {
	stage.Test(t, Stage2, []stage.TestCase{
		{
			Name:   "example",
			Input:  example,
			Result: 5905,
		},
		{
			Name:   "input",
			Input:  input,
			Result: 248750699,
		},
	})
}
