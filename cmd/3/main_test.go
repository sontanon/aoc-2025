package main

import (
	"os"
	"strings"
	"testing"
)

func TestMaximumJoltage(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{
			"987654321111111",
			98,
		}, {
			"811111111111119",
			89,
		},
		{
			"234234234234278",
			78,
		},
		{
			"818181911112111",
			92,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.input,
			func(t *testing.T) {
				result, err := ParseBank(tt.input)
				if err != nil {
					t.Fatalf("ParseBank() error = %v", err)
				}
				if result != tt.expected {
					t.Errorf("ParseBank() = %v, want %v", result, tt.expected)
				}
			},
		)
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
			`987654321111111
811111111111119
234234234234278
818181911112111
`,
			357,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := Part1(strings.NewReader(tt.input))
			if err != nil {
				t.Fatalf("Part1() error = %v", err)
			}
			if r != tt.expected {
				t.Errorf("Part1() = %v, want %v", r, tt.expected)
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
	expected := 17359

	b.Run("Part1", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			result, err := Part1(strings.NewReader(input))
			if err != nil {
				b.Fatalf("benchmark failed: %v", err)
			}
			if result != expected {
				b.Fatalf("expected %d, got %d", expected, result)
			}
		}
	})
}

func BenchmarkParseBank(b *testing.B) {
	data, err := os.ReadFile(getInputPath())
	if err != nil {
		b.Fatalf("failed to read input file: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	b.Run("ParseBank", func(b *testing.B) {
		b.ReportAllocs()
		idx := 0
		for b.Loop() {
			line := lines[idx%len(lines)]
			result, err := ParseBank(line)
			if err != nil {
				b.Fatalf("ParseBank failed: %v", err)
			}
			_ = result
			idx++
		}
	})
}
