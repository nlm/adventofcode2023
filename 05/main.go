package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"runtime"
	"strings"
	"sync"
	"time"

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
	maps := make(Maps)
	for _, part := range parts[1:] {
		k, cm, err := ProcessMap(part)
		if err != nil {
			return nil, nil, err
		}
		// fmt.Println("KEY", k, "MAP", cm)
		maps[*k] = cm
	}
	return seeds, maps, nil
}

type Range struct {
	Source int
	Dest   int
	Len    int
}
type ConversionMap struct {
	Ranges []*Range
}
type Maps map[Key]*ConversionMap
type Key struct {
	Source string
	Dest   string
}

func (r *Range) Contains(value int) bool {
	return value >= r.Source && value <= r.Source+r.Len
}

func (cm *ConversionMap) Value(key int) int {
	for _, r := range cm.Ranges {
		if r.Contains(key) {
			return r.Dest + (key - r.Source)
		}
	}
	return key
}

func (cm *ConversionMap) Values(key int) []int {
	var values []int
	for _, r := range cm.Ranges {
		if r.Contains(key) {
			values = append(values, r.Dest+(key-r.Source))
		}
	}
	if len(values) > 0 {
		return values
	}
	return []int{key}
}

func ProcessMap(data []byte) (*Key, *ConversionMap, error) {
	parts := strings.Split(string(data), "\n")
	headers := strings.Split(parts[0], " ")
	if len(headers) != 2 || headers[1] != "map:" {
		return nil, nil, fmt.Errorf("invalid header format: >%v<", string(parts[0]))
	}
	mapName := headers[0]
	mapKey, err := NameToKey(mapName)
	if err != nil {
		return nil, nil, err
	}
	// fmt.Println("MAPNAME", mapName)
	cm := &ConversionMap{Ranges: make([]*Range, 0, len(parts)-1)}
	for _, line := range parts[1:] {
		if len(line) == 0 {
			continue
		}
		corr := strings.Split(line, " ")
		if len(corr) != 3 {
			return nil, nil, fmt.Errorf("invalid correspondance line format: >%v<", line)
		}
		destRangeStart := utils.MustAtoi(corr[0])
		sourceRangeStart := utils.MustAtoi(corr[1])
		rangeLength := utils.MustAtoi(corr[2])
		// fmt.Println("RANGE", "SOURCE", sourceRangeStart, "DEST", destRangeStart, "LENGTH", rangeLength)
		cm.Ranges = append(cm.Ranges, &Range{
			Source: sourceRangeStart,
			Dest:   destRangeStart,
			Len:    rangeLength,
		})
	}
	return mapKey, cm, nil
}

func NameToKey(name string) (*Key, error) {
	parts := strings.Split(name, "-to-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid name: %v", name)
	}
	return &Key{
		Source: parts[0],
		Dest:   parts[1],
	}, nil
}

const (
	TEmpty       = ""
	TSeed        = "seed"
	TSoil        = "soil"
	TFertilizer  = "fertilizer"
	TWater       = "water"
	TLight       = "light"
	TTemperature = "temperature"
	THumidity    = "humidity"
	TLocation    = "location"
)

var NextTypeMap = map[string]string{
	TEmpty:       TSeed,
	TSeed:        TSoil,
	TSoil:        TFertilizer,
	TFertilizer:  TWater,
	TWater:       TLight,
	TLight:       TTemperature,
	TTemperature: THumidity,
	THumidity:    TLocation,
	TLocation:    TEmpty,
}

func (k *Key) Next() bool {
	var ok bool
	k.Source = k.Dest
	k.Dest, ok = NextTypeMap[k.Dest]
	if !ok {
		panic("invalid next type")
	}
	return k.Dest != ""
}

func InitialKey() Key {
	return Key{Source: "", Dest: "seed"}
}

// func (maps Maps) ResolveSeedLocation(seed int) int {
// 	key := InitialKey()
// 	value := seed
// 	for key.Next() {
// 		cm, ok := maps[key]
// 		if !ok {
// 			panic("error ResolveSeedLocation map lookup")
// 		}
// 		value = cm.Value(value)
// 	}
// 	return value
// }

func (maps Maps) ResolveSeedLocations(seeds []int) []int {
	key := InitialKey()
	for key.Next() {
		var newseeds []int
		cm, ok := maps[key]
		if !ok {
			panic("error ResolveSeedLocation map lookup")
		}
		for _, seed := range seeds {
			newseeds = append(newseeds, cm.Values(seed)...)
		}
		seeds = newseeds
	}
	return seeds
}

func Stage1(input io.Reader) (any, error) {
	seeds, maps, err := ParseInput(input)
	if err != nil {
		return nil, err
	}
	nearestLocation := MaxInt
	for _, seed := range seeds {
		locs := maps.ResolveSeedLocations([]int{seed})
		for _, loc := range locs {
			if loc < nearestLocation {
				nearestLocation = loc
			}
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

const (
	MaxUint = ^uint(0)
	MinUint = 0
	MaxInt  = int(MaxUint >> 1)
	MinInt  = -MaxInt - 1
)

func Stage2(input io.Reader) (any, error) {
	seeds, maps, err := ParseInput(input)
	if err != nil {
		return nil, err
	}
	seedRanges := SeedRanges(seeds)
	results := make(chan int, 10000)

	startTime := time.Now()
	fmt.Println("START", startTime)

	// Reduce
	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	nearestLocation := MaxInt
	go func(result <-chan int) {
		for r := range results {
			if r < nearestLocation {
				nearestLocation = r
			}
			fmt.Print(".")
		}
		wg2.Done()
	}(results)

	// Map
	limiter := make(chan struct{}, runtime.GOMAXPROCS(-1))
	fmt.Println("LIMIT", cap(limiter))
	wg := sync.WaitGroup{}
	for _, sr := range seedRanges {
		wg.Add(1)
		go func(sr SeedRange) {
			startTime := time.Now()
			nearestLocation := MaxInt
			for i := 0; i < sr.Length; i++ {
				seed := sr.Start + i
				locs := maps.ResolveSeedLocations([]int{seed})
				for _, loc := range locs {
					if loc < nearestLocation {
						nearestLocation = loc
					}
				}
			}
			results <- nearestLocation
			wg.Done()
			fmt.Println("done", time.Since(startTime))
		}(sr)
	}
	wg.Wait()
	close(results)
	wg2.Wait()

	fmt.Println(time.Since(startTime))
	return nearestLocation, nil
}

func (maps Maps) AsArray() []*ConversionMap {
	key := InitialKey()
	var mapsArray = make([]*ConversionMap, 0, len(maps))
	for key.Next() {
		mapsArray = append(mapsArray, maps[key])
	}
	return mapsArray
}

type PipeFunc func(in <-chan int, out chan<- int)

func (cm *ConversionMap) SendValues(key int, out chan<- int) {
	var found bool
	for _, r := range cm.Ranges {
		if r.Contains(key) {
			if !found {
				found = true
			}
			out <- r.Dest + (key - r.Source)
		}
	}
	if !found {
		out <- key
	}
}

// Stage3 is Stage2 with different run approach
func Stage3(input io.Reader) (any, error) {
	seeds, maps, err := ParseInput(input)
	if err != nil {
		return nil, err
	}
	seedRanges := SeedRanges(seeds)

	// Build Pipeline
	var pipeline []PipeFunc
	pipeline = append(pipeline, func(in <-chan int, out chan<- int) {
		for _, sr := range seedRanges {
			for i := 0; i < sr.Length; i++ {
				out <- sr.Start + i
			}
		}
		close(out)
	})
	for _, m := range maps.AsArray() {
		m := m
		pipeline = append(pipeline, func(in <-chan int, out chan<- int) {
			for i := range in {
				m.SendValues(i, out)
			}
			close(out)
		})
	}
	pipeline = append(pipeline, func(in <-chan int, out chan<- int) {
		var nearestLocation = MaxInt
		for location := range in {
			if location < nearestLocation {
				nearestLocation = location
			}
		}
		out <- nearestLocation
		close(out)
	})

	// Run the Pipeline
	chanSize := 1000
	wg := sync.WaitGroup{}
	in := make(chan int, chanSize)
	close(in)
	for _, fn := range pipeline {
		wg.Add(1)
		out := make(chan int, chanSize)
		go func(fn PipeFunc, in <-chan int, out chan<- int) {
			fn(in, out)
			wg.Done()
		}(fn, in, out)
		in = out
	}
	result := <-in
	wg.Wait()
	return result, nil
}

func main() {
	stage.RunCLI(input, Stage1, Stage2, Stage3)
}
