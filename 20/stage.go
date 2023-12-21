package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/nlm/adventofcode2023/internal/stage"
)

type ModuleType int

const (
	Receiver ModuleType = iota
	Broadcaster
	FlipFlop
	Conjunction
)

type Module struct {
	Name          string
	Targets       []string
	Type          ModuleType
	High          bool
	LastPulse     map[string]bool
	LowPulseCount int
}

type System struct {
	Modules        map[string]*Module
	Pulses         chan Pulse
	LowPulseCount  int
	HighPulseCount int
}

var lineRE = regexp.MustCompile(`^([&%]?)(\w+) -> (([a-z](, )?)+)$`)

func ParseLine(data []byte) (*Module, error) {
	matches := lineRE.FindSubmatch(data)
	if matches == nil {
		return nil, fmt.Errorf("ParseInput: no match: %v", string(data))
	}
	moduleType := Receiver
	var lastPulse map[string]bool
	switch string(matches[1]) {
	case "":
		moduleType = Broadcaster
	case "&":
		moduleType = Conjunction
		lastPulse = make(map[string]bool)
	case "%":
		moduleType = FlipFlop
	}
	return &Module{
		Name:      string(matches[2]),
		Targets:   strings.Split(strings.ReplaceAll(string(matches[3]), " ", ""), ","),
		Type:      moduleType,
		LastPulse: lastPulse,
	}, nil
}

func ParseInput(input io.Reader) (*System, error) {
	system := &System{
		Modules: map[string]*Module{},
		Pulses:  make(chan Pulse, 100),
	}
	s := bufio.NewScanner(input)
	for s.Scan() {
		module, err := ParseLine(s.Bytes())
		if err != nil {
			return nil, err
		}
		system.Modules[module.Name] = module
		// create missing receive-only targets
		for _, target := range module.Targets {
			if _, ok := system.Modules[target]; !ok {
				system.Modules[target] = &Module{
					Name: target,
					Type: Receiver,
				}
			}
		}
	}
	if s.Err() != nil {
		return nil, s.Err()
	}
	for _, module := range system.Modules {
		for _, target := range module.Targets {
			if system.Modules[target].Type == Conjunction {
				system.Modules[target].LastPulse[module.Name] = false
			}
		}
	}
	return system, nil
}

type Pulse struct {
	Sender string
	Target string
	High   bool
}

func (pulse Pulse) String() string {
	return fmt.Sprintf(
		"%s -%s-> %s",
		pulse.Sender,
		map[bool]string{false: "low", true: "high"}[pulse.High],
		pulse.Target,
	)
}

// type Pulser struct {
// 	system *System
// 	pulses chan Pulse
// }

// func NewPulser(s *System) *Pulser {
// 	return &Pulser{
// 		system: s,
// 		pulses: make(chan Pulse, 100),
// 	}
// }

func (s *System) SendPulse(pulse Pulse) {
	if pulse.High {
		s.HighPulseCount++
	} else {
		s.LowPulseCount++
	}
	s.Pulses <- pulse
}

func (s *System) PushButton() {
	s.SendPulse(Pulse{
		Sender: "button",
		Target: "broadcaster",
		High:   false,
	})
}

func (s *System) PropagatePulses() bool {
	var pulse Pulse
	// read a pulse from the queue
	select {
	case pulse = <-s.Pulses:
		break
	default:
		return false
	}
	if stage.Verbose() {
		fmt.Println("PULSE", pulse)
	}
	// lookup module
	mod := s.Modules[pulse.Target]
	if mod == nil {
		fmt.Println("unknown module:", pulse.Target)
		return false
	}
	// get new state after pulse + shall we emit
	high, emit := mod.ReceivePulse(pulse)
	if !emit {
		return true
	}
	for _, target := range mod.Targets {
		// if stage.Verbose() {
		// 	fmt.Println("PULSE", pulse.Target, "-"+map[bool]string{false: "low", true: "high"}[high]+"->", target)
		// }
		s.SendPulse(Pulse{
			Sender: mod.Name,
			Target: target,
			High:   high,
		})
	}
	return true
}

// Pulse -> state, emit bool
func (m *Module) ReceivePulse(pulse Pulse) (bool, bool) {
	if !pulse.High {
		m.LowPulseCount++
	}
	switch m.Type {
	case Receiver:
		// panic("can't send pulse to receiver")
		return false, false
	case Broadcaster:
		return pulse.High, true
	case FlipFlop:
		switch pulse.High {
		case true:
			return false, false
		case false:
			m.High = !m.High
			return m.High, true
		}
	case Conjunction:
		// Conjunction modules (prefix &) remember the type
		// of the most recent pulse received from each of
		// their connected input modules; they initially
		// default to remembering a low pulse for each input.
		// When a pulse is received, the conjunction module first
		// updates its memory for that input. Then, if it
		// remembers high pulses for all inputs, it sends a
		// low pulse; otherwise, it sends a high pulse.
		m.LastPulse[pulse.Sender] = pulse.High
		// FIXME
		for _, v := range m.LastPulse {
			if !v {
				return true, true
			}
		}
		return false, true
	}
	panic("unknown type")
}

func Stage1(input io.Reader) (any, error) {
	system, err := ParseInput(input)
	if err != nil {
		return nil, err
	}
	// for _, mod := range system.Modules {
	// 	fmt.Printf("%+v\n", *mod)
	// }
	for i := 0; i < 1000; i++ {
		for system.PushButton(); system.PropagatePulses(); {
		}
	}
	return system.LowPulseCount * system.HighPulseCount, nil
}

func Stage2(input io.Reader) (any, error) {
	system, err := ParseInput(input)
	if err != nil {
		return nil, err
	}
	for i := 0; ; i++ {
		for system.PushButton(); system.PropagatePulses(); {
		}
		if system.Modules["rx"].LowPulseCount > 0 {
			return i + 1, nil
		}
	}
}

func Stage3(input io.Reader) (any, error) {
	system, err := ParseInput(input)
	if err != nil {
		return nil, err
	}
	fmt.Println("digraph G {")
	fmt.Println("  overlap=false;")
	for _, module := range system.Modules {
		var color string
		switch module.Type {
		case Receiver:
			color = "orange"
		case Broadcaster:
			color = "green"
		case FlipFlop:
			color = "red"
		case Conjunction:
			color = "blue"
		}
		fmt.Printf("  %s [color=%s];\n", module.Name, color)
		for _, target := range module.Targets {
			fmt.Printf("  %s -> %s;\n", module.Name, target)
		}
	}
	fmt.Println("}")
	return nil, nil
}
