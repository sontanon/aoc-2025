package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"sync"
)

const (
	numWorkers = 8
)

type Lights []bool

func (l Lights) Equals(r Lights) bool {
	if len(l) != len(r) {
		return false
	}
	for i := range l {
		if l[i] != r[i] {
			return false
		}
	}
	return true
}

func (l Lights) String() string {
	buf := bytes.Buffer{}
	for _, light := range l {
		if light {
			buf.WriteByte('#')
		} else {
			buf.WriteByte('.')
		}
	}
	return buf.String()
}

type Button func(Lights) (Lights, error)

type ActionSpace struct {
	Goal    Lights
	Buttons []Button
}

func Solve(actionSpace ActionSpace) (int, error) {
	start := make(Lights, len(actionSpace.Goal))
	visited := make(map[string]bool)
	type State struct {
		Lights Lights
		Steps  int
	}
	queue := []State{{Lights: start, Steps: 0}}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if current.Lights.Equals(actionSpace.Goal) {
			return current.Steps, nil
		}
		key := current.Lights.String()
		if visited[key] {
			continue
		}
		visited[key] = true
		for _, button := range actionSpace.Buttons {
			newLights, err := button(current.Lights)
			if err != nil {
				return 0, err
			}
			queue = append(queue, State{Lights: newLights, Steps: current.Steps + 1})
		}
	}
	return 0, errors.New("no solution found")
}

func ParseLights(input []byte) (Lights, error) {
	if len(input) < 3 {
		return nil, errors.New("input is too short to be valid")
	}
	if input[0] != '[' && input[len(input)-1] != ']' {
		return nil, errors.New("input is not enclosed in square brackets")
	}
	output := make([]bool, len(input)-2)
	for i, b := range input[1 : len(input)-1] {
		if b != '.' && b != '#' {
			return nil, fmt.Errorf("unrecognized character '%v' at position %d", b, i)
		}
		if b == '#' {
			output[i] = true
		}
	}
	return output, nil
}

func ParseButton(input []byte) (Button, error) {
	if len(input) < 3 {
		return nil, errors.New("input is too short to be valid")
	}
	if input[0] != '(' && input[len(input)-1] != ')' {
		return nil, errors.New("input is not enclosed in parentheses")
	}
	subSlice := bytes.Split(input[1:len(input)-1], []byte{','})
	indices := make([]int, 0, len(subSlice))
	for i, b := range subSlice {
		idx, err := strconv.Atoi(string(b))
		if err != nil {
			return nil, fmt.Errorf("error processing field '%v' (%d): %w", b, i, err)
		}
		indices = append(indices, idx)
	}
	button := func(input Lights) (Lights, error) {
		output := slices.Clone(input)
		for _, idx := range indices {
			if idx > len(input) {
				return nil, fmt.Errorf("requested toggle index %d is outside the lights array length of %d", idx, len(input))
			}
			output[idx] = !input[idx]
		}
		return output, nil
	}
	return button, nil
}

func ParseLine(input []byte) (ActionSpace, error) {
	fields := bytes.Fields(input)
	if len(fields) < 3 {
		return ActionSpace{}, errors.New("invalid input does not have at least 3 whitespace separated fields")
	}
	goal, err := ParseLights(fields[0])
	if err != nil {
		return ActionSpace{}, fmt.Errorf("failed to build goal: %w", err)
	}

	buttons := make([]Button, 0, len(fields)-2)
	for i, f := range fields[1 : len(fields)-1] {
		button, err := ParseButton(f)
		if err != nil {
			return ActionSpace{}, fmt.Errorf("failed to build button %d: %w", i+1, err)
		}
		buttons = append(buttons, button)
	}
	return ActionSpace{Goal: goal, Buttons: buttons}, nil
}

func ParseInput(input io.Reader) ([]ActionSpace, error) {
	lineScanner := bufio.NewScanner(input)
	result := make([]ActionSpace, 0)
	numLine := 0
	for lineScanner.Scan() {
		numLine++
		line := lineScanner.Bytes()
		actionSpace, err := ParseLine(line)
		if err != nil {
			return nil, fmt.Errorf("line %d is invalid: %w", numLine, err)
		}
		result = append(result, actionSpace)
	}
	if err := lineScanner.Err(); err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.New("input does not contain any action spaces")
	}
	return result, nil
}

func Part1(input io.Reader) (int, error) {
	actionSpaces, err := ParseInput(input)
	if err != nil {
		return 0, err
	}
	solutions := make([]int, len(actionSpaces))

	chunkPerWorker := (len(actionSpaces) + numWorkers - 1) / numWorkers
	wg := sync.WaitGroup{}
	// log.Printf("Solving %d action spaces with %d workers (chunk size %d)", len(actionSpaces), numWorkers, chunkPerWorker)
	for w := range numWorkers {
		wg.Go(
			func() {
				startIdx := w * chunkPerWorker
				endIdx := min((w+1)*chunkPerWorker, len(actionSpaces))
				for i, actionSpace := range actionSpaces[startIdx:endIdx] {
					solution, err := Solve(actionSpace)
					if err != nil {
						panic(fmt.Sprintf("failed to solve action space %d: %v", startIdx+i+1, err))
					}
					solutions[startIdx+i] = solution
				}
			})
	}
	wg.Wait()

	total := 0
	for _, s := range solutions {
		total += s
	}
	return total, nil
}

func getInputPath() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	return filepath.Join(dir, "input.txt")
}

func main() {
	file, err := os.Open(getInputPath())
	if err != nil {
		panic(err)
	}
	defer file.Close()

	result1, err := Part1(file)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 1:", result1)

}
