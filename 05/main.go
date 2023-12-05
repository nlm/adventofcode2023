package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/nlm/adventofcode2023/internal/stage"
	"github.com/nlm/adventofcode2023/internal/utils"
)

//go:embed data/input.txt
var input []byte

func ParseInput(input io.Reader) ([]int, Maps, error) {
	data, err := io.ReadAll(input)
	if err != nil {
		return nil, nil, err
	}
	parts := bytes.Split(data, []byte("\n\n"))
	if len(parts) != 8 {
		return nil, nil, fmt.Errorf("invalid len: %d", len(parts))
	}
	seedsData := strings.Split(string(parts[0]), " ")
	if seedsData[0] != "seeds:" {
		return nil, nil, fmt.Errorf("invalid header: %v", seedsData[0])
	}
	var seeds []int
	for _, s := range seedsData[1:] {
		seeds = append(seeds, utils.MustAtoi(s))
	}
	fmt.Println("SEEDS", seeds)
	maps := make(Maps)
	for _, part := range parts[1:] {
		k, cm, err := ProcessMap(part)
		if err != nil {
			return nil, nil, err
		}
		fmt.Println("KEY", k, "MAP", cm)
		maps[k] = cm
	}
	return seeds, maps, nil
}

type Range struct {
	Source int
	Dest   int
	Len    int
}
type ConversionMap []Range
type Maps map[Key]ConversionMap
type Key struct {
	Source string
	Dest   string
}

func (r Range) Contains(value int) bool {
	return value >= r.Source && value <= r.Source+r.Len
}

func (cm ConversionMap) Value(key int) int {
	for _, r := range cm {
		// if key >= r.Source && key <= r.Source+r.Len {
		if r.Contains(key) {
			return r.Dest + (key - r.Source)
		}
	}
	return key
}

func ProcessMap(data []byte) (Key, ConversionMap, error) {
	parts := strings.Split(string(data), "\n")
	headers := strings.Split(parts[0], " ")
	if len(headers) != 2 || headers[1] != "map:" {
		return Key{}, nil, fmt.Errorf("invalid header format: >%v<", string(parts[0]))
	}
	mapName := headers[0]
	mapKey, err := NameToKey(headers[0])
	if err != nil {
		return Key{}, nil, err
	}
	fmt.Println("MAPNAME", mapName)
	cm := make(ConversionMap, 0, len(parts)-1)
	for _, line := range parts[1:] {
		if len(line) == 0 {
			continue
		}
		corr := strings.Split(line, " ")
		if len(corr) != 3 {
			return Key{}, nil, fmt.Errorf("invalid correspondance line format: >%v<", line)
		}
		destRangeStart := utils.MustAtoi(corr[0])
		sourceRangeStart := utils.MustAtoi(corr[1])
		rangeLength := utils.MustAtoi(corr[2])
		fmt.Println("RANGE", "SOURCE", sourceRangeStart, "DEST", destRangeStart, "LENGTH", rangeLength)
		cm = append(cm, Range{
			Source: sourceRangeStart,
			Dest:   destRangeStart,
			Len:    rangeLength,
		})
	}
	return mapKey, cm, nil
}

func NameToKey(name string) (Key, error) {
	parts := strings.Split(name, "-to-")
	if len(parts) != 2 {
		return Key{}, fmt.Errorf("invalid name: %v", name)
	}
	return Key{
		Source: parts[0],
		Dest:   parts[1],
	}, nil
}

func GetNextType(s string) string {
	switch s {
	case "":
		return "seed"
	case "seed":
		return "soil"
	case "soil":
		return "fertilizer"
	case "fertilizer":
		return "water"
	case "water":
		return "light"
	case "light":
		return "temperature"
	case "temperature":
		return "humidity"
	case "humidity":
		return "location"
	case "location":
		return ""
	default:
		panic("invalid NextType")
	}
}

func (k *Key) Next() bool {
	k.Source = k.Dest
	k.Dest = GetNextType(k.Dest)
	return k.Dest != ""
}

func InitialKey() Key {
	return Key{Source: "", Dest: "seed"}
}

func (maps Maps) ResolveSeedLocation(seed int) int {
	key := InitialKey()
	value := seed
	for key.Next() {
		cm, ok := maps[key]
		if !ok {
			panic("error ResolveSeedLocation map lookup")
		}
		// fmt.Print("-> ", key.Source, " ", value)
		value = cm.Value(value)
	}
	// fmt.Println("-> ", key.Source, " ", value)
	return value
}

func Stage1(input io.Reader) (any, error) {
	seeds, maps, err := ParseInput(input)
	if err != nil {
		return nil, err
	}
	nearestLocation := 99999999999999999
	for _, seed := range seeds {
		loc := maps.ResolveSeedLocation(seed)
		if loc < nearestLocation {
			nearestLocation = loc
		}
	}
	return nearestLocation, nil
}

type SeedRange struct {
	Start  int
	Length int
}

func SeedRanges(seeds []int) []SeedRange {
	var ranges []SeedRange
	for i := 0; i < len(seeds); i += 2 {
		ranges = append(ranges, SeedRange{Start: seeds[i], Length: seeds[i+1]})
	}
	return ranges
}

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

// WARNING: does not work
func Stage2(input io.Reader) (any, error) {
	seeds, maps, err := ParseInput(input)
	seedRanges := SeedRanges(seeds)
	results := make(chan int, len(seedRanges))
	wg := sync.WaitGroup{}
	fmt.Println("RANGE COUNT", len(seedRanges))
	for _, sr := range seedRanges {
		wg.Add(1)
		go func(sr SeedRange) {
			nearestLocation := MaxInt
			for i := 0; i < sr.Length; i++ {
				seed := sr.Start + i
				loc := maps.ResolveSeedLocation(seed)
				if loc < nearestLocation {
					nearestLocation = loc
				}
			}
			results <- nearestLocation
			wg.Done()
			fmt.Println("done")
		}(sr)
	}
	wg.Wait()
	close(results)
	nearestLocation := MaxInt
	for r := range results {
		if r < nearestLocation {
			nearestLocation = r
		}
	}
	return nearestLocation, err
}

func main() {
	stage.RunCLI(input, Stage1, Stage2)
}
