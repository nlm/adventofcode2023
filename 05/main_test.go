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
			Result: 35,
		},
		{
			Name:   "input",
			Input:  input,
			Result: 324724204,
		},
	})
}

func TestStage2(t *testing.T) {
	stage.Test(t, Stage2, []stage.TestCase{
		{
			Name:   "example",
			Input:  example2,
			Result: 46,
		},
		// not testing input due to computational cost
	})
}

func TestStage3(t *testing.T) {
	stage.Test(t, Stage3, []stage.TestCase{
		{
			Name:   "example",
			Input:  example2,
			Result: 46,
		},
		// not testing input due to computational cost
	})
}
