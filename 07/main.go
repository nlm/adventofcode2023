package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"sort"

	"github.com/nlm/adventofcode2023/internal/stage"
	"github.com/nlm/adventofcode2023/internal/utils"
)

//go:embed data/input.txt
var input []byte

type Entry struct {
	Hand  Hand
	Value int
	Bid   int
}

func (e Entry) String() string {
	return fmt.Sprintf(
		"([%2d %2d %2d %2d %2d] value=%x bid=%d)",
		e.Hand[0],
		e.Hand[1],
		e.Hand[2],
		e.Hand[3],
		e.Hand[4],
		e.Value,
		e.Bid,
	)
}

func ParseInput(input io.Reader, jokers bool) []Entry {
	var entries []Entry
	s := bufio.NewScanner(input)
	for s.Scan() {
		e := Entry{}
		b := s.Bytes()
		cards := b[0:5]
		for i := 0; i < len(cards); i++ {
			e.Hand[i] = CardValue(cards[i], jokers)
		}
		e.Bid = utils.MustAtoi(string(b[6:]))
		e.Value = e.Hand.Value()
		entries = append(entries, e)
	}
	return entries
}

func Stage(input io.Reader, jokers bool) (any, error) {
	entries := ParseInput(input, jokers)
	sort.Slice(entries, func(i, j int) bool {
		// if entries[i].Value == entries[j].Value {
		// 	panic("duplicate")
		// }
		return entries[i].Value < entries[j].Value
	})
	result := 0
	for i := 0; i < len(entries); i++ {
		// fmt.Println("E:", entries[i], "*", i+1, "=>", entries[i].Bid*(i+1))
		result += entries[i].Bid * (i + 1)
	}
	return result, nil
}

func Stage1(input io.Reader) (any, error) {
	return Stage(input, false)
}

func Stage2(input io.Reader) (any, error) {
	return Stage(input, true)
}

func main() {
	stage.RunCLI(input, Stage1, Stage2)
}
