package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"strings"

	"github.com/nlm/adventofcode2023/internal/stage"
	"github.com/nlm/adventofcode2023/internal/utils"
)

//go:embed data/input.txt
var input []byte

type Sequence []int
type History []Sequence

func (seq Sequence) IsAllZeroes() bool {
	for i := 0; i < len(seq); i++ {
		if seq[i] != 0 {
			return false
		}
	}
	return true
}

func (seq Sequence) Diffs() Sequence {
	var diffs = make(Sequence, 0, len(seq)-1)
	for i, last := 1, seq[0]; i < len(seq); i++ {
		diffs = append(diffs, seq[i]-last)
		last = seq[i]
	}
	return diffs
}

func (seq Sequence) LastValue() int {
	return seq[len(seq)-1]
}

func (seq Sequence) FirstValue() int {
	return seq[0]
}

func (seq Sequence) GenHistory() History {
	var hist = make(History, 0, 32)
	hist = append(hist, seq)
	for !seq.IsAllZeroes() {
		seq = seq.Diffs()
		hist = append(hist, seq)
	}
	return hist
}

func (hist History) String() string {
	var (
		str strings.Builder
		pad int
	)
	for i := 0; i < len(hist); i++ {
		fmt.Fprint(&str, strings.Repeat(" ", pad))
		for j := 0; j < len(hist[i]); j++ {
			fmt.Fprintf(&str, "%-2d  ", hist[i][j])
		}
		fmt.Fprint(&str, "\n")
		pad += 2
	}
	return str.String()
}

func (hist History) Extrapolate() History {
	histLast := len(hist) - 1
	hist[histLast] = append(hist[histLast], 0)
	hist[histLast] = append(Sequence{0}, hist[histLast]...)
	for i := histLast - 1; i >= 0; i-- {
		hist[i] = append(hist[i], hist[i].LastValue()+hist[i+1].LastValue())
		hist[i] = append(Sequence{hist[i].FirstValue() - hist[i+1].FirstValue()}, hist[i]...)
	}
	return hist
}

func (hist History) LastValue() int {
	return hist[0].LastValue()
}

func (hist History) FirstValue() int {
	return hist[0].FirstValue()
}

func ParseInput(input io.Reader) []Sequence {
	s := bufio.NewScanner(input)
	var sequences []Sequence
	for s.Scan() {
		fields := strings.Fields(s.Text())
		values := make([]int, 0, len(fields))
		for _, v := range fields {
			values = append(values, utils.MustAtoi(v))
		}
		sequences = append(sequences, values)
	}
	return sequences
}

func Stage1(input io.Reader) (any, error) {
	sequences := ParseInput(input)
	sum := 0
	for _, seq := range sequences {
		hist := seq.GenHistory()
		// fmt.Println(hist)
		hist.Extrapolate()
		// fmt.Println(hist)
		sum += hist.LastValue()
	}
	return sum, nil
}

func Stage2(input io.Reader) (any, error) {
	sequences := ParseInput(input)
	sum := 0
	for _, seq := range sequences {
		hist := seq.GenHistory()
		// fmt.Println(hist)
		hist.Extrapolate()
		// fmt.Println(hist)
		sum += hist.FirstValue()
	}
	return sum, nil
}

func main() {
	stage.RunCLI(input, Stage1, Stage2)
}
