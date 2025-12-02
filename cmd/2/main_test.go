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
			`11-22,95-115,998-1012,1188511880-1188511890,222220-222224,1698522-1698528,446443-446449,38593856-38593862,565653-565659,824824821-824824827,2121212118-2121212124`,
			1227775554,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Part1(strings.NewReader(tt.input))
			if err != nil {
				t.Fatalf("Part 1 error = %v", err)
			}
			if result != tt.expected {
				t.Errorf("Part 1 got = %v, want %v", result, tt.expected)
			}
		})
	}
}
