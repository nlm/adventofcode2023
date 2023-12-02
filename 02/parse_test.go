package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLine(t *testing.T) {
	for _, tc := range []struct {
		value []byte
		id    int
		draws []Draw
	}{
		{
			[]byte("Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green"),
			1,
			[]Draw{{Blue: 3, Red: 4}, {Red: 1, Green: 2, Blue: 6}, {Green: 2}},
		},
		{
			[]byte("Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue"),
			2,
			[]Draw{{Blue: 1, Green: 2}, {Green: 3, Blue: 4, Red: 1}, {Green: 1, Blue: 1}},
		},
		{
			[]byte("Game 11: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red"),
			11,
			[]Draw{{Green: 8, Blue: 6, Red: 20}, {Blue: 5, Red: 4, Green: 13}, {Green: 5, Red: 1}},
		},
	} {
		t.Run(string(tc.value), func(t *testing.T) {
			id, draws, err := ParseLine(tc.value)
			assert.NoError(t, err)
			assert.Equal(t, tc.id, id)
			assert.Equal(t, tc.draws, draws)
		})
	}
}
