package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindCycle(t *testing.T) {
	data := []int{87, 69, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68}
	ref := []int{69, 69, 65, 64, 65, 63, 68}
	// data := []int{87,69,69,69} // bad?
	t.Run("small", func(t *testing.T) {
		size, offset := FindCycle(data[:20], 5, 10)
		assert.Equal(t, 0, size)
		assert.Equal(t, 0, offset)
	})
	t.Run("ref", func(t *testing.T) {
		size, offset := FindCycle(data, 1, 3)
		assert.Equal(t, 2, offset)
		assert.Equal(t, ref, data[offset:offset+size])
	})
}

func TestCycle(t *testing.T) {
	data := []int{87, 69, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68}
	ref := []int{69, 69, 65, 64, 65, 63, 68}
	size, offset := FindCycle(data, 5, 5)
	t.Run("ref", func(t *testing.T) {
		assert.Equal(t, ref, data[offset:offset+size])
	})
	t.Run("offset", func(t *testing.T) {
		assert.Equal(t, 2, offset)
	})
	for _, tc := range []struct {
		Idx   int
		Value int
	}{
		{1, 87},
		{2, 69},
		{3, 69},
		{4, 69},
		{5, 65},
		{10, 69},
		{12, 65},
		{13, 64},
		{17, 69},
	} {
		t.Run(fmt.Sprint(tc.Idx), func(t *testing.T) {
			assert.Equal(t, tc.Value, ValueAt(data, offset, size, tc.Idx))
		})
	}
}
