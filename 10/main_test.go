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

//go:embed data/example3.txt
var example3 []byte

//go:embed data/example4.txt
var example4 []byte

//go:embed data/example5.txt
var example5 []byte

//go:embed data/example6.txt
var example6 []byte

//go:embed data/example8.txt
var example8 []byte

//go:embed data/example9.txt
var example9 []byte

//go:embed data/example10.txt
var example10 []byte

func TestStage1(t *testing.T) {
	stage.Test(t, Stage1, []stage.TestCase{
		{
			Name:   "example1",
			Input:  example1,
			Result: 4,
		},
		{
			Name:   "example2",
			Input:  example2,
			Result: 8,
		},
		{
			Name:   "example3",
			Input:  example3,
			Result: 8,
		},
		{
			Name:   "example4",
			Input:  example4,
			Result: 8,
		},
		{
			Name:   "example5",
			Input:  example5,
			Result: 2,
		},
	})
}

func BenchmarkStage1(b *testing.B) {
	stage.Benchmark(b, Stage1, []stage.TestCase{
		{
			Name:   "example1",
			Input:  example1,
			Result: 4,
		},
		{
			Name:   "example2",
			Input:  example2,
			Result: 8,
		},
		{
			Name:   "example3",
			Input:  example3,
			Result: 8,
		},
		{
			Name:   "example4",
			Input:  example4,
			Result: 8,
		},
		{
			Name:   "example5",
			Input:  example5,
			Result: 2,
		},
		{
			Name:   "input",
			Input:  input,
			Result: 6942,
		},
	})
}

func TestStage2(t *testing.T) {
	stage.Test(t, Stage2, []stage.TestCase{
		{
			Name:   "example1",
			Input:  example1,
			Result: 1,
		},
		{
			Name:   "example6",
			Input:  example6,
			Result: 4,
		},
		{
			Name:   "example8",
			Input:  example8,
			Result: 10,
		},
		{
			Name:   "example9",
			Input:  example9,
			Result: 4,
		},
		{
			Name:   "example10",
			Input:  example10,
			Result: 7,
		},
		{
			Name:   "input",
			Input:  input,
			Result: 297,
		},
	})
}
