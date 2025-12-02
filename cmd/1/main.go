package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

const (
	inputFile        = "/home/santiago/Projects/aoc-2025/cmd/1/input.txt"
	startingPosition = 50
	dialLength       = 100
)

type Direction byte

const (
	DirectionLeft  Direction = 'L'
	DirectionRight Direction = 'R'
)

type Rotation struct {
	Direction      Direction
	Steps          int
	ExtraRotations int
}

func (r Rotation) Apply(startPosition, dialLength int) (int, bool) {
	switch r.Direction {
	case DirectionLeft:
		return (startPosition - r.Steps + dialLength) % dialLength, startPosition != 0 && startPosition-r.Steps < 0
	case DirectionRight:
		return (startPosition + r.Steps) % dialLength, startPosition != 0 && startPosition+r.Steps > dialLength
	default:
		panic("invalid direction in rotation")
	}
}

func ParseRotation(s string, dialLength int) (Rotation, error) {
	if len(s) < 2 {
		return Rotation{}, fmt.Errorf("invalid rotation string: %s", s)
	}

	dir := Direction(s[0])
	if dir != DirectionLeft && dir != DirectionRight {
		return Rotation{}, fmt.Errorf("invalid direction in rotation string: %s", s)
	}

	steps, err := strconv.Atoi(s[1:])
	if err != nil {
		return Rotation{}, fmt.Errorf("invalid steps in rotation string: %s", s)
	}

	normalizedSteps := steps % dialLength
	extraRotations := steps / dialLength

	return Rotation{
		Direction:      dir,
		Steps:          normalizedSteps,
		ExtraRotations: extraRotations,
	}, nil
}

func Part1(r io.Reader, startingPosition, dialLength int) (int, error) {
	scanner := bufio.NewScanner(r)
	lineNum := 0
	currentPosition := startingPosition
	zeroCounts := 0

	for scanner.Scan() {
		lineNum++
		rotation, err := ParseRotation(scanner.Text(), dialLength)
		if err != nil {
			return 0, fmt.Errorf("error parsing line %d, '%s', %w", lineNum, scanner.Text(), err)
		}
		currentPosition, _ = rotation.Apply(currentPosition, dialLength)
		if currentPosition == 0 {
			zeroCounts++
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error reading input: %w", err)
	}
	return zeroCounts, nil
}

func Part2(r io.Reader, startingPosition, dialLength int) (int, error) {
	scanner := bufio.NewScanner(r)
	lineNum := 0
	currentPosition := startingPosition
	crossedZero := false
	zeroCounts := 0

	for scanner.Scan() {
		lineNum++
		rotation, err := ParseRotation(scanner.Text(), dialLength)
		if err != nil {
			return 0, fmt.Errorf("error parsing line %d, '%s', %w", lineNum, scanner.Text(), err)
		}
		currentPosition, crossedZero = rotation.Apply(currentPosition, dialLength)
		zeroCounts += rotation.ExtraRotations
		if crossedZero {
			zeroCounts++
		}
		if currentPosition == 0 {
			zeroCounts++
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error reading input: %w", err)
	}
	return zeroCounts, nil
}

func main() {
	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	zeroCounts, err := Part1(file, startingPosition, dialLength)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Stage 1: %d\n", zeroCounts)

	file.Seek(0, io.SeekStart)
	zeroCounts, err = Part2(file, startingPosition, dialLength)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Stage 2: %d\n", zeroCounts)
}
