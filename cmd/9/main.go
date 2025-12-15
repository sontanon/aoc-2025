package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
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
	batch, err := ParseInput(input)
	if err != nil {
		return 0, err
	}
	log.Printf("Loaded %d vectors in batch\n", len(batch))

	xMin, xMax := batch[0][0], batch[0][0]
	yMin, yMax := batch[0][1], batch[0][1]
	for _, v := range batch[1:] {
		x, y := v[0], v[1]
		if x < xMin {
			xMin = x
		}
		if x > xMax {
			xMax = x
		}
		if y < yMin {
			yMin = y
		}
		if y > yMax {
			yMax = y
		}
	}
	log.Printf("X range: %d to %d\n", xMin, xMax)
	log.Printf("Y range: %d to %d\n", yMin, yMax)

	for k, u := range batch {
		batch[k] = Vector{u[0] - xMin, u[1] - yMin}
	}
	log.Printf("Normalized vectors to origin\n")

	yRange := yMax - yMin + 1
	xRange := xMax - xMin + 1
	mask := make([][]int, yRange)
	for j := range mask {
		mask[j] = make([]int, xRange)
	}
	for k, u := range batch {
		v := batch[(k+1)%len(batch)]
		xi, yi := u[0], u[1]
		xf, yf := v[0], v[1]
		switch {
		case xf-xi > 0:
			for x := xi; x < xf; x++ {
				mask[yi][x] = 1
			}
		case xi-xf > 0:
			for x := xi; x > xf; x-- {
				mask[yi][x] = 1
			}
		case yf-yi > 0:
			for y := yi; y < yf; y++ {
				mask[y][xi] = 1
			}
		case yi-yf > 0:
			for y := yi; y > yf; y-- {
				mask[y][xi] = 1
			}
		}
	}
	log.Printf("Constructed mask of polygon edges\n")
	// printMask(mask, xRange, yRange)

	for j := range yRange {
		for i := 1; i < xRange; i++ {
			if mask[j][i] == 0 && mask[j][i-1] == 1 {
				// Fill the entire run of zeros in one go
				for i < xRange && mask[j][i] == 0 {
					mask[j][i] = 1
					i++
				}
			}
		}
	}
	log.Printf("Filled interior of polygon\n")
	// printMask(mask, xRange, yRange)

	maxArea := 0
	for k := range len(batch) {
		u := batch[k]
		xi, yi := u[0], u[1]
	CandidateLoop:
		for l := k + 1; l < len(batch); l++ {
			v := batch[l]
			xf, yf := v[0], v[1]
			deltaX := xf - xi
			switch {
			case deltaX > 0:
				for x := 1; x < deltaX; x++ {
					if mask[yi][xi+x] == 0 || mask[yf][xf-x] == 0 {
						continue CandidateLoop
					}
				}
			case -deltaX > 0:
				for x := 1; x < -deltaX; x++ {
					if mask[yi][xi-x] == 0 || mask[yf][xf+x] == 0 {
						continue CandidateLoop
					}
				}
			}
			deltaY := yf - yi
			switch {
			case deltaY > 0:
				for y := 1; y < deltaY; y++ {
					if mask[yi+y][xi] == 0 || mask[yf-y][xf] == 0 {
						continue CandidateLoop
					}
				}
			case -deltaY > 0:
				for y := 1; y < -deltaY; y++ {
					if mask[yi-y][xi] == 0 || mask[yf+y][xf] == 0 {
						continue CandidateLoop
					}
				}
			}

			if deltaX < 0 {
				deltaX = -deltaX
			}
			if deltaY < 0 {
				deltaY = -deltaY
			}
			area := (deltaX + 1) * (deltaY + 1)
			if area > maxArea {
				maxArea = area
			}
		}
	}

	return maxArea, nil
}
func printMask(mask [][]int, xRange, yRange int) {
	for j := range yRange {
		for i := range xRange {
			if mask[j][i] == 1 {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
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
