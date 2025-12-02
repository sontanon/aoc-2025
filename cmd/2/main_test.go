package main

import (
	"os"
	"strings"
	"testing"
)

func TestGetInvalidIds(t *testing.T) {
	tests := []struct {
		name     string
		expected []int
	}{
		{
			"11-22",
			[]int{11, 22},
		},
		{
			"95-115",
			[]int{99},
		},
		{
			"998-1012",
			[]int{1010},
		},
		{
			"1188511880-1188511890",
			[]int{1188511885},
		},
		{
			"222220-222224",
			[]int{222222},
		},
		{
			"1698522-1698528",
			[]int{},
		},
		{
			"446443-446449",
			[]int{446446},
		},
		{
			"38593856-38593862",
			[]int{38593859},
		},
		{
			"565653-565659",
			[]int{},
		},
		{
			"824824821-824824827",
			[]int{},
		},
		{
			"2121212118-2121212124",
			[]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			span, err := ParseSpan(tt.name)
			if err != nil {
				t.Fatalf("ParseSpan error = %v", err)
			}
			result := span.GetInvalidIds()
			if len(result) != len(tt.expected) {
				t.Fatalf("GetInvalidIds unexpected length, got = %v, want %v", result, tt.expected)
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("GetInvalidIds[%d] got = %v, want %v", i, result[i], tt.expected[i])
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

func BenchmarkPart1(b *testing.B) {
	data, err := os.ReadFile(getInputPath())
	if err != nil {
		b.Fatalf("failed to read input file: %v", err)
	}
	input := string(data)

	tests := []struct {
		name     string
		fn       func(string) (int, error)
		expected int
	}{
		{
			name: "Part 1",
			fn: func(input string) (int, error) {
				return Part1(strings.NewReader(input))
			},
			expected: 44487518055,
		},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				result, err := tt.fn(input)
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
