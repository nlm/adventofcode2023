package main

import (
	"testing"

	"github.com/nlm/adventofcode2023/internal/stage"
)

var Stage1TestCases = []stage.TestCase{
	{
		Name:   "example",
		Result: 46,
	},
	{
		Name:   "tsplit",
		Result: 9,
	},
	{
		Name:   "diagonal",
		Result: 18,
	},
	{
		Name:   "infinitefirst",
		Result: 16,
	},
	{
		Name:   "multiinf",
		Result: 41,
	},
}

var Stage2TestCases = []stage.TestCase{
	{
		Name:   "example",
		Result: 51,
	},
	{
		Name:   "input",
		Result: 7493,
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
