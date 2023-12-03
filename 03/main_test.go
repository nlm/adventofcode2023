package main

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLine(t *testing.T) {
	for _, tc := range []struct {
		Line    string
		Symbols []Symbol
		Parts   []Part
	}{
		{
			Line: `467..114..`,
			Parts: []Part{
				{
					Number:   467,
					Line:     1,
					BeginCol: 0,
					EndCol:   2,
				},
				{
					Number:   114,
					Line:     1,
					BeginCol: 5,
					EndCol:   7,
				},
			},
		},
		{
			Line: `...*......`,
			Symbols: []Symbol{
				{
					Col:    3,
					Line:   1,
					IsGear: true,
				},
			},
		},
		{
			Line: `617?......`,
			Parts: []Part{
				{
					Number:   617,
					BeginCol: 0,
					EndCol:   2,
					Line:     1,
				},
			},
			Symbols: []Symbol{
				{
					Col:    3,
					Line:   1,
					IsGear: false,
				},
			},
		},
		// `..35..633.`
		// ......#...
		// .....+.58.
		// ..592.....
		// ......755.
		// ...$.*....
		// .664.598..'
	} {
		t.Run(tc.Line, func(t *testing.T) {
			parts, symbols, err := ParseLine(1, []byte(tc.Line))
			assert.NoError(t, err)
			assert.Equal(t, tc.Parts, parts)
			assert.Equal(t, tc.Symbols, symbols)
		})
	}
}

// func TestIsAdjascent(t *testing.T) {
// }

//go:embed data/test_stage1.txt
var stage1Data []byte

func TestStage1(t *testing.T) {
	parts, symbols, err := Parse(bytes.NewReader(stage1Data))
	assert.NoError(t, err)
	sum := CalulatePartsSum(parts, symbols)
	assert.Equal(t, 4361, sum)
}

func TestStage2(t *testing.T) {
	parts, symbols, err := Parse(bytes.NewReader(stage1Data))
	assert.NoError(t, err)
	sum := CalulateGearRatio(parts, symbols)
	assert.Equal(t, 467835, sum)
}
