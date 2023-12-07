package main

import (
	_ "embed"
	"sort"
)

type Hand [5]byte

const (
	HighCard     = 0x100000
	OnePair      = 0x200000
	TwoPairs     = 0x300000
	ThreeOfAKind = 0x400000
	FullHouse    = 0x500000
	FourOfAKind  = 0x600000
	FiveOfAKind  = 0x700000
)

func CardValue(b byte, joker bool) byte {
	if b >= '1' && b <= '9' {
		return b - '0'
	}
	switch b {
	case 'T':
		return 0xa
	case 'J':
		if joker {
			return 0x1
		}
		return 0xb
	case 'Q':
		return 0xc
	case 'K':
		return 0xd
	case 'A':
		return 0xe
	}
	panic("oops")
}

func (h Hand) cardsValue() int {
	sum := 0
	for i := 0; i < len(h); i++ {
		sum <<= 4
		sum += int(h[i])
	}
	return sum
}

// func (h Hand) PokerValue() int {
// 	cards := make(map[byte]int)
// 	for i := 0; i < len(h); i++ {
// 		cards[h[i]]++
// 	}
// 	for k, v := range cards {
// 		if v == 5 {
// 			return FiveOfAKind + int(k)<<16
// 		}
// 	}
// 	for k, v := range cards {
// 		if v == 4 {
// 			for k2, v2 := range cards {
// 				if v2 == 1 {
// 					return FourOfAKind + int(k)<<16 + int(k2)<<12
// 				}
// 			}
// 			return -1
// 		}
// 	}
// 	for k, v := range cards {
// 		if v == 3 {
// 			for k2, v2 := range cards {
// 				if v2 == 2 {
// 					return FullHouse + int(k)<<16 + int(k2)<<12
// 				}
// 			}
// 			var next []int
// 			for k2, v2 := range cards {
// 				if v2 == 1 {
// 					next = append(next, int(k2))
// 				}
// 			}
// 			if len(next) != 2 {
// 				return -1
// 			}
// 			if next[0] > next[1] {
// 				return ThreeOfAKind + int(k)<<16 + int(next[0])<<12 + int(next[1])<<8
// 			}
// 			return ThreeOfAKind + int(k)<<16 + int(next[1])<<12 + int(next[0])<<8
// 		}
// 	}
// 	for k, v := range cards {
// 		if v == 2 {
// 			for k2, v2 := range cards {
// 				if v2 == 2 && k != k2 {
// 					next := 0
// 					for k3, v3 := range cards {
// 						if v3 == 1 {
// 							next = int(k3)
// 						}
// 					}
// 					if next == 0 {
// 						return -1
// 					}
// 					if k > k2 {
// 						return TwoPairs + int(k)<<16 + int(k2)<<12 + int(next)<<8
// 					}
// 					return TwoPairs + int(k2)<<16 + int(k)<<12 + int(next)<<8
// 				}
// 			}
// 			var next []int
// 			for k2, v2 := range cards {
// 				if v2 == 1 {
// 					next = append(next, int(k2))
// 				}
// 			}
// 			if len(next) != 3 {
// 				return -1
// 			}
// 			sort.Slice(next, func(i, j int) bool {
// 				return next[i] > next[j]
// 			})
// 			return OnePair + int(k)<<16 + int(next[0])<<12 + int(next[1])<<8 + int(next[2])<<4
// 		}
// 	}
// 	return HighCard + h.cardsValue()
// }

func (h Hand) Value() int {
	// Count cards
	cards := make(map[byte]int)
	for i := 0; i < len(h); i++ {
		cards[h[i]]++
	}
	var (
		jokers int
		major  int
		vs     []int
	)
	// Render and sort duplicates
	for k, v := range cards {
		if k == 0x1 {
			jokers += v
		} else {
			vs = append(vs, v)
		}
	}
	sort.Slice(vs, func(i, j int) bool {
		return vs[i] > vs[j]
	})
	// Evaluate value
	if len(vs) == 0 {
		return FiveOfAKind + h.cardsValue()
	}
	switch vs[0] {
	case 5:
		major = FiveOfAKind
	case 4:
		switch jokers {
		case 1:
			major = FiveOfAKind
		default:
			major = FourOfAKind
		}
	case 3:
		switch jokers {
		case 2:
			major = FiveOfAKind
		case 1:
			major = FourOfAKind
		default:
			switch vs[1] {
			case 2:
				major = FullHouse
			default:
				major = ThreeOfAKind
			}
		}
	case 2:
		switch jokers {
		case 3:
			major = FiveOfAKind
		case 2:
			major = FourOfAKind
		case 1:
			switch vs[1] {
			case 2:
				major = FullHouse
			default:
				major = ThreeOfAKind
			}
		default:
			switch vs[1] {
			case 2:
				major = TwoPairs
			default:
				major = OnePair
			}
		}
	default:
		switch jokers {
		case 4:
			major = FiveOfAKind
		case 3:
			major = FourOfAKind
		case 2:
			major = ThreeOfAKind
		case 1:
			major = OnePair
		default:
			major = HighCard
		}
	}
	return major + h.cardsValue()
}
