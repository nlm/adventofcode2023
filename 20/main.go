package main

import (
	"embed"

	"github.com/nlm/adventofcode2023/internal/stage"
)

// Do not edit this file

//go:embed data
var inputFS embed.FS

func init() {
	stage.SetFS(inputFS)
}

func main() {
	stage.RunCLI(nil, Stage1, Stage2, Stage3)
}
