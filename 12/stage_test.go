package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSolve(t *testing.T) {
	for _, tc := range []struct {
		Buffer string
		Stats  [][]int
	}{
		{"?", [][]int{{}, {1}}},
	} {
		t.Run(tc.Buffer, func(t *testing.T) {
			assert.Equal(t, tc.Stats, Solve([]byte(tc.Buffer)))
		})
	}
}
