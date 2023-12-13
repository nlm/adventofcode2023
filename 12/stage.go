package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	_ "net/http/pprof"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/nlm/adventofcode2023/internal/utils"
)

type ConditionRecord struct {
	Record Record
	Stats  []int
}

func (cr ConditionRecord) Unfold() ConditionRecord {
	var stats []int
	for i := 0; i < 5; i++ {
		stats = append(stats, cr.Stats[:]...)
	}
	return ConditionRecord{
		Record: cr.Record.Unfold(),
		Stats:  stats,
	}
}

func ParseLine(line []byte) ConditionRecord {
	var cr ConditionRecord
	fields := bytes.Fields(line)
	if len(fields) != 2 {
		panic("field mismatch")
	}
	// fmt.Println("DATA", string(fields[0]), string(fields[1]))
	cr.Record = NewRecord(bytes.Clone(fields[0]))
	// fmt.Println(cr.Record)
	for _, v := range strings.Split(string(fields[1]), ",") {
		cr.Stats = append(cr.Stats, utils.MustAtoi(v))
	}
	return cr
}

func CalculateArrangement(line []byte) []int {
	var inHash = false
	var count = 0
	var counts = make([]int, 0, 16)
	for i := 0; i < len(line); i++ {
		switch line[i] {
		case '#':
			if !inHash {
				inHash = true
				count = 0
			}
			count++
		case '.':
			if inHash {
				inHash = false
				counts = append(counts, count)
			}
		default:
			fmt.Println("invalid:", line[i])
		}
	}
	if inHash {
		inHash = false
		counts = append(counts, count)
	}
	return counts
}

type BitField []bool

func (bf *BitField) Increment() bool {
	var carry bool
	var i int
	for i = len(*bf) - 1; i >= 0; i-- {
		carry = (*bf)[i]
		(*bf)[i] = !(*bf)[i]
		if !carry {
			break
		}
	}
	return !(i == -1 && carry)
}

type Record struct {
	Data   []byte
	Blocks [][]byte
}

func NewRecord(data []byte) Record {
	return Record{
		Data:   data,
		Blocks: bytes.FieldsFunc(data, func(r rune) bool { return r == '.' }),
	}
}

type Block struct {
	Data         []byte
	RenderedData []byte
	Indexes      []int
	Unknowns     int
}

func NewBlock(data []byte) Block {
	r := Block{
		Data:         bytes.Clone(data),
		RenderedData: bytes.Clone(data),
	}
	r.Unknowns = bytes.Count(r.Data, []byte{'?'})
	re := regexp.MustCompile(`\?`)
	for _, v := range re.FindAllIndex(r.Data, -1) {
		r.Indexes = append(r.Indexes, v[0])
	}
	return r
}

func (r Record) Unfold() Record {
	bytes := make([]byte, 0, 5*(len(r.Data)+1))
	for i := 0; i < 5; i++ {
		bytes = append(bytes, r.Data...)
		bytes = append(bytes, '?')
	}
	return NewRecord(bytes[:len(bytes)-1])
}

func (r *Block) Render(bf BitField) {
	if len(bf) != r.Unknowns {
		panic("length mismatch")
	}
	var b byte
	for i := 0; i < len(r.Indexes); i++ {
		if bf[i] {
			b = '#'
		} else {
			b = '.'
		}
		r.RenderedData[r.Indexes[i]] = b
	}
}

func SliceEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

var arrCache = make(map[string][][]int)
var arrCacheMut = sync.RWMutex{}

func Solve(block []byte) [][]int {
	key := string(block)
	arrCacheMut.RLock()
	if v, ok := arrCache[key]; ok {
		arrCacheMut.RUnlock()
		return v
	}
	arrCacheMut.RUnlock()
	var arrangements [][]int
	rec := NewBlock(block)
	if rec.Unknowns > 0 {
		bf := make(BitField, rec.Unknowns)
		for ok := true; ok; ok = bf.Increment() {
			rec.Render(bf)
			arrangements = append(arrangements, CalculateArrangement(rec.RenderedData))
		}
	} else {
		arrangements = append(arrangements, CalculateArrangement(block))
	}
	arrCacheMut.Lock()
	defer arrCacheMut.Unlock()
	arrCache[key] = arrangements
	return arrangements
}

func ParseInput(input io.Reader) ([]ConditionRecord, error) {
	var crs []ConditionRecord
	s := bufio.NewScanner(input)
	for s.Scan() {
		crs = append(crs, ParseLine(s.Bytes()))
	}
	if s.Err() != nil {
		return nil, s.Err()
	}
	return crs, nil
}

type Indexes struct {
	Values  [][][]int
	Current []int
	Max     []int
	values  []int
}

func NewIndexes(values [][][]int) Indexes {
	idxs := Indexes{
		Values:  values,
		Current: make([]int, len(values)),
		Max:     make([]int, len(values)),
		values:  make([]int, 0),
	}
	for i := 0; i < len(values); i++ {
		idxs.Max[i] = len(values[i]) - 1
	}
	return idxs
}

func (idxs *Indexes) Value() []int {
	// idxs.values = make([]int, 0, len(idxs.Values))
	idxs.values = idxs.values[:0]
	for i := 0; i < len(idxs.Values); i++ {
		if len(idxs.Values[i]) > 0 {
			idxs.values = append(idxs.values, idxs.Values[i][idxs.Current[i]]...)
		}
	}
	return idxs.values
}

func (idxs *Indexes) IsMax() bool {
	for i := 0; i < len(idxs.Values); i++ {
		if idxs.Current[i] != idxs.Max[i] {
			return false
		}
	}
	return true
}

func (idxs *Indexes) Incr() {
	for i := len(idxs.Values) - 1; i >= 0; i-- {
		var carry bool
		idxs.Current[i] += 1
		if idxs.Current[i] > idxs.Max[i] {
			idxs.Current[i] = 0
			carry = true
		}
		if !carry {
			break
		}
	}
}

func (cr *ConditionRecord) CountOptions() int {
	var options [][][]int
	// fmt.Println("data ", string(cr.Record.Data))
	for _, block := range cr.Record.Blocks {
		solved := Solve(block)
		// fmt.Println("solve", string(block), "->", solved)
		options = append(options, solved)
	}
	count := 0
	idxs := NewIndexes(options)
	for {
		// fmt.Println("CMP", idxs.Value(), cr.Stats)
		if SliceEqual(idxs.Value(), cr.Stats) {
			count++
		}
		if idxs.IsMax() {
			break
		}
		idxs.Incr()
	}
	return count
}

func Stage1(input io.Reader) (any, error) {
	crs, err := ParseInput(input)
	if err != nil {
		return nil, err
	}
	var count int
	for _, cr := range crs {
		optcount := cr.CountOptions()
		count += optcount
	}
	return count, nil
}

func Stage2(input io.Reader) (any, error) {
	crs, err := ParseInput(input)
	if err != nil {
		return nil, err
	}
	var (
		count int
		m     sync.Mutex
		wg    sync.WaitGroup
	)
	for _, cr := range crs {
		wg.Add(1)
		go func(cr ConditionRecord) {
			start := time.Now()
			cr = cr.Unfold()
			fmt.Println(string(cr.Record.Data), cr.Stats)
			optcount := cr.CountOptions()
			fmt.Println("DONE", time.Since(start))
			m.Lock()
			defer m.Unlock()
			count += optcount
			wg.Done()
		}(cr)
	}
	wg.Wait()
	return count, nil
}
