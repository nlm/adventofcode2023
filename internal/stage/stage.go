package stage

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var stages = make(map[int]StageFunc)
var flagStage = flag.Int("stage", 1, "stage to run")

var ErrUnimplemented = fmt.Errorf("unimplemented")

// StageFunc is the stage function signature.
type StageFunc func(input io.Reader) (any, error)

// Register registers a stage in the registry.
func Register(stage int, f StageFunc) {
	stages[stage] = f
}

// Run runs a registered stage function.
// It returns an error if the stage has not been registered.
func Run(stage int, input io.Reader) (any, error) {
	if f, ok := stages[stage]; ok {
		return f(input)
	}
	return nil, fmt.Errorf("stage %d not found", stage)
}

// RunCLI is a CLI helper to run stages.
func RunCLI(input any, fns ...StageFunc) {
	if !flag.Parsed() {
		flag.Parse()
	}
	for i, fn := range fns {
		Register(i+1, fn)
	}
	reader, err := Reader(input)
	if err != nil {
		log.Fatal(err)
	}
	res, err := Run(*flagStage, reader)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", res)
}

// TestCase represents the input and expected result of a test.
// Tests will individually run under the provided name.
// Input supports different types that will be converted to io.Reader.
// Result can be of any type, but it has to match the function result to succeed.
// Err is usually nil, but if we expect one, it can be given here.
type TestCase struct {
	// Name of the test
	Name string
	// Input can be []byte, string, or io.Reader
	Input any
	// Result can be of any type
	Result any
	// Error
	Err error
}

// Test runs the provided test cases against the StageFunc.
func Test(t *testing.T, fn StageFunc, cases []TestCase) {
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			reader, err := Reader(tc.Input)
			if !assert.NoError(t, err) {
				return
			}
			res, err := fn(reader)
			if tc.Err != nil && !assert.ErrorIs(t, err, tc.Err) {
				return
			}
			assert.Equal(t, tc.Result, res)
		})
	}
}

// Reader converts some classic input type as io.Reader.
func Reader(input any) (io.Reader, error) {
	if input == nil {
		return nil, fmt.Errorf("nil input")
	}
	if reader, ok := input.(io.Reader); ok {
		return reader, nil
	}
	switch in := input.(type) {
	case []byte:
		return bytes.NewReader(in), nil
	case string:
		return strings.NewReader(in), nil
	default:
		return nil, fmt.Errorf("unsupported input type: %t", input)
	}
}
