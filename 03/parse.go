package main

import (
	"bufio"
	"io"
	"math"
	"unicode"
)

type Part struct {
	Number   int
	BeginCol int
	EndCol   int
	Line     int
}

type Symbol struct {
	Col    int
	Line   int
	IsGear bool
}

func (p Part) IsAdjascent(s Symbol) bool {
	if math.Abs(float64(p.Line-s.Line)) > 1 {
		return false
	}
	if p.BeginCol-1 > s.Col {
		return false
	}
	if p.EndCol+1 < s.Col {
		return false
	}
	return true
}

func Parse(r io.Reader) ([]Part, []Symbol, error) {
	var (
		parts   []Part
		symbols []Symbol
	)
	s := bufio.NewScanner(r)
	lineNumber := 0
	for s.Scan() {
		p, s, err := ParseLine(lineNumber, s.Bytes())
		if err != nil {
			return nil, nil, err
		}
		parts = append(parts, p...)
		symbols = append(symbols, s...)
		lineNumber++
	}
	return parts, symbols, nil
}

func ParseLine(n int, line []byte) ([]Part, []Symbol, error) {
	var (
		inNumber   = false
		startIdx   = 0
		partNumber = 0
		parts      []Part
		symbols    []Symbol
	)
	// trick to simplify the problem
	line = append(line, '.')
	for i, r := range line {
		switch {
		case unicode.IsDigit(rune(r)):
			if !inNumber {
				partNumber = 0
				startIdx = i
				inNumber = true
			}
			partNumber *= 10
			partNumber += int(r - '0')
		default:
			if inNumber {
				parts = append(parts, Part{
					Number:   partNumber,
					Line:     n,
					BeginCol: startIdx,
					EndCol:   i - 1,
				})
				inNumber = false
			}
			// anything that is not a point is a symbol
			if r != '.' {
				symbols = append(symbols, Symbol{
					Col:    i,
					Line:   n,
					IsGear: r == '*',
				})
			}
		}
	}
	return parts, symbols, nil
}
