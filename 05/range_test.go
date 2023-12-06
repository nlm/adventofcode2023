package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConversionMap(t *testing.T) {
	cm := &ConversionMap{Ranges: make([]*Range, 0)}
	cm.Ranges = append(cm.Ranges, &Range{
		Source: 10,
		Dest:   20,
		Len:    5,
	})
	for _, tc := range []struct {
		Key   int
		Value int
	}{
		{0, 0},
		{9, 9},
		{10, 20},
		{14, 24},
		{15, 25},
		{16, 16},
		{25, 25},
	} {
		t.Run(fmt.Sprintf("%+v", tc.Key), func(t *testing.T) {
			assert.Equal(t, tc.Value, cm.Value(tc.Key))
		})
	}
}

func TestSeedRange(t *testing.T) {
	for _, tc := range []struct {
		Range  SeedRange
		Result []int
	}{
		{
			Range:  SeedRange{Start: 79, Length: 14},
			Result: []int{79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92},
		},
		{
			Range:  SeedRange{Start: 55, Length: 13},
			Result: []int{55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67},
		},
	} {
		t.Run(fmt.Sprint(tc.Range), func(t *testing.T) {
			var r []int
			for i := 0; i < tc.Range.Length; i++ {
				r = append(r, tc.Range.Start+i)
			}
			assert.Equal(t, tc.Result, r)
		})
	}
}

func TestResolveSeedLocation(t *testing.T) {
	seeds, maps, err := ParseInput(bytes.NewReader(example1))
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, []int{79, 14, 55, 13}, seeds)
	_ = maps
}
