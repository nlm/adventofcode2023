package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"io"

	"github.com/nlm/adventofcode2023/internal/stage"
	"github.com/nlm/adventofcode2023/internal/tokenizer"
)

const (
	TZero tokenizer.Key = iota
	TOne
	TTwo
	TThree
	TFour
	TFive
	TSix
	TSeven
	TEight
	TNine
)

func ProcessLine(parser *tokenizer.Tokenizer, data []byte) (int, error) {
	stream := parser.Parse(data)
	var tokens []tokenizer.Key
	for stream.Scan() {
		if stream.Token() > 0 {
			tokens = append(tokens, stream.Token())
		}
	}
	if len(tokens) < 1 {
		return 0, errors.New("invalid input")
	}
	return int(tokens[0]*10 + tokens[len(tokens)-1]), nil
}

func ProcessInput(parser *tokenizer.Tokenizer, r io.Reader) (int, error) {
	var sum = 0
	s := bufio.NewScanner(r)
	for s.Scan() {
		v, err := ProcessLine(parser, s.Bytes())
		if err != nil {
			return 0, err
		}
		sum += v
	}
	return sum, nil
}

//go:embed data/input.txt
var input []byte

func Stage1Tokenizer() *tokenizer.Tokenizer {
	parser := tokenizer.New().WithOverlap(true)
	parser.DefineTokensString(TZero, "0")
	parser.DefineTokensString(TOne, "1")
	parser.DefineTokensString(TTwo, "2")
	parser.DefineTokensString(TThree, "3")
	parser.DefineTokensString(TFour, "4")
	parser.DefineTokensString(TFive, "5")
	parser.DefineTokensString(TSix, "6")
	parser.DefineTokensString(TSeven, "7")
	parser.DefineTokensString(TEight, "8")
	parser.DefineTokensString(TNine, "9")
	return parser
}

func Stage1() error {
	parser := Stage1Tokenizer()
	v, err := ProcessInput(parser, bytes.NewReader(input))
	if err != nil {
		return err
	}
	fmt.Println(v)
	return nil
}

func Stage2Tokenizer() *tokenizer.Tokenizer {
	parser := tokenizer.New().WithOverlap(true)
	parser.DefineTokensString(TZero, "0", "zero")
	parser.DefineTokensString(TOne, "1", "one")
	parser.DefineTokensString(TTwo, "2", "two")
	parser.DefineTokensString(TThree, "3", "three")
	parser.DefineTokensString(TFour, "4", "four")
	parser.DefineTokensString(TFive, "5", "five")
	parser.DefineTokensString(TSix, "6", "six")
	parser.DefineTokensString(TSeven, "7", "seven")
	parser.DefineTokensString(TEight, "8", "eight")
	parser.DefineTokensString(TNine, "9", "nine")
	return parser
}

func Stage2() error {
	parser := Stage2Tokenizer()
	v, err := ProcessInput(parser, bytes.NewReader(input))
	if err != nil {
		return err
	}
	fmt.Println(v)
	return nil
}

func main() {
	stage.RunCLI(Stage1, Stage2)
}
