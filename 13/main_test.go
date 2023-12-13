package main

import (
	"testing"

	"github.com/nlm/adventofcode2023/internal/stage"
)

var Stage1TestCases = []stage.TestCase{
	{
		Name:   "example",
		Result: 405,
	},
	{
		Name:   "input",
		Result: 37561,
	},
}

var Stage2TestCases = []stage.TestCase{
	{
		Name:   "example",
		Result: 400,
	},
	{
		Name:   "input",
		Result: 31108,
	},
}

// Do not edit below

func TestStage1(t *testing.T) {
	stage.Test(t, Stage1, Stage1TestCases)
}

func TestStage2(t *testing.T) {
	stage.Test(t, Stage2, Stage2TestCases)
}

func BenchmarkStage1(b *testing.B) {
	stage.Benchmark(b, Stage1, Stage1TestCases)
}

func BenchmarkStage2(b *testing.B) {
	stage.Benchmark(b, Stage2, Stage2TestCases)
}
