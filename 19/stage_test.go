package main

import (
	"testing"

	"github.com/nlm/adventofcode2023/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestRangeSplit(t *testing.T) {
	for _, tc := range []struct {
		Condition    string
		RangeMatch   Ranges
		RangeNoMatch Ranges
	}{
		{
			Condition:    "x<2000",
			RangeMatch:   [8]int{1, 1999, 1, 4000, 1, 4000, 1, 4000},
			RangeNoMatch: [8]int{2000, 4000, 1, 4000, 1, 4000, 1, 4000},
		},
		{
			Condition:    "x>2000",
			RangeMatch:   [8]int{2001, 4000, 1, 4000, 1, 4000, 1, 4000},
			RangeNoMatch: [8]int{1, 2000, 1, 4000, 1, 4000, 1, 4000},
		},
		{
			Condition:    "s<1000",
			RangeMatch:   [8]int{1, 4000, 1, 4000, 1, 4000, 1, 999},
			RangeNoMatch: [8]int{1, 4000, 1, 4000, 1, 4000, 1000, 4000},
		},
		{
			Condition:    "s>1000",
			RangeMatch:   [8]int{1, 4000, 1, 4000, 1, 4000, 1001, 4000},
			RangeNoMatch: [8]int{1, 4000, 1, 4000, 1, 4000, 1, 1000},
		},
	} {
		t.Run(tc.Condition, func(t *testing.T) {
			cond := utils.Must(ParseCondition([]byte(tc.Condition)))
			ranges := NewRanges()
			rMatch, rNoMatch := ranges.Split(*cond)
			t.Run("Match", func(t *testing.T) {
				assert.Equal(t, tc.RangeMatch, rMatch)
			})
			t.Run("NoMatch", func(t *testing.T) {
				assert.Equal(t, tc.RangeNoMatch, rNoMatch)
			})
		})
	}
}
