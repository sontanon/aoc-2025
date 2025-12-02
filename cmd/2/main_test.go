package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestGetInvalidIdsPart1(t *testing.T) {
	methods := []struct {
		name string
		fn   func(Span) []int
	}{
		{
			"Arithmetic",
			Span.GetInvalidIdsPart1,
		},
		{
			"Direct",
			func(s Span) []int {
				allInvalidsPart1 := generateAllInvalidsPart1()
				return s.GetInvalidIdsPart1Direct(allInvalidsPart1)
			},
		},
	}
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
	for _, method := range methods {
		for _, tt := range tests {
			t.Run(
				fmt.Sprintf("%s/%s", method.name, tt.name),
				func(t *testing.T) {
					span, err := ParseSpan(tt.name)
					if err != nil {
						t.Fatalf("ParseSpan error = %v", err)
					}
					result := method.fn(span)
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
}

func TestPart1(t *testing.T) {
	methods := []struct {
		name string
		fn   func(io.Reader) (int, error)
	}{
		{
			"Arithmetic",
			func(i io.Reader) (int, error) {
				return processSpansWithWorkers(i, Span.GetInvalidIdsPart1, numWorkers)
			},
		},
		{
			"Direct",
			func(i io.Reader) (int, error) {
				allInvalidsPart1 := generateAllInvalidsPart1()
				return processSpansWithWorkers(i, func(s Span) []int {
					return s.GetInvalidIdsPart1Direct(allInvalidsPart1)
				}, numWorkers)
			},
		},
	}
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
	for _, method := range methods {
		for _, tt := range tests {
			t.Run(fmt.Sprintf("%s/%s", method.name, tt.name),
				func(t *testing.T) {
					result, err := method.fn(strings.NewReader(tt.input))
					if err != nil {
						t.Fatalf("Part 1 error = %v", err)
					}
					if result != tt.expected {
						t.Errorf("Part 1 got = %v, want %v", result, tt.expected)
					}
				})
		}
	}
}

func TestGetInvalidIdsPart2(t *testing.T) {
	methods := []struct {
		name string
		fn   func(Span) []int
	}{
		{
			"Arithmetic",
			Span.GetInvalidIdsPart2,
		},
		{
			"Direct",
			func(s Span) []int {
				allInvalidsPart2 := generateAllInvalidsPart2()
				return s.GetInvalidIdsPart2Direct(allInvalidsPart2)
			},
		},
	}
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
			[]int{99, 111},
		},
		{
			"998-1012",
			[]int{999, 1010},
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
			[]int{565656},
		},
		{
			"824824821-824824827",
			[]int{824824824},
		},
		{
			"2121212118-2121212124",
			[]int{2121212121},
		},
	}
	for _, method := range methods {
		for _, tt := range tests {
			t.Run(fmt.Sprintf("%s/%s", method.name, tt.name), func(t *testing.T) {
				span, err := ParseSpan(tt.name)
				if err != nil {
					t.Fatalf("ParseSpan error = %v", err)
				}
				result := method.fn(span)
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
}

func TestPart2(t *testing.T) {
	methods := []struct {
		name string
		fn   func(io.Reader) (int, error)
	}{
		{
			"Arithmetic",
			func(i io.Reader) (int, error) {
				return processSpansWithWorkers(i, Span.GetInvalidIdsPart2, numWorkers)
			},
		},
		{
			"Direct",
			func(i io.Reader) (int, error) {
				allInvalidsPart2 := generateAllInvalidsPart2()
				return processSpansWithWorkers(i, func(s Span) []int {
					return s.GetInvalidIdsPart2Direct(allInvalidsPart2)
				}, numWorkers)
			},
		},
	}
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			"Provided example",
			`11-22,95-115,998-1012,1188511880-1188511890,222220-222224,1698522-1698528,446443-446449,38593856-38593862,565653-565659,824824821-824824827,2121212118-2121212124`,
			4174379265,
		},
	}

	for _, method := range methods {
		for _, tt := range tests {
			t.Run(fmt.Sprintf("%s/%s", method.name, tt.name), func(t *testing.T) {
				result, err := method.fn(strings.NewReader(tt.input))
				if err != nil {
					t.Fatalf("Part 2 error = %v", err)
				}
				if result != tt.expected {
					t.Errorf("Part 2 got = %v, want %v", result, tt.expected)
				}
			})
		}
	}
}

func BenchmarkPart1(b *testing.B) {
	data, err := os.ReadFile(getInputPath())
	if err != nil {
		b.Fatalf("failed to read input file: %v", err)
	}
	input := string(data)
	expected := 44487518055

	b.Run("Arithmetic", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			result, err := processSpans(strings.NewReader(input), Span.GetInvalidIdsPart1)
			if err != nil {
				b.Fatalf("benchmark failed: %v", err)
			}
			if result != expected {
				b.Fatalf("expected %d, got %d", expected, result)
			}
		}
	})

	b.Run("Direct", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			allInvalidsPart1 := generateAllInvalidsPart1()
			result, err := processSpans(strings.NewReader(input), func(s Span) []int {
				return s.GetInvalidIdsPart1Direct(allInvalidsPart1)
			})
			if err != nil {
				b.Fatalf("benchmark failed: %v", err)
			}
			if result != expected {
				b.Fatalf("expected %d, got %d", expected, result)
			}
		}
	})
}

func BenchmarkPart2(b *testing.B) {
	data, err := os.ReadFile(getInputPath())
	if err != nil {
		b.Fatalf("failed to read input file: %v", err)
	}
	input := string(data)
	expected := 53481866137

	b.Run("Arithmetic", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			result, err := processSpans(strings.NewReader(input), Span.GetInvalidIdsPart2)
			if err != nil {
				b.Fatalf("benchmark failed: %v", err)
			}
			if result != expected {
				b.Fatalf("expected %d, got %d", expected, result)
			}
		}
	})

	b.Run("Direct", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			allInvalidsPart2 := generateAllInvalidsPart2()
			result, err := processSpans(strings.NewReader(input), func(s Span) []int {
				return s.GetInvalidIdsPart2Direct(allInvalidsPart2)
			})
			if err != nil {
				b.Fatalf("benchmark failed: %v", err)
			}
			if result != expected {
				b.Fatalf("expected %d, got %d", expected, result)
			}
		}
	})
}

func BenchmarkInit(b *testing.B) {
	b.Run("GenerateAllInvalidsPart1", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = generateAllInvalidsPart1()
		}
	})

	b.Run("GenerateAllInvalidsPart2", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = generateAllInvalidsPart2()
		}
	})
}
