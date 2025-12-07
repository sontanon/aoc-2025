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
		expected    Worksheet
		expectError bool
	}{
		{
			"Provided example",
			`123 328  51 64 
 45 64  387 23 
  6 98  215 314
*   +   *   +  
`,
			Worksheet{
				Columns: [][]int{
					{123, 45, 6},
					{328, 64, 98},
					{51, 387, 215},
					{64, 23, 314},
				},
				Operations: []Operation{
					OperationMul,
					OperationAdd,
					OperationMul,
					OperationAdd,
				},
			},
			false,
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
				if !tt.expectError && !worksheetsEqual(result, tt.expected) {
					t.Errorf("ParseInput() = %v, want %v", result, tt.expected)
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
			`123 328  51 64 
 45 64  387 23 
  6 98  215 314
*   +   *   +  
`,
			4277556,
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
	expected := 7229350537438
	b.Run("Part1", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			result, err := Part1(strings.NewReader(string(data)))
			if err != nil {
				b.Fatalf("Part1() error = %v", err)
			}
			if result != expected {
				b.Fatalf("Part1() = %v, want %v", result, expected)
			}
		}
	})
}

func worksheetsEqual(a, b Worksheet) bool {
	if len(a.Columns) != len(b.Columns) {
		return false
	}
	for i := range a.Columns {
		if len(a.Columns[i]) != len(b.Columns[i]) {
			return false
		}
		for j := range a.Columns[i] {
			if a.Columns[i][j] != b.Columns[i][j] {
				return false
			}
		}
	}
	if len(a.Operations) != len(b.Operations) {
		return false
	}
	for i := range a.Operations {
		if a.Operations[i] != b.Operations[i] {
			return false
		}
	}
	return true
}
