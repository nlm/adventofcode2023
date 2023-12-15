package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindCycle(t *testing.T) {
	data := []int{87, 69, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68, 69, 69, 65, 64, 65, 63, 68}
	ref := []int{69, 69, 65, 64, 65, 63, 68}
	// data := []int{87,69,69,69} // bad?
	t.Run("small", func(t *testing.T) {
		assert.Nil(t, FindCycle(data, 200, 300))
	})
	t.Run("ref", func(t *testing.T) {
		cycle := FindCycle(data, 1, 3)
		assert.Equal(t, ref, cycle)
		assert.Equal(t, 2, len(data)%len(cycle))
	})
}
