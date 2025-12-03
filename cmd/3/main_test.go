package main

import (
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
			92,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.input,
			func(t *testing.T) {
				bank, err := ParseBank(tt.input)
				if err != nil {
					t.Fatalf("ParseBank() error = %v", err)
				}
				result := bank.MaximumJoltage()
				if result != tt.expected {
					t.Errorf("MaximumJoltage() = %v, want %v", result, tt.expected)
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
