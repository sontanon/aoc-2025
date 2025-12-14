package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type Vector [2]int

type Batch []Vector

func ParseInput(input io.Reader) (Batch, error) {
	lineScanner := bufio.NewScanner(input)
	result := make([]Vector, 0)
	numLine := 0
	for lineScanner.Scan() {
		numLine++
		line := lineScanner.Text()
		splits := strings.Split(line, ",")
		if len(splits) != 2 {
			return nil, fmt.Errorf("line %d does not contain valid input: %s", numLine, line)
		}
		var v Vector
		for i, s := range splits {
			val, err := strconv.Atoi(s)
			if err != nil {
				return nil, fmt.Errorf("line %d contains invalid integer: %s", numLine, s)
			}
			v[i] = val
		}
		result = append(result, v)
	}
	if err := lineScanner.Err(); err != nil {
		return nil, err
	}
	if len(result) < 2 {
		return nil, errors.New("input does not contain at least a pair of vectors")
	}

	return result, nil
}

func Part1(input io.Reader) (int, error) {
	batch, err := ParseInput(input)
	if err != nil {
		return 0, err
	}
	maxArea := 0
	for i := range len(batch) {
		u := batch[i]
		for j := i + 1; j < len(batch); j++ {
			v := batch[j]
			l := u[0] - v[0]
			w := u[1] - v[1]
			if l == 0 || w == 0 {
				continue
			}
			if l < 0 {
				l = -l
			}
			if w < 0 {
				w = -w
			}
			l++
			w++
			area := l * w
			if area > maxArea {
				maxArea = area
			}
		}
	}
	return maxArea, nil
}

func Part2(input io.Reader) (int, error) {
	return 0, nil
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
	fmt.Println("Part 1:", result1)

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		panic(err)
	}

	result2, err := Part2(file)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 2:", result2)
}
