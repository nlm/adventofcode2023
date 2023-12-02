package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type Draw struct {
	Red   int
	Green int
	Blue  int
}

func (d Draw) Power() int {
	if d.Red+d.Green+d.Blue == 0 {
		return 0
	}
	v := 1
	if d.Red > 0 {
		v *= d.Red
	}
	if d.Green > 0 {
		v *= d.Green
	}
	if d.Blue > 0 {
		v *= d.Blue
	}
	return v
}

type Game struct {
	Id    int
	Draws []Draw
}

var (
	gameRe = regexp.MustCompile(`^Game (\d+): (.+)$`)
	drawRe = regexp.MustCompile(`^(\d+) (red|green|blue)$`)
)

func ParseLine(line []byte) (int, []Draw, error) {
	game := gameRe.FindSubmatch(line)
	if game == nil {
		return 0, nil, errors.New("game format mismatch")
	}
	id := MustAtoi(string(game[1]))
	draws, err := ParseDraws(string(game[2]))
	if err != nil {
		return 0, nil, err
	}
	return id, draws, nil
}

func ParseDraws(data string) ([]Draw, error) {
	var draws []Draw
	for _, draw := range strings.Split(data, "; ") {
		d := Draw{}
		for _, item := range strings.Split(draw, ", ") {
			items := drawRe.FindSubmatch([]byte(item))
			if items == nil {
				return nil, errors.New("error matching items")
			}
			v := MustAtoi(string(items[1]))
			switch string(items[2]) {
			case "red":
				d.Red = v
			case "green":
				d.Green = v
			case "blue":
				d.Blue = v
			}
		}
		draws = append(draws, d)
	}
	return draws, nil
}

func MustAtoi(s string) int {
	id, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return id
}
