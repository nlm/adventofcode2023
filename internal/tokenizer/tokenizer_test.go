package tokenizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ProcessLine(parser *Tokenizer, data []byte) ([]Key, error) {
	stream := parser.Parse(data)
	var tokens []Key
	for stream.Scan() {
		if stream.Token() > 0 {
			tokens = append(tokens, stream.Token())
		}
	}
	return tokens, nil
}

func TestTokenizer(t *testing.T) {
	p := New()
	p.DefineTokens(0, [][]byte{[]byte("0"), []byte("zero")})
	p.DefineTokens(1, [][]byte{[]byte("1"), []byte("one")})
	p.DefineTokens(2, [][]byte{[]byte("2"), []byte("two")})
	p.DefineTokens(3, [][]byte{[]byte("3"), []byte("three")})
	p.DefineTokens(4, [][]byte{[]byte("4"), []byte("four")})
	p.DefineTokens(5, [][]byte{[]byte("5"), []byte("five")})
	p.DefineTokens(6, [][]byte{[]byte("6"), []byte("six")})
	p.DefineTokens(7, [][]byte{[]byte("7"), []byte("seven")})
	p.DefineTokens(8, [][]byte{[]byte("8"), []byte("eight")})
	p.DefineTokens(9, [][]byte{[]byte("9"), []byte("nine")})
	for _, tc := range []struct {
		input []byte
		value []Key
	}{
		{[]byte("four9tbnqhjlbmqnjq4gpzpvjtl2"), []Key{4, 9, 4, 2}},
		{[]byte("8three75sevenbbsbxjscvseven6mhpx"), []Key{8, 3, 7, 5, 7, 7, 6}},
		{[]byte("fivetmxkjczpjninefive5pss3onetwonetmq"), []Key{5, 9, 5, 5, 3, 1, 2, 1}},
	} {
		v, err := ProcessLine(p, tc.input)
		if assert.NoError(t, err) {
			assert.Equal(t, tc.value, v)
		}
	}
}
