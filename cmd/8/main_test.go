package main

import (
	"strings"
	"testing"
)

func TestParseInput(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    Batch
		expectError bool
	}{
		{
			"Provided example",
			`162,817,812
57,618,57
906,360,560
592,479,940
352,342,300
466,668,158
542,29,236
431,825,988
739,650,466
52,470,668
216,146,977
819,987,18
117,168,530
805,96,715
346,949,466
970,615,88
941,993,340
862,61,35
984,92,344
425,690,689
`,
			Batch{
				{162, 817, 812},
				{57, 618, 57},
				{906, 360, 560},
				{592, 479, 940},
				{352, 342, 300},
				{466, 668, 158},
				{542, 29, 236},
				{431, 825, 988},
				{739, 650, 466},
				{52, 470, 668},
				{216, 146, 977},
				{819, 987, 18},
				{117, 168, 530},
				{805, 96, 715},
				{346, 949, 466},
				{970, 615, 88},
				{941, 993, 340},
				{862, 61, 35},
				{984, 92, 344},
				{425, 690, 689},
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
				if !tt.expectError {
					if len(result) != len(tt.expected) {
						t.Fatalf("ParseInput() = %v, want %v", result, tt.expected)
					}
					for i := range result {
						if result[i] != tt.expected[i] {
							t.Errorf("ParseInput()[%d] = %v, want %v", i, result[i], tt.expected[i])
						}
					}
				}
			})
	}
}

func TestPart1(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		numConnections int
		expected       int
	}{
		{
			"Provided example",
			`162,817,812
57,618,57
906,360,560
592,479,940
352,342,300
466,668,158
542,29,236
431,825,988
739,650,466
52,470,668
216,146,977
819,987,18
117,168,530
805,96,715
346,949,466
970,615,88
941,993,340
862,61,35
984,92,344
425,690,689
`,
			10,
			40,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				result, err := Part1(
					strings.NewReader(tt.input), tt.numConnections)
				if err != nil {
					t.Fatalf("Part1() error = %v", err)
				}
				if result != tt.expected {
					t.Errorf("Part1() = %v, want %v", result, tt.expected)
				}
			})
	}
}
