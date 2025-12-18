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

type Button struct {
	Indices []int
}

func (b Button) Apply(initial State) (State, error) {
	outputLights := slices.Clone(initial.Lights)
	outputJoltage := slices.Clone(initial.Joltage)

	for _, idx := range b.Indices {
		if idx >= len(initial.Lights) || idx >= len(initial.Joltage) {
			return State{}, fmt.Errorf("requested toggle index %d is outside the lights array length of %d", idx, len(initial.Lights))
		}
		outputLights[idx] = !initial.Lights[idx]
		outputJoltage[idx] = initial.Joltage[idx] + 1
	}
	return State{Lights: outputLights, Joltage: outputJoltage}, nil
}

func (b Button) WithinJoltageLimits(state State, maxJoltage []int) (bool, error) {
	for _, idx := range b.Indices {
		if idx >= len(state.Joltage) || idx >= len(maxJoltage) {
			return false, fmt.Errorf("requested toggle index %d is outside the joltage array length of %d", idx, len(state.Joltage))
		}
		if state.Joltage[idx]+1 > maxJoltage[idx] {
			return false, nil
		}
	}
	return true, nil
}

type ActionSpace struct {
	Goal    State
	Buttons []Button
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
	if len(joltage) != len(lightsGoal) {
		return ActionSpace{}, errors.New("joltage array length does not match lights array length")
	}

	return ActionSpace{Goal: State{Lights: lightsGoal, Joltage: joltage}, Buttons: buttons}, nil
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
		return Button{}, errors.New("input is too short to be valid")
	}
	if input[0] != '(' && input[len(input)-1] != ')' {
		return Button{}, errors.New("input is not enclosed in parentheses")
	}
	subSlice := bytes.Split(input[1:len(input)-1], []byte{','})
	indices := make([]int, 0, len(subSlice))
	for i, b := range subSlice {
		idx, err := strconv.Atoi(string(b))
		if err != nil {
			return Button{}, fmt.Errorf("error processing field '%v' (%d): %w", b, i, err)
		}
		indices = append(indices, idx)
	}
	button := Button{Indices: indices}
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

func MultiSolve(
	actionSpaces []ActionSpace,
	maxIterations int,
	equalityFunc func(a, b State) bool,
	hashFunc func(s State) string,
	filterButtons func(buttons []Button, state State, actionSpace ActionSpace) []int,
) (int, error) {
	solutions := make([]int, len(actionSpaces))
	wg := sync.WaitGroup{}
	chunkPerWorker := (len(actionSpaces) + numWorkers - 1) / numWorkers
	// log.Printf("solving %d action spaces with %d workers, %d chunks each", len(actionSpaces), numWorkers, chunkPerWorker)
	for w := range numWorkers {
		startIdx := w * chunkPerWorker
		endIdx := min((w+1)*chunkPerWorker, len(actionSpaces))
		wg.Go(
			func() {
				for i := startIdx; i < endIdx; i++ {
					solution, err := Solve(actionSpaces[i], maxIterations, equalityFunc, hashFunc, filterButtons)
					if err != nil {
						log.Printf("failed to solve action space %d: %v", i, err)
						continue
					}
					// log.Printf("solved action space %d: %d", i, solution)
					solutions[i] = solution
				}
			},
		)
	}
	wg.Wait()

	// for i, actionSpace := range actionSpaces {
	// 	solution, err := Solve(actionSpace, maxIterations, equalityFunc, hashFunc, filterButtons)
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

func Solve(
	actionSpace ActionSpace,
	maxIterations int,
	equalityFunc func(a, b State) bool,
	hashFunc func(s State) string,
	filterButtons func(buttons []Button, state State, actionSpace ActionSpace) []int,
) (int, error) {
	start := State{
		Lights:  make(Lights, len(actionSpace.Goal.Lights)),
		Joltage: make([]int, len(actionSpace.Goal.Lights)),
	}
	visited := make(map[string]bool)
	type Tracker struct {
		State State
		Steps int
	}
	queue := []Tracker{{State: start, Steps: 0}}
	iterations := 0
	buttonIndices := make([]int, len(actionSpace.Buttons))
	for i := range actionSpace.Buttons {
		buttonIndices[i] = i
	}
	for ; iterations < maxIterations && len(queue) > 0; iterations++ {
		current := queue[0]
		queue = queue[1:]
		if equalityFunc(current.State, actionSpace.Goal) {
			return current.Steps, nil
		}
		key := hashFunc(current.State)
		if visited[key] {
			continue
		}
		visited[key] = true
		if filterButtons != nil {
			buttonIndices = filterButtons(actionSpace.Buttons, current.State, actionSpace)
		}
		for _, buttonIdx := range buttonIndices {
			button := actionSpace.Buttons[buttonIdx]
			newState, err := button.Apply(current.State)
			if err != nil {
				return 0, err
			}
			queue = append(queue, Tracker{State: newState, Steps: current.Steps + 1})
		}
	}
	return 0, fmt.Errorf("no solution found after %d iterations, visited %d states", iterations, len(visited))
}

func VerifySolution(
	actionSpace ActionSpace,
	buttonSequence []int,
	equalityFunc func(a, b State) bool,
) (bool, error) {
	currentState := State{
		Lights:  make(Lights, len(actionSpace.Goal.Lights)),
		Joltage: make([]int, len(actionSpace.Goal.Joltage)),
	}

	for step, buttonIdx := range buttonSequence {
		if buttonIdx < 0 || buttonIdx >= len(actionSpace.Buttons) {
			return false, fmt.Errorf("button index %d at step %d is out of range", buttonIdx, step)
		}
		button := actionSpace.Buttons[buttonIdx]
		newState, err := button.Apply(currentState)
		if err != nil {
			return false, fmt.Errorf("error applying button %d at step %d: %w", buttonIdx, step, err)
		}
		currentState = newState
	}
	return equalityFunc(currentState, actionSpace.Goal), nil
}

func Part1(input io.Reader) (int, error) {
	actionSpaces, err := ParseInput(input)
	if err != nil {
		return 0, err
	}
	equalityFunc := func(a, b State) bool {
		if len(a.Lights) != len(b.Lights) {
			return false
		}
		for i := range a.Lights {
			if a.Lights[i] != b.Lights[i] {
				return false
			}
		}
		return true
	}
	hashFunc := func(s State) string {
		buf := bytes.Buffer{}
		buf.Grow(len(s.Lights))
		for _, light := range s.Lights {
			if light {
				buf.WriteByte('#')
			} else {
				buf.WriteByte('.')
			}
		}
		return buf.String()
	}
	var filterButtons func(buttons []Button, state State, actionSpace ActionSpace) []int = nil
	return MultiSolve(actionSpaces, maxIterations, equalityFunc, hashFunc, filterButtons)
}

func Part2(input io.Reader) (int, error) {
	actionSpaces, err := ParseInput(input)
	if err != nil {
		return 0, err
	}
	equalityFunc := func(a, b State) bool {
		if len(a.Joltage) != len(b.Joltage) {
			return false
		}
		for i := range a.Joltage {
			if a.Joltage[i] != b.Joltage[i] {
				return false
			}
		}
		return true
	}
	hashFunc := func(s State) string {
		buf := bytes.Buffer{}
		buf.Grow(len(s.Joltage) * 3)
		for i, joltage := range s.Joltage {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(strconv.Itoa(joltage))
		}
		return buf.String()
	}
	filterButtons := func(buttons []Button, state State, actionSpace ActionSpace) []int {
		validButtons := make([]int, 0, len(buttons))
		for i := range buttons {
			valid, err := buttons[i].WithinJoltageLimits(state, actionSpace.Goal.Joltage)
			if err != nil {
				panic(fmt.Errorf("error checking joltage limits for button %d: %w", i, err))
			}
			if valid {
				validButtons = append(validButtons, i)
			}
		}
		return validButtons
	}
	return MultiSolve(actionSpaces, maxIterations, equalityFunc, hashFunc, filterButtons)
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
