package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

type BeamSplitters struct {
	StartingBeam int
	Splits       [][]int
	Width        int
}

func ParseInput(input io.Reader) (BeamSplitters, error) {
	lineScanner := bufio.NewScanner(input)

	result := BeamSplitters{}
	numLine := 0

	splits := make([][]int, 0)
	for lineScanner.Scan() {
		numLine++
		if numLine%2 == 0 {
			continue
		}
		line := lineScanner.Bytes()
		if numLine == 1 {
			result.Width = len(line)
			if result.Width == 0 {
				return BeamSplitters{}, errors.New("first line is empty")
			}
		}
		if len(line) != result.Width {
			return BeamSplitters{}, fmt.Errorf("unexpected line length of %d, expected %d", len(line), result.Width)
		}
		allocateSplits := true
		for i, b := range line {
			switch b {
			case 'S':
				result.StartingBeam = i
			case '^':
				if allocateSplits {
					splits = append(splits, make([]int, 0))
					allocateSplits = false
				}
				splits[len(splits)-1] = append(splits[len(splits)-1], i)
			}
		}
	}
	if err := lineScanner.Err(); err != nil {
		return BeamSplitters{}, err
	}
	if result.StartingBeam == 0 {
		return BeamSplitters{}, errors.New("invalid starting position or no starting position found")
	}
	result.Splits = splits
	return result, nil
}

type Beams struct {
	Positions map[int]struct{}
	Width     int
}

func (b *Beams) Split(splitters []int) (int, error) {
	if len(splitters) == 0 {
		return 0, nil
	}

	count := 0
	for _, splitterPos := range splitters {
		if splitterPos+1 > b.Width || splitterPos-1 < 0 {
			return 0, fmt.Errorf("invalid position %d would overflow bounds", splitterPos)
		}
		_, beamExists := b.Positions[splitterPos]
		if !beamExists {
			continue
		}
		count++
		delete(b.Positions, splitterPos)
		b.Positions[splitterPos+1] = struct{}{}
		b.Positions[splitterPos-1] = struct{}{}
	}
	return count, nil
}

func Part1(input io.Reader) (int, error) {
	bSps, err := ParseInput(input)
	if err != nil {
		return 0, err
	}

	beams := Beams{
		Positions: map[int]struct{}{
			bSps.StartingBeam: {},
		},
		Width: bSps.Width,
	}
	count := 0
	for _, splitters := range bSps.Splits {
		new, err := beams.Split(splitters)
		if err != nil {
			return 0, err
		}
		count += new
	}
	return count, nil
}

func getInputPath() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	return filepath.Join(dir, "input.txt")
}

func main() {
	file, err := os.Open(getInputPath())
	if err != nil {
		panic(err)
	}
	defer file.Close()

	result1, err := Part1(file)
	if err != nil {
		panic(err)
	}
	println("Part 1:", result1)

	// if _, err := file.Seek(0, io.SeekStart); err != nil {
	// 	panic(err)
	// }
	// result2, err := Part2(file)
	// println("Part 2:", result2)
}
