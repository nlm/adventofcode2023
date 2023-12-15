package main

import (
	"bytes"
	"fmt"
	"io"
	"regexp"

	"github.com/nlm/adventofcode2023/internal/utils"
)

func Report(b byte) {
	if b < 33 || b > 123 {
		fmt.Println("SUSPECT", b)
	}
}

func Hash(bytes []byte) int {
	currentValue := 0
	for _, b := range bytes {
		if b == '\n' || b == '\r' {
			continue
		}
		// Report(b)
		currentValue += int(b)
		currentValue = currentValue * 17
		currentValue %= 256
	}
	return currentValue
}

func Stage1(input io.Reader) (any, error) {
	// Result is > 186
	data := utils.Must(io.ReadAll(input))
	sum := 0
	for _, step := range bytes.Split(data, []byte{','}) {
		sum += Hash(step)
	}
	return sum, nil
}

type Lens struct {
	Label string
	Focal int
}

type Box struct {
	Lenses []Lens
}

func (b *Box) Has(label string) bool {
	return b.GetIndex(label) >= 0
}

func (b *Box) GetIndex(label string) int {
	for i := 0; i < len(b.Lenses); i++ {
		if b.Lenses[i].Label == label {
			return i
		}
	}
	return -1
}

func (b *Box) Remove(label string) bool {
	i := b.GetIndex(label)
	if i < 0 {
		return false
	}
	b.Lenses = append(b.Lenses[:i], b.Lenses[i+1:]...)
	return true
}

func (b *Box) Add(lens Lens) bool {
	i := b.GetIndex(lens.Label)
	if i < 0 {
		// Lens is not in the box
		b.Lenses = append(b.Lenses, lens)
		return false
	}
	b.Lenses[i] = lens
	return true
}

var (
	reLensEqual = regexp.MustCompile(`^(\w+)=(\d+)$`)
	reLensMinus = regexp.MustCompile(`^(\w+)-$`)
)

type Instruction struct {
	Remove bool
	Lens   Lens
}

func ParseInstruction(instruction []byte) *Instruction {
	instruction = bytes.ReplaceAll(instruction, []byte("\n"), []byte(""))
	res := reLensEqual.FindSubmatch(instruction)
	if res != nil {
		return &Instruction{
			Remove: false,
			Lens: Lens{
				Label: string(res[1]),
				Focal: utils.MustAtoi(string(res[2])),
			},
		}
	}
	res = reLensMinus.FindSubmatch(instruction)
	if res != nil {
		return &Instruction{
			Remove: true,
			Lens: Lens{
				Label: string(res[1]),
			},
		}
	}
	fmt.Println(">>", string(instruction), "<<")
	panic("unparseable instruction")
}

func Stage2(input io.Reader) (any, error) {
	data := utils.Must(io.ReadAll(input))
	boxes := make([]Box, 256)
	for _, step := range bytes.Split(data, []byte{','}) {
		inst := ParseInstruction(step)
		boxId := Hash([]byte(inst.Lens.Label))
		// fmt.Printf("%+v -> %v\n", inst, boxId)
		if inst.Remove {
			boxes[boxId].Remove(inst.Lens.Label)
		} else {
			boxes[boxId].Add(inst.Lens)
		}
		// for i, box := range boxes {
		// 	if len(box.Lenses) > 0 {
		// 		fmt.Println(i, box.Lenses)
		// 	}
		// }
	}
	result := 0
	for boxi, box := range boxes {
		for lensi, lens := range box.Lenses {
			// Focus Power
			result += (boxi + 1) * (lensi + 1) * lens.Focal
		}
	}
	return result, nil
}
