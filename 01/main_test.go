package main

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed data/test_stage1.txt
var stageOneInput []byte

func TestStageOne(t *testing.T) {
	p := Tokenizer(false)
	s := bytes.NewReader(stageOneInput)
	v, err := ProcessInput(p, s)
	if assert.NoError(t, err) {
		assert.Equal(t, 142, v)
	}
}

//go:embed data/test_stage2.txt
var stageTwoInput []byte

func TestStageTwo(t *testing.T) {
	p := Tokenizer(true)
	s := bytes.NewReader(stageTwoInput)
	v, err := ProcessInput(p, s)
	if assert.NoError(t, err) {
		assert.Equal(t, 281, v)
	}
}

func TestOther(t *testing.T) {
	p := Tokenizer(true)
	for _, tc := range []struct {
		input []byte
		value int
	}{
		{[]byte("four9tbnqhjlbmqnjq4gpzpvjtl2"), 42},
		{[]byte("8three75sevenbbsbxjscvseven6mhpx"), 86},
		{[]byte("oneseveneightr6onecseven3"), 13},
		{[]byte("r81fourbnskcrnn1twobcfpvqqtdd"), 82},
		{[]byte("seven23"), 73},
		{[]byte("xttzsrkjjcvlgrm584qfjjhzlrhccj9"), 59},
		{[]byte("fivetmxkjczpjninefive5pss3onetwonetmq"), 51},
	} {
		v, err := ProcessLine(p, tc.input)
		if assert.NoError(t, err) {
			assert.Equal(t, tc.value, v)
		}
	}
}
