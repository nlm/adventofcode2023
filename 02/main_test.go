package main

import (
	_ "embed"
	"testing"
)

//go:embed data/test_stage1.txt
var exampleData []byte

func TestParse(t *testing.T) {
	
}
