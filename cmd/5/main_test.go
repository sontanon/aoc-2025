package main

import (
	"os"
	"slices"
	"strings"
	"testing"
)

func TestParseInput(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedRange SparseRange
		expectedIds   []int
		expectError   bool
	}{
		{
			"Provided example",
			`3-5
10-14
16-20
12-18

1
5
8
11
17
32`,
			SparseRange{
				SubRanges: []Range{
					{Start: 3, End: 5},
					{Start: 10, End: 20},
				},
				GlobalMinimum: 3,
				GlobalMaximum: 20,
			},
			[]int{1, 5, 8, 11, 17, 32},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				sr, ids, err := ParseInput(
					strings.NewReader(tt.input))
				if (err != nil) != tt.expectError {
					t.Fatalf("ParseInput() error = %v, expectError %v", err, tt.expectError)
				}
				if !tt.expectError {
					if !equalInputs(sr, tt.expectedRange) {
						t.Errorf("ParseInput() range = %v, want %v", sr, tt.expectedRange)
					}
					if !slices.Equal(ids, tt.expectedIds) {
						t.Errorf("ParseInput() ids = %v, want %v", ids, tt.expectedIds)
					}
				}
			})
	}
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			"Provided example",
			`3-5
10-14
16-20
12-18

1
5
8
11
17
32`,
			3,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				result, err := Part1(
					strings.NewReader(tt.input), 1)
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
		b.Fatalf("Failed to read input file: %v", err)
	}
	expected := 761
	b.Run("Part1", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			result, err := Part1(strings.NewReader(string(data)), numWorkers)
			if err != nil {
				b.Fatalf("Part1() error = %v", err)
			}
			if result != expected {
				b.Fatalf("Part1() = %v, want %v", result, expected)
			}
		}
	})
}

func equalInputs(a, b SparseRange) bool {
	if len(a.SubRanges) != len(b.SubRanges) {
		return false
	}
	for i := range a.SubRanges {
		if a.SubRanges[i] != b.SubRanges[i] {
			return false
		}
	}
	if a.GlobalMinimum != b.GlobalMinimum || a.GlobalMaximum != b.GlobalMaximum {
		return false
	}
	return true
}
