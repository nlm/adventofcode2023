package main

import (
	_ "embed"
	"testing"

	"github.com/nlm/adventofcode2023/internal/stage"
)

//go:embed data/example1.txt
var example1 []byte

//go:embed data/example2.txt
var example2 []byte

func TestStage1(t *testing.T) {
	stage.Test(t, Stage1, []stage.TestCase{
		{
			Name:   "example",
			Input:  example1,
			Result: 6,
		},
		{
			Name:   "input",
			Input:  input,
			Result: 15871,
		},
	})
}

func TestStage2(t *testing.T) {
	stage.Test(t, Stage2, []stage.TestCase{
		{
			Name:   "example",
			Input:  example2,
			Result: 6,
		},
		{
			Name:   "input",
			Input:  input,
			Result: 11283670395017,
		},
	})
}
