package tokenizer

import "bytes"

type Key int

type Tokenizer struct {
	tokens map[Key][][]byte
}

func New() *Tokenizer {
	return &Tokenizer{
		tokens: make(map[Key][][]byte),
	}
}

func (t *Tokenizer) DefineTokens(key Key, tokens [][]byte) {
	t.tokens[key] = tokens
}

func (t *Tokenizer) Parse(data []byte) *Stream {
	return &Stream{
		buffer:    data,
		tokenizer: t,
		valid:     true,
		current:   -1,
	}
}

type Stream struct {
	tokenizer *Tokenizer
	buffer    []byte
	valid     bool
	current   Key
}

func (s *Stream) Scan() bool {
	if len(s.buffer) == 0 {
		return false
	}
	for key, values := range s.tokenizer.tokens {
		for _, value := range values {
			if bytes.HasPrefix(s.buffer, value) {
				s.current = key
				// s.buffer = s.buffer[len(value):]
				s.buffer = s.buffer[1:]
				return true
			}
		}
	}
	s.current = -1
	s.buffer = s.buffer[1:]
	return true
}

func (s *Stream) Token() Key {
	return s.current
}
