package main

import (
	_ "embed"
	"testing"

	"github.com/nlm/adventofcode2023/internal/stage"
)

//go:embed data/test_stage1.txt
var exampleData []byte

func TestStage1(t *testing.T) {
	stage.Test(t, Stage1, []stage.TestCase{
		{
			Name:   "example",
			Input:  exampleData,
			Result: 8,
			Err:    nil,
		},
		{
			Name:   "input",
			Input:  input,
			Result: 2348,
			Err:    nil,
		},
	})
}

func TestStage2(t *testing.T) {
	stage.Test(t, Stage2, []stage.TestCase{
		{
			Name:   "example",
			Input:  exampleData,
			Result: 2286,
			Err:    nil,
		},
		{
			Name:   "input",
			Input:  input,
			Result: 76008,
			Err:    nil,
		},
	})
}
