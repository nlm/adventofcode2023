package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	// "github.com/bzick/tokenizer"
	"aoc2023/internal/tokenizer"
)

var flagFile = flag.String("file", "input.txt", "filename")
var flagStageTwo = flag.Bool("stage2", false, "stage2")

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

func Tokenizer(stageTwo bool) *tokenizer.Tokenizer {
	parser := tokenizer.New()
	if stageTwo {
		parser.DefineTokens(TZero, [][]byte{[]byte("0"), []byte("zero")})
		parser.DefineTokens(TOne, [][]byte{[]byte("1"), []byte("one")})
		parser.DefineTokens(TTwo, [][]byte{[]byte("2"), []byte("two")})
		parser.DefineTokens(TThree, [][]byte{[]byte("3"), []byte("three")})
		parser.DefineTokens(TFour, [][]byte{[]byte("4"), []byte("four")})
		parser.DefineTokens(TFive, [][]byte{[]byte("5"), []byte("five")})
		parser.DefineTokens(TSix, [][]byte{[]byte("6"), []byte("six")})
		parser.DefineTokens(TSeven, [][]byte{[]byte("7"), []byte("seven")})
		parser.DefineTokens(TEight, [][]byte{[]byte("8"), []byte("eight")})
		parser.DefineTokens(TNine, [][]byte{[]byte("9"), []byte("nine")})
	} else {
		parser.DefineTokens(TZero, [][]byte{[]byte("0")})
		parser.DefineTokens(TOne, [][]byte{[]byte("1")})
		parser.DefineTokens(TTwo, [][]byte{[]byte("2")})
		parser.DefineTokens(TThree, [][]byte{[]byte("3")})
		parser.DefineTokens(TFour, [][]byte{[]byte("4")})
		parser.DefineTokens(TFive, [][]byte{[]byte("5")})
		parser.DefineTokens(TSix, [][]byte{[]byte("6")})
		parser.DefineTokens(TSeven, [][]byte{[]byte("7")})
		parser.DefineTokens(TEight, [][]byte{[]byte("8")})
		parser.DefineTokens(TNine, [][]byte{[]byte("9")})
	}
	return parser
}

func ProcessLine(parser *tokenizer.Tokenizer, data []byte) (int, error) {
	stream := parser.Parse([]byte(strings.ToLower(string(data))))
	var tokens []tokenizer.Key
	for stream.Scan() {
		if stream.Token() > 0 {
			tokens = append(tokens, stream.Token())
		}
	}
	if len(tokens) < 1 {
		return 0, errors.New("invalid input")
	}
	v, err := strconv.Atoi(fmt.Sprintf("%d%d", tokens[0], tokens[len(tokens)-1]))
	if err != nil {
		return 0, err
	}
	return v, nil
}

func ProcessInput(parser *tokenizer.Tokenizer, r io.Reader) (int, error) {
	var sum = 0
	s := bufio.NewScanner(r)
	for s.Scan() {
		v, err := ProcessLine(parser, s.Bytes())
		if err != nil {
			return 0, err
		}
		fmt.Println(s.Text(), "->", v)
		sum += v
	}
	return sum, nil
}

func main() {
	flag.Parse()

	// Parser
	parser := Tokenizer(*flagStageTwo)

	// File
	f, err := os.Open(*flagFile)
	if err != nil {
		log.Fatal(err)
	}

	// Scan
	v, err := ProcessInput(parser, f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(v)
}
