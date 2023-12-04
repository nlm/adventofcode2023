package main

import (
	"bufio"
	_ "embed"
	"io"

	"github.com/nlm/adventofcode2023/internal/stage"
)

//go:embed data/input.txt
var input []byte

func Stage1(input io.Reader) (any, error) {
	s := bufio.NewScanner(input)
	points := 0
	for s.Scan() {
		card := ParseLine(s.Bytes())
		points += card.Value()
	}
	return points, nil
}

func Stage2(input io.Reader) (any, error) {
	s := bufio.NewScanner(input)
	var cards []Card
	for s.Scan() {
		cards = append(cards, ParseLine(s.Bytes()))
	}
	return CountMultiCards(cards), nil
}

func CountMultiCards(cards []Card) int {
	var multipliers = make([]int, len(cards))
	for ci, card := range cards {
		multipliers[ci] += 1
		for i := 0; i < card.WinCount(); i++ {
			multipliers[ci+1+i] += multipliers[ci]
		}
	}
	sum := 0
	for _, value := range multipliers {
		sum += value
	}
	return sum
}

func main() {
	stage.RunCLI(input, Stage1, Stage2)
}
