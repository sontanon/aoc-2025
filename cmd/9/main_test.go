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
			`7,1
11,1
11,7
9,7
9,5
2,5
2,3
7,3
`,
			50,
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
