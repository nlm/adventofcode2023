package utils

import (
	"strconv"
	"strings"
)

func MustAtoi(s string) int {
	v, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		panic(err)
	}
	return v
}
