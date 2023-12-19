package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/nlm/adventofcode2023/internal/stage"
	"github.com/nlm/adventofcode2023/internal/utils"
)

type Condition struct {
	Symbol string
	Op     byte
	Value  int
}

func (c Condition) Match(p Part) bool {
	value := p[c.Symbol]
	switch c.Op {
	case '<':
		return value < c.Value
	case '>':
		return value > c.Value
	default:
		panic("unknown operator")
	}
}

type Instruction struct {
	Condition *Condition
	Target    string
}

type Workflows map[string][]Instruction

// func (w Workflows) String() string {
// 	return fmt.Sprintf("%")
// }

func (i Instruction) String() string {
	sb := strings.Builder{}
	if i.Condition != nil {
		fmt.Fprintf(&sb, "%s%s%d:%s", i.Condition.Symbol, string(i.Condition.Op), i.Condition.Value, i.Target)
	} else {
		fmt.Fprintf(&sb, "%s", i.Target)
	}
	return sb.String()
}

type Part map[string]int

func (p Part) Rating() int {
	var sum = 0
	for _, v := range p {
		sum += v
	}
	return sum
}

type Input struct {
	Workflows Workflows
	Parts     []Part
}

func (p Part) String() string {
	sb := strings.Builder{}
	fmt.Fprint(&sb, "{")
	for k, v := range p {
		fmt.Fprintf(&sb, "%s=%d,", k, v)
	}
	fmt.Fprint(&sb, "}")
	return sb.String()
}

var wflowRE = regexp.MustCompile(`^([a-z]+){(.+)}$`)
var instrRE = regexp.MustCompile(`^((\w+)([><])(\d+):)?(\w+)$`)

func ParseInstruction(data []byte) (Instruction, error) {
	parts := instrRE.FindSubmatch(data)
	if parts == nil {
		return Instruction{}, fmt.Errorf("ParseInstruction: no match: %v", string(data))
	}
	if len(parts[1]) == 0 {
		return Instruction{Target: string(parts[5])}, nil
	}
	return Instruction{
		Condition: &Condition{
			Symbol: string(parts[2]),
			Op:     parts[3][0],
			Value:  utils.MustAtoi(string(parts[4])),
		},
		Target: string(parts[5]),
	}, nil
}

func ParseInstructions(data []byte) ([]Instruction, error) {
	parts := bytes.Split(data, []byte{','})
	if len(parts) == 0 {
		return nil, fmt.Errorf("ParseInstructions: no instructions")
	}
	instructions := make([]Instruction, 0, len(parts))
	for _, part := range parts {
		instruction, err := ParseInstruction(part)
		if err != nil {
			return nil, err
		}
		instructions = append(instructions, instruction)
	}
	return instructions, nil
}

func ParseWorkflows(s *bufio.Scanner) (Workflows, error) {
	workflows := make(Workflows)
	for s.Scan() && len(s.Bytes()) > 0 {
		matches := wflowRE.FindSubmatch(s.Bytes())
		if matches == nil {
			return nil, fmt.Errorf("no match: %v", s.Text())
		}
		instructions, err := ParseInstructions(matches[2])
		if err != nil {
			return nil, err
		}
		workflows[string(matches[1])] = instructions
	}
	return workflows, nil
}

var attrRe = regexp.MustCompile(`^(\w+)=(\d+)$`)
var partRe = regexp.MustCompile(`^\{(.+)\}$`)

func ParsePart(data []byte) (Part, error) {
	matches := partRe.FindSubmatch(data)
	if matches == nil {
		return nil, fmt.Errorf("ParsePart: no match: %v", string(data))
	}
	attrs := bytes.Split(matches[1], []byte{','})
	if len(attrs) == 0 {
		return nil, fmt.Errorf("ParsePart: no attributes: %v", string(matches[1]))
	}
	part := make(Part)
	for _, attr := range attrs {
		amatches := attrRe.FindSubmatch(attr)
		if amatches == nil {
			return nil, fmt.Errorf("ParsePart: no match attribute: %v", string(attr))
		}
		part[string(amatches[1])] = utils.MustAtoi(string(amatches[2]))

	}
	return part, nil
}

func ParseParts(s *bufio.Scanner) ([]Part, error) {
	parts := make([]Part, 0)
	for s.Scan() {
		part, err := ParsePart(s.Bytes())
		if err != nil {
			return nil, err
		}
		parts = append(parts, part)
	}
	if s.Err() != nil {
		return nil, s.Err()
	}
	return parts, nil
}

func ParseInput(input io.Reader) (*Input, error) {
	s := bufio.NewScanner(input)
	workflows, err := ParseWorkflows(s)
	if err != nil {
		return nil, err
	}
	if s.Err() != nil {
		return nil, err
	}
	parts, err := ParseParts(s)
	if err != nil {
		return nil, err
	}
	if s.Err() != nil {
		return nil, err
	}
	return &Input{
		Workflows: workflows,
		Parts:     parts,
	}, nil
}

func (w Workflows) RunPart(part Part) (bool, error) {
	wfName := "in"
	for wfName != "R" && wfName != "A" {
		for _, inst := range w[wfName] {
			condition := inst.Condition
			if condition == nil {
				wfName = inst.Target
				break
			}
			if condition.Match(part) {
				wfName = inst.Target
				break
			}
		}
		// return false, fmt.Errorf("RunPart: no matching instruction\n  %s\n  %s", w, part)
	}
	return wfName == "A", nil
}

func Stage1(input io.Reader) (any, error) {
	in, err := ParseInput(input)
	if err != nil {
		return nil, err
	}
	sum := 0
	for _, part := range in.Parts {
		accepted, err := in.Workflows.RunPart(part)
		if err != nil {
			return nil, err
		}
		if accepted {
			sum += part.Rating()
		}
	}
	return sum, nil
}

// func (w Workflows) Explore(name string) [][]*Condition {
// 	conditionsList := make([][]*Condition, 0)
// 	for _, instruction := range w[name] {
// 		conditions := make([]*Condition, 0)
// 		if instruction.Target == "R" {
// 			continue
// 		}
// 		if instruction.Target == "A" {
// 			conditions = append(conditions, nil)
// 			continue
// 		}
// 		conditions = append(conditions, instruction.Condition)
// 	}
// }

func Stage2(input io.Reader) (any, error) {
	// in, err := ParseInput(input)
	// if err != nil {
	// 	return nil, err
	// }
	// in.Workflows.Explore(in)
	return nil, stage.ErrUnimplemented
}
