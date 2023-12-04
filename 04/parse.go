package main

import (
	"math"
	"strings"

	"github.com/nlm/adventofcode2023/internal/utils"
)

type Card struct {
	Id             int
	WinningNumbers []int
	Numbers        []int
}

func (c Card) CurrentWinningNumbers() []int {
	var win []int
	for _, n := range c.Numbers {
		for _, w := range c.WinningNumbers {
			if n == w {
				win = append(win, n)
				continue
			}
		}
	}
	return win
}

func (c Card) Value() int {
	winCount := c.WinCount()
	if winCount > 0 {
		return int(math.Pow(2, float64(winCount-1)))
	}
	return 0
}

func (c Card) WinCount() int {
	return len(c.CurrentWinningNumbers())
}

func ParseLine(line []byte) Card {
	l := string(line)
	items := strings.Split(l, ":")
	cardNumber := items[0][5:]
	numbers := strings.Split(items[1], "|")
	winningNumbers := strings.Split(strings.ReplaceAll(numbers[0], "  ", " "), " ")
	cardNumbers := strings.Split(strings.ReplaceAll(numbers[1], "  ", " "), " ")
	card := Card{
		Id: utils.MustAtoi(cardNumber),
	}
	for _, n := range winningNumbers {
		if len(n) > 0 {
			card.WinningNumbers = append(card.WinningNumbers, utils.MustAtoi(n))
		}
	}
	for _, n := range cardNumbers {
		if len(n) > 0 {
			card.Numbers = append(card.Numbers, utils.MustAtoi(n))
		}
	}
	return card
}
