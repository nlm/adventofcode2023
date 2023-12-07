package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestCardValue(t *testing.T) {
// 	assert.Equal(t, 2, int(CardValue('K')))
// 	assert.Equal(t, 13, int(CardValue('K')))
// }

// func TestCardsValue(t *testing.T) {
// 	assert.Equal(t, 0xddddd, Hand([]byte{13, 13, 13, 13, 13}).cardsValue())
// 	// assert.Equal(t, 13, int(CardValue('K')))
// }

func TestHandValue(t *testing.T) {
	for _, tc := range []struct {
		Hand  Hand
		Value int
	}{
		{
			Hand:  [5]byte{1, 1, 1, 1, 1},
			Value: 0x711111,
		},
		{
			Hand:  [5]byte{14, 1, 14, 1, 14},
			Value: 0x7e1e1e,
		},
		{
			Hand:  [5]byte{8, 1, 3, 3, 8},
			Value: 0x581338,
		},
	} {
		t.Run(fmt.Sprint(tc.Hand), func(t *testing.T) {
			assert.Equal(t, fmt.Sprintf("%x", tc.Value), fmt.Sprintf("%x", tc.Hand.Value()))
		})
	}
}
