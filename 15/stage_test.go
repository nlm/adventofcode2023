package main

import (
	"io"
	"testing"

	"github.com/nlm/adventofcode2023/internal/stage"
	"github.com/nlm/adventofcode2023/internal/utils"
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
		{"qp-", true, Lens{"qp", 0}},
	} {
		t.Run(tc.Instruction, func(t *testing.T) {
			inst := ParseInstruction([]byte(tc.Instruction))
			assert.Equal(t, tc.Remove, inst.Remove)
			assert.Equal(t, tc.Lens, inst.Lens)
		})
	}
}

func HashRange(bytes []byte) int {
	currentValue := 0
	for _, b := range bytes {
		if b == '\n' || b == '\r' {
			continue
		}
		currentValue += int(b)
		currentValue *= 17
		currentValue %= 256
	}
	return currentValue
}
func HashIndex(bytes []byte) int {
	currentValue := 0
	for i := 0; i < len(bytes); i++ {
		if bytes[i] == '\n' || bytes[i] == '\r' {
			continue
		}
		currentValue += int(bytes[i])
		currentValue *= 17
		currentValue %= 256
	}
	return currentValue
}

func BenchmarkHash(b *testing.B) {
	data := utils.Must(io.ReadAll(stage.Open("example.txt")))
	for _, tc := range []struct {
		Name string
		Func func([]byte) int
	}{
		{"HashIndex", HashIndex},
		{"HashRange", HashRange},
	} {
		b.Run(tc.Name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				tc.Func(data)
			}
		})
	}
}
