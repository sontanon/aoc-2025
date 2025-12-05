package main

import (
	"os"
	"strings"
	"testing"
)

func TestParseInput(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    Matrix
		expectError bool
	}{
		{
			"Valid input",
			`..@@.
@@@.@
@@@@@
@.@@@
@@.@@`,
			Matrix{[][]int{
				{0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 1, 1, 0, 0},
				{0, 1, 1, 1, 0, 1, 0},
				{0, 1, 1, 1, 1, 1, 0},
				{0, 1, 0, 1, 1, 1, 0},
				{0, 1, 1, 0, 1, 1, 0},
				{0, 0, 0, 0, 0, 0, 0},
			}, 7, 7},
			false,
		},
		{
			"Invalid character",
			`..@@.
@@A.@
@@@@@
@.@@@
@@.@@`,
			Matrix{},
			true,
		},
		{
			"Inconsistent row lengths",
			`..@@.
@@@.@
@@@@
@.@@@
@@.@@`,
			Matrix{},
			true,
		},
		{
			"Empty input",
			``,
			Matrix{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				result, err := ParseInput(
					strings.NewReader(tt.input))
				if (err != nil) != tt.expectError {
					t.Fatalf("ParseInput() error = %v, expectError %v", err, tt.expectError)
				}
				if !tt.expectError && !equalMatrices(result, tt.expected) {
					t.Errorf("ParseInput() = %v, want %v", result, tt.expected)
				}
			})
	}
}

func equalMatrices(a, b Matrix) bool {
	if a.rows != b.rows || a.cols != b.cols {
		return false
	}
	for i := 0; i < a.rows; i++ {
		for j := 0; j < a.cols; j++ {
			if a.data[i][j] != b.data[i][j] {
				return false
			}
		}
	}
	return true
}

func TestAccessibleRollsPart1(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			"Provided example",
			`..@@.@@@@.
@@@.@.@.@@
@@@@@.@.@@
@.@@@@..@.
@@.@@@@.@@
.@@@@@@@.@
.@.@.@.@@@
@.@@@.@@@@
.@@@@@@@@.
@.@.@@@.@.`,
			13,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				result, err := Part1(
					strings.NewReader(tt.input))
				if err != nil {
					t.Fatalf("Part1() error = %v", err)
				}
				if result != tt.expected {
					t.Errorf("Part1() = %v, want %v", result, tt.expected)
				}
			})
	}
}

func BenchmarkPart1(b *testing.B) {
	data, err := os.ReadFile(getInputPath())
	if err != nil {
		b.Fatalf("failed to read input file: %v", err)
	}
	input := string(data)
	expected := 1467

	b.Run("ParseWithWorkersChannels", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			matrix, err := ParseInput(strings.NewReader(input))
			if err != nil {
				b.Fatalf("benchmark failed: %v", err)
			}
			result := ParseWithWorkersChannels(matrix, numWorkers)
			if result != expected {
				b.Fatalf("expected %d, got %d", expected, result)
			}
		}
	})

	b.Run("ParseWithWorkersStatic", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			matrix, err := ParseInput(strings.NewReader(input))
			if err != nil {
				b.Fatalf("benchmark failed: %v", err)
			}
			result := ParseWithWorkersStatic(matrix, numWorkers)
			if result != expected {
				b.Fatalf("expected %d, got %d", expected, result)
			}
		}
	})

	b.Run("ParseWithWorkersMutex", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			matrix, err := ParseInput(strings.NewReader(input))
			if err != nil {
				b.Fatalf("benchmark failed: %v", err)
			}
			result := ParseWithWorkersMutex(matrix, numWorkers)
			if result != expected {
				b.Fatalf("expected %d, got %d", expected, result)
			}
		}
	})
}
