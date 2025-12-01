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

type Direction string

const (
	DirectionInvalid Direction = ""
	DirectionLeft    Direction = "L"
	DirectionRight   Direction = "R"
)

type Rotation struct {
	Direction Direction
	Steps     int
}

func (r Rotation) Apply(currentPosition, dialLength int) int {
	switch r.Direction {
	case DirectionLeft:
		return (currentPosition - r.Steps + dialLength) % dialLength
	case DirectionRight:
		return (currentPosition + r.Steps) % dialLength
	default:
		panic("invalid direction in rotation")
	}
}

func ParseRotation(s string, dialLength int) (Rotation, error) {
	if len(s) < 2 {
		return Rotation{}, fmt.Errorf("invalid rotation string: %s", s)
	}

	dir := Direction(s[0:1])
	if dir != DirectionLeft && dir != DirectionRight {
		return Rotation{}, fmt.Errorf("invalid direction in rotation string: %s", s)
	}

	steps, err := strconv.Atoi(s[1:])
	if err != nil {
		return Rotation{}, fmt.Errorf("invalid steps in rotation string: %s", s)
	}

	return Rotation{
		Direction: dir,
		Steps:     steps % dialLength,
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
		currentPosition = rotation.Apply(currentPosition, dialLength)
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
}
