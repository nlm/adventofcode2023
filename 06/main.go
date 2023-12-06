package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"io"
	"regexp"

	"github.com/nlm/adventofcode2023/internal/stage"
	"github.com/nlm/adventofcode2023/internal/utils"
)

//go:embed data/input.txt
var input []byte

type RaceData struct {
	Len       int
	Durations []int
	Records   []int
}

func ParseInput(input io.Reader, trimSpaces bool) RaceData {
	rd := RaceData{}
	s := bufio.NewScanner(input)
	s.Scan()
	rd.Durations = ParseLine(s.Bytes(), trimSpaces)
	s.Scan()
	rd.Records = ParseLine(s.Bytes(), trimSpaces)
	rd.Len = len(rd.Records)
	return rd
}

var NumberRE = regexp.MustCompile(`\b\d+\b`)

func ParseLine(line []byte, trimSpaces bool) []int {
	var elts []int
	line = line[bytes.IndexByte(line, ':')+1:]
	if trimSpaces {
		line = bytes.ReplaceAll(line, []byte(" "), []byte(""))
	}
	for _, elt := range NumberRE.FindAll(line, -1) {
		elts = append(elts, utils.MustAtoi(string(elt)))
	}
	return elts
}

func (rd *RaceData) CountWinningHolds(race int) int {
	duration := rd.Durations[race]
	record := rd.Records[race]
	var winningHolds int
	for i := 0; i < duration; i++ {
		speed := i
		runtime := duration - i
		distance := speed * runtime
		if distance > record {
			winningHolds++
		}
	}
	return winningHolds
}

func (rd *RaceData) WaysToBeatRecords() int {
	var sum int = 1
	for i := 0; i < rd.Len; i++ {
		sum *= rd.CountWinningHolds(i)
	}
	return sum
}

func Stage1(input io.Reader) (any, error) {
	rd := ParseInput(input, false)
	return rd.WaysToBeatRecords(), nil
}

func Stage2(input io.Reader) (any, error) {
	rd := ParseInput(input, true)
	return rd.WaysToBeatRecords(), nil
}

func main() {
	stage.RunCLI(input, Stage1, Stage2)
}
