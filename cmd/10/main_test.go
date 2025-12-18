package main

import (
	"strings"
	"testing"
)

func TestPart1(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			"Provided example",
			`[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}
[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}
`,
			7,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				result, err := Part1(strings.NewReader(tt.input))
				if err != nil {
					t.Fatalf("Part1() error = %v", err)
				}
				if result != tt.expected {
					t.Errorf("Part1() = %v, want %v", result, tt.expected)
				}
			})
	}
}

func TestVerifySolution(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		buttonSequence []int
		ignoreJoltage  bool
		expected       bool
		errorExpected  bool
	}{
		{
			"Provided example 1 with joltage",
			`[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}`,
			[]int{0, 1, 1, 1, 3, 3, 3, 4, 5, 5},
			false,
			true,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				actionSpace, err := ParseLine([]byte(tt.input))
				if err != nil {
					t.Fatalf("ParseLine() error = %v", err)
				}
				result, err := VerifySolution(actionSpace, tt.buttonSequence, tt.ignoreJoltage)
				if (err != nil) != tt.errorExpected {
					t.Fatalf("VerifySolution() error = %v, errorExpected %v", err, tt.errorExpected)
				}
				if result != tt.expected {
					t.Errorf("VerifySolution() = %v, want %v", result, tt.expected)
				}
			})
	}
}

func TestPart2(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			"Provided example",
			`[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}
[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}
`,
			33,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				result, err := Part2(strings.NewReader(tt.input))
				if err != nil {
					t.Fatalf("Part2() error = %v", err)
				}
				if result != tt.expected {
					t.Errorf("Part2() = %v, want %v", result, tt.expected)
				}
			})
	}
}
