package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/nlm/adventofcode2023/internal/stage"
)

//go:embed data/input.txt
var input []byte

var NodeRE = regexp.MustCompile(`^(\w{3}) = \((\w{3}), (\w{3})\)$`)

type Node struct {
	Name  string
	Left  string
	Right string
}

type Input struct {
	Instructions []byte
	Nodes        map[string]*Node
}

func ParseInput(input io.Reader) (*Input, error) {
	s := bufio.NewScanner(input)
	s.Scan()
	data := Input{}
	data.Instructions = bytes.Clone(s.Bytes())
	data.Nodes = make(map[string]*Node)
	s.Scan()
	if len(s.Bytes()) != 0 {
		return nil, fmt.Errorf("no empty line")
	}
	for s.Scan() {
		line := s.Bytes()
		matches := NodeRE.FindSubmatch(line)
		if len(matches) != 4 {
			return nil, fmt.Errorf("parse error: %s", string(line))
		}
		// if _, ok := data.Nodes[string(fields[0])]; ok {
		// 	panic("duplicate")
		// }
		data.Nodes[string(matches[1])] = &Node{
			Name:  string(matches[1]),
			Left:  string(matches[2]),
			Right: string(matches[3]),
		}
	}
	return &data, nil
}

func Stage1(input io.Reader) (any, error) {
	data, err := ParseInput(input)
	if err != nil {
		return nil, err
	}
	return data.CountSteps("AAA", func(n *Node) bool { return n.Name == "ZZZ" }), nil
}

func (data *Input) CountSteps(start string, stopFunc func(*Node) bool) int {
	counter := 0
	for node := data.Nodes[start]; !stopFunc(node); {
		instruction := data.Instructions[counter%len(data.Instructions)]
		node = data.NextNode(node, instruction)
		counter++
	}
	return counter
}

func (data *Input) NextNode(node *Node, instruction byte) *Node {
	switch instruction {
	case 'L':
		return data.Nodes[node.Left]
	case 'R':
		return data.Nodes[node.Right]
	default:
		panic("invalid instruction")
	}
}

func Stage2(input io.Reader) (any, error) {
	data, err := ParseInput(input)
	if err != nil {
		return nil, err
	}
	var nodes []*Node
	for k := range data.Nodes {
		if strings.HasSuffix(k, "A") {
			nodes = append(nodes, data.Nodes[k])
		}
	}
	results := make([]int, len(nodes))
	for i := 0; i < len(nodes); i++ {
		results[i] = data.CountSteps(nodes[i].Name, func(node *Node) bool {
			return strings.HasSuffix(node.Name, "Z")
		})
	}
	return LCM(results[0], results[1], results[2:]...), nil
}

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)
	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}
	return result
}

func main() {
	stage.RunCLI(input, Stage1, Stage2)
}
