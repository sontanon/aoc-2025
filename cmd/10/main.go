package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"sync"
)

const (
	numWorkers    = 8
	maxIterations = 1_000_000_000
)

type Lights []bool

type State struct {
	Lights  Lights
	Joltage []int
}

func (s State) Equals(r State, lightEquality bool) bool {
	if lightEquality {
		if len(s.Lights) != len(r.Lights) {
			return false
		}
		for i := range s.Lights {
			if s.Lights[i] != r.Lights[i] {
				return false
			}
		}
	} else {
		if len(s.Joltage) != len(r.Joltage) {
			return false
		}
		for i := range s.Joltage {
			if s.Joltage[i] != r.Joltage[i] {
				return false
			}
		}
	}
	return true
}

func (s State) String() string {
	buf := bytes.Buffer{}
	// buf.Grow(4 * len(s.Lights))
	buf.WriteByte('[')
	for _, light := range s.Lights {
		if light {
			buf.WriteByte('#')
		} else {
			buf.WriteByte('.')
		}
	}
	buf.WriteByte(']')
	if len(s.Joltage) > 0 {
		buf.WriteByte(' ')
		buf.WriteByte('{')
		for i, j := range s.Joltage {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(strconv.Itoa(j))
		}
		buf.WriteByte('}')
	}
	return buf.String()
}

type Button func(initial State, maximumJoltage []int) (State, bool, error)

type ActionSpace struct {
	Goal    State
	Buttons []Button
}

func Solve(actionSpace ActionSpace, ignoreJoltage bool) (int, error) {
	start := State{
		Lights: make(Lights, len(actionSpace.Goal.Lights)),
	}
	if !ignoreJoltage {
		start.Joltage = make([]int, len(actionSpace.Goal.Joltage))
	}
	visited := make(map[string]bool)
	type Tracker struct {
		State State
		Steps int
	}
	queue := []Tracker{{State: start, Steps: 0}}
	iterations := 0
	for ; iterations < maxIterations && len(queue) > 0; iterations++ {
		current := queue[0]
		queue = queue[1:]
		if current.State.Equals(actionSpace.Goal, ignoreJoltage) {
			return current.Steps, nil
		}
		key := current.State.String()
		if visited[key] {
			continue
		}
		visited[key] = true
		for _, button := range actionSpace.Buttons {
			var maxJoltage []int
			if !ignoreJoltage {
				maxJoltage = actionSpace.Goal.Joltage
			}
			newState, valid, err := button(current.State, maxJoltage)
			if err != nil {
				return 0, err
			}
			if !valid {
				continue
			}
			queue = append(queue, Tracker{State: newState, Steps: current.Steps + 1})
		}
	}
	return 0, fmt.Errorf("no solution found after %d iterations, visited %d states", iterations, len(visited))
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
	button := func(initial State, desiredJoltage []int) (State, bool, error) {
		outputLights := slices.Clone(initial.Lights)
		outputJoltage := slices.Clone(initial.Joltage)

		for _, idx := range indices {
			if idx > len(initial.Lights) {
				return State{}, false, fmt.Errorf("requested toggle index %d is outside the lights array length of %d", idx, len(initial.Lights))
			}
			outputLights[idx] = !initial.Lights[idx]

			if len(desiredJoltage) > 0 {
				if idx > len(initial.Joltage) {
					return State{}, false, fmt.Errorf("requested toggle index %d is outside the joltage array length of %d", idx, len(initial.Joltage))
				}
				outputJoltage[idx] = initial.Joltage[idx] + 1
				if outputJoltage[idx] > desiredJoltage[idx] {
					return State{}, false, nil
				}
			}
		}
		return State{Lights: outputLights, Joltage: outputJoltage}, true, nil
	}
	return button, nil
}

func ParseJoltage(input []byte) ([]int, error) {
	if len(input) < 3 {
		return nil, errors.New("input is too short to be valid")
	}
	if input[0] != '{' && input[len(input)-1] != '}' {
		return nil, errors.New("input is not enclosed in curly braces")
	}
	subSlice := bytes.Split(input[1:len(input)-1], []byte{','})
	output := make([]int, 0, len(subSlice))
	for i, b := range subSlice {
		joltage, err := strconv.Atoi(string(b))
		if err != nil {
			return nil, fmt.Errorf("error processing field '%v' (%d): %w", b, i, err)
		}
		output = append(output, joltage)
	}
	return output, nil
}

func ParseLine(input []byte) (ActionSpace, error) {
	fields := bytes.Fields(input)
	if len(fields) < 3 {
		return ActionSpace{}, errors.New("invalid input does not have at least 3 whitespace separated fields")
	}
	lightsGoal, err := ParseLights(fields[0])
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
	joltage, err := ParseJoltage(fields[len(fields)-1])
	if err != nil {
		return ActionSpace{}, fmt.Errorf("failed to build maximum joltage: %w", err)
	}

	return ActionSpace{Goal: State{Lights: lightsGoal, Joltage: joltage}, Buttons: buttons}, nil
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

func MultiSolve(input io.Reader, ignoreJoltage bool) (int, error) {
	actionSpaces, err := ParseInput(input)
	if err != nil {
		return 0, err
	}
	solutions := make([]int, len(actionSpaces))

	wg := sync.WaitGroup{}
	chunkPerWorker := (len(actionSpaces) + numWorkers - 1) / numWorkers
	log.Printf("solving %d action spaces with %d workers, %d chunks each", len(actionSpaces), numWorkers, chunkPerWorker)
	for w := range numWorkers {
		startIdx := w * chunkPerWorker
		endIdx := min((w+1)*chunkPerWorker, len(actionSpaces))
		wg.Go(
			func() {
				for i := startIdx; i < endIdx; i++ {
					solution, err := Solve(actionSpaces[i], ignoreJoltage)
					if err != nil {
						log.Printf("failed to solve action space %d: %v", i, err)
						continue
					}
					log.Printf("solved action space %d: %d", i, solution)
					solutions[i] = solution
				}
			},
		)
	}
	wg.Wait()

	// for i, actionSpace := range actionSpaces {
	// 	solution, err := Solve(actionSpace, ignoreJoltage)
	// 	if err != nil {
	// 		return 0, fmt.Errorf("failed to solve action space %d: %w", i, err)
	// 	}
	// 	solutions[i] = solution
	// }

	total := 0
	for _, s := range solutions {
		total += s
	}
	return total, nil
}

func VerifySolution(actionSpace ActionSpace, buttonSequence []int, ignoreJoltage bool) (bool, error) {
	currentState := State{
		Lights: make(Lights, len(actionSpace.Goal.Lights)),
	}
	if !ignoreJoltage {
		currentState.Joltage = make([]int, len(actionSpace.Goal.Joltage))
	}
	for step, buttonIdx := range buttonSequence {
		if buttonIdx < 0 || buttonIdx >= len(actionSpace.Buttons) {
			return false, fmt.Errorf("button index %d at step %d is out of range", buttonIdx, step)
		}
		button := actionSpace.Buttons[buttonIdx]
		newState, valid, err := button(currentState, actionSpace.Goal.Joltage)
		if err != nil {
			return false, fmt.Errorf("error applying button %d at step %d: %w", buttonIdx, step, err)
		}
		if !valid {
			return false, fmt.Errorf("button %d at step %d produced an invalid state", buttonIdx, step)
		}
		currentState = newState
	}
	return currentState.Equals(actionSpace.Goal, ignoreJoltage), nil
}

func Part1(input io.Reader) (int, error) {
	return MultiSolve(input, true)
}

func Part2(input io.Reader) (int, error) {
	return MultiSolve(input, false)
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

	file.Seek(0, io.SeekStart)
	result2, err := Part2(file)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 2:", result2)

}
