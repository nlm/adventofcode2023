package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/nlm/adventofcode2023/internal/utils"
)

type Condition struct {
	Symbol string
	Op     byte
	Value  int
}

func (c *Condition) Match(p Part) bool {
	if c == nil {
		return true
	}
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

func (c Condition) String() string {
	return fmt.Sprintf("%s%s%d", c.Symbol, string(c.Op), c.Value)
}

type Instruction struct {
	Condition *Condition
	Target    string
}

type Workflows map[string][]*Instruction

func (w Workflows) Visitor() *WorkflowsVisitor {
	return &WorkflowsVisitor{
		Workflows: w,
	}
}

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
var instrRE = regexp.MustCompile(`^(([\w+][><]\d+):)?(\w+)$`)
var condRE = regexp.MustCompile(`^([\w+])([><])(\d+)$`)

func ParseCondition(data []byte) (*Condition, error) {
	parts := condRE.FindSubmatch(data)
	if parts == nil {
		return nil, fmt.Errorf("ParseCondition: no match: %v", string(data))
	}
	return &Condition{
		Symbol: string(parts[1]),
		Op:     parts[2][0],
		Value:  utils.MustAtoi(string(parts[3])),
	}, nil
}

func ParseInstruction(data []byte) (*Instruction, error) {
	parts := instrRE.FindSubmatch(data)
	if parts == nil {
		return nil, fmt.Errorf("ParseInstruction: no match: %v", string(data))
	}
	if len(parts[1]) == 0 {
		return &Instruction{Target: string(parts[3])}, nil
	}
	cond, err := ParseCondition(parts[2])
	if err != nil {
		return nil, err
	}
	return &Instruction{
		Condition: cond,
		Target:    string(parts[3]),
	}, nil
}

func ParseInstructions(data []byte) ([]*Instruction, error) {
	parts := bytes.Split(data, []byte{','})
	if len(parts) == 0 {
		return nil, fmt.Errorf("ParseInstructions: no instructions")
	}
	instructions := make([]*Instruction, 0, len(parts))
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

func (w Workflows) CheckPart(part Part) (bool, error) {
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
		accepted, err := in.Workflows.CheckPart(part)
		if err != nil {
			return nil, err
		}
		if accepted {
			sum += part.Rating()
		}
	}
	return sum, nil
}

type WorkflowsVisitor struct {
	Workflows Workflows
	Ranges    []Ranges
}

func (we *WorkflowsVisitor) Visit(name string, ranges Ranges) {
	// fmt.Println("NA", name, "ST", stack)
	if name == "R" {
		// fmt.Println(stack, "->", "Rejected")
		return
	}
	if name == "A" {
		// fmt.Println(ranges, "->", "Accepted")
		we.Ranges = append(we.Ranges, ranges)
		return
	}
	wf := we.Workflows[name]
	for _, inst := range wf {
		if inst.Condition != nil {
			rMatch, rNoMatch := ranges.Split(*inst.Condition)
			we.Visit(inst.Target, rMatch)
			ranges = rNoMatch
		} else {
			we.Visit(inst.Target, ranges)
		}
	}
}

func (r Ranges) Split(c Condition) (Ranges, Ranges) {
	rMatch, rNoMatch := r, r
	offset := 0
	switch c.Symbol {
	case "x":
		offset = 0
	case "m":
		offset = 2
	case "a":
		offset = 4
	case "s":
		offset = 6
	}
	switch c.Op {
	// 1..4000
	// x < 2000
	// rMatch = 1..1999
	// rNoMatch = 2000...4000
	case '<':
		if rMatch[offset+1] > c.Value {
			rMatch[offset+1] = c.Value - 1
		}
		if rNoMatch[offset] < c.Value {
			rNoMatch[offset] = c.Value
		}
	// 1..4000
	// x > 2000
	// rMatch = 2001..4000
	// rNoMatch = 1..2000
	case '>':
		if rMatch[offset] < c.Value+1 {
			rMatch[offset] = c.Value + 1
		}
		if rNoMatch[offset+1] > c.Value {
			rNoMatch[offset+1] = c.Value
		}
	default:
		panic("unknown operator")
	}
	return rMatch, rNoMatch
}

type Ranges [8]int

func NewRanges() Ranges {
	return [8]int{1, 4000, 1, 4000, 1, 4000, 1, 4000}
}

func (r Ranges) Possibilities() int {
	product := 1
	for i := 0; i < 4; i++ {
		product *= r[i*2+1] - r[i*2] + 1
	}
	return product
}

func Stage2(input io.Reader) (any, error) {
	in, err := ParseInput(input)
	if err != nil {
		return nil, err
	}
	wv := in.Workflows.Visitor()
	wv.Visit("in", NewRanges())
	sum := 0
	for _, ranges := range wv.Ranges {
		sum += ranges.Possibilities()
	}
	return sum, nil
}
