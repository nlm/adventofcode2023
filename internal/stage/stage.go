package stage

import (
	"flag"
	"fmt"
	"log"
)

var stages = make(map[int]func() error)
var flagStage = flag.Int("stage", 1, "stage to run")

func Register(stage int, f func() error) {
	stages[stage] = f
}

func Run(stage int) error {
	if f, ok := stages[stage]; ok {
		return f()
	}
	return fmt.Errorf("stage %d not found", stage)
}

func RunCLI(fns ...func() error) {
	if !flag.Parsed() {
		flag.Parse()
	}
	for i, fn := range fns {
		Register(i+1, fn)
	}
	err := Run(*flagStage)
	if err != nil {
		log.Fatal(err)
	}
}
