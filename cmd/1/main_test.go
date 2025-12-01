package main

import (
	"os"
	"strings"
	"testing"
)

func TestPart1(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		startingPosition   int
		dialLength         int
		expectedZeroCounts int
	}{
		{
			"Provided example",
			`L68
L30
R48
L5
R60
L55
L1
L99
R14
L82`,
			startingPosition,
			dialLength,
			3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zeroCounts, err := Part1(
				strings.NewReader(tt.input),
				tt.startingPosition,
				tt.dialLength,
			)
			if err != nil {
				t.Fatalf("stage1() error = %v", err)
			}
			if zeroCounts != tt.expectedZeroCounts {
				t.Errorf("stage1() = %v, want %v", zeroCounts, tt.expectedZeroCounts)
			}
		})
	}
}

func BenchmarkPart1(b *testing.B) {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		b.Fatalf("failed to read input file: %v", err)
	}
	input := string(data)

	tests := []struct {
		name     string
		fn       func(string, int, int) (int, error)
		expected int
	}{
		{
			name: "Part 1",
			fn: func(input string, startPosition, dialLength int) (int, error) {
				return Part1(strings.NewReader(input), startPosition, dialLength)
			},
			expected: 1180,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs() // Show allocation stats
			for b.Loop() {
				result, err := tt.fn(input, startingPosition, dialLength)
				if err != nil {
					b.Fatalf("benchmark failed: %v", err)
				}
				if result != tt.expected {
					b.Fatalf("expected %d, got %d", tt.expected, result)
				}
			}
		})
	}
}
