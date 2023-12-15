package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	for _, tc := range []struct {
		String string
		Hash   int
	}{
		{"rn=1", 30},
		{"cm-", 253},
		{"qp=3", 97},
	} {
		t.Run(tc.String, func(t *testing.T) {
			r := Hash([]byte(tc.String))
			assert.Equal(t, tc.Hash, r)
		})
	}

}

func TestParseInstruction(t *testing.T) {
	for _, tc := range []struct {
		Instruction string
		Remove      bool
		Lens        Lens
	}{
		{"rn=1", false, Lens{"rn", 1}},
		{"cm-", true, Lens{"cm", 0}},
		{"qp=3", false, Lens{"qp", 3}},
	} {
		t.Run(tc.Instruction, func(t *testing.T) {
			inst := ParseInstruction([]byte(tc.Instruction))
			assert.Equal(t, tc.Remove, inst.Remove)
			assert.Equal(t, tc.Lens, inst.Lens)
		})
	}

}
