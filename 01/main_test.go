package main

import (
	_ "embed"
	"testing"

	"github.com/nlm/adventofcode2023/internal/stage"
	"github.com/stretchr/testify/assert"
)

//go:embed data/test_stage1.txt
var stageOneInput []byte

func TestStageOne(t *testing.T) {
	stage.Test(t, Stage1, []stage.TestCase{
		{
			Name:   "example",
			Input:  stageOneInput,
			Result: 142,
			Err:    nil,
		},
		{
			Name:   "input",
			Input:  input,
			Result: 55123,
			Err:    nil,
		},
	})
}

//go:embed data/test_stage2.txt
var stageTwoInput []byte

func TestStageTwo(t *testing.T) {
	stage.Test(t, Stage2, []stage.TestCase{
		{
			Name:   "example",
			Input:  stageTwoInput,
			Result: 281,
			Err:    nil,
		},
		{
			Name:   "input",
			Input:  input,
			Result: 55260,
			Err:    nil,
		},
	})
}

func TestStage2Tokenizer(t *testing.T) {
	p := Stage2Tokenizer()
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
