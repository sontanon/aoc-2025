package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
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
	// log.Printf("Loaded %d vectors in batch\n", len(batch))

	uniqueX := make(map[int]struct{})
	uniqueY := make(map[int]struct{})
	for _, v := range batch {
		uniqueX[v[0]] = struct{}{}
		uniqueY[v[1]] = struct{}{}
	}
	// log.Printf("Found %d unique X and %d unique Y coordinates\n", len(uniqueX), len(uniqueY))

	Xs := make([]int, 0, len(uniqueX)+2)
	Ys := make([]int, 0, len(uniqueY)+2)
	Xs = append(Xs, 0)
	Ys = append(Ys, 0)
	for x := range uniqueX {
		Xs = append(Xs, x)
	}
	for y := range uniqueY {
		Ys = append(Ys, y)
	}
	slices.Sort(Xs)
	slices.Sort(Ys)
	Xs = append(Xs, Xs[len(Xs)-1]+1)
	Ys = append(Ys, Ys[len(Ys)-1]+1)
	// log.Printf("Sorted unique coordinates\n")

	xToIndex := make(map[int]int, len(Xs))
	for i, x := range Xs {
		xToIndex[x] = i
	}
	yToIndex := make(map[int]int, len(Ys))
	for j, y := range Ys {
		yToIndex[y] = j
	}
	// log.Printf("Constructed coordinate to index maps\n")

	mask := make([][]int, len(Ys))
	for j := range mask {
		mask[j] = make([]int, len(Xs))
	}
	for k, u := range batch {
		v := batch[(k+1)%len(batch)]
		xi, yi := xToIndex[u[0]], yToIndex[u[1]]
		xf, yf := xToIndex[v[0]], yToIndex[v[1]]
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
	// log.Printf("Constructed optimized mask of polygon edges\n")
	// printMask(mask, len(Xs), len(Ys))

	type vertEdge struct{ x, yMin, yMax int }
	verticalEdges := make([]vertEdge, 0, len(batch)/2)
	for k, u := range batch {
		v := batch[(k+1)%len(batch)]
		if u[0] == v[0] { // Vertical edge
			xi := xToIndex[u[0]]
			yi, yf := yToIndex[u[1]], yToIndex[v[1]]
			if yi > yf {
				yi, yf = yf, yi
			}
			verticalEdges = append(verticalEdges, vertEdge{xi, yi, yf})
		}
	}

	crossings := make([]int, 0, len(verticalEdges))
	for j := 0; j < len(Ys)-1; j++ {
		crossings = crossings[:0]
		for _, e := range verticalEdges {
			if e.yMin <= j && j < e.yMax {
				crossings = append(crossings, e.x)
			}
		}
		slices.Sort(crossings)
		for c := 0; c+1 < len(crossings); c += 2 {
			for i := crossings[c]; i < crossings[c+1]; i++ {
				mask[j][i] = 1
			}
		}
	}

	// log.Printf("Filled interior of polygon\n")
	// printMask(mask, len(Xs), len(Ys))

	maxArea := 0
	// candidate1 := 0
	// candidate2 := 1
	for k := range batch {
		u := batch[k]
		xi, yi := xToIndex[u[0]], yToIndex[u[1]]
	CandidateLoop:
		for l := k + 1; l < len(batch); l++ {
			v := batch[l]
			xf, yf := xToIndex[v[0]], yToIndex[v[1]]
			for y := min(yi, yf); y <= max(yi, yf); y++ {
				for x := min(xi, xf); x <= max(xi, xf); x++ {
					if mask[y][x] == 0 {
						continue CandidateLoop
					}
				}
			}
			width := u[0] - v[0]
			height := u[1] - v[1]
			if width < 0 {
				width = -width
			}
			if height < 0 {
				height = -height
			}
			area := (width + 1) * (height + 1)
			if area > maxArea {
				maxArea = area
				// candidate1 = k
				// candidate2 = l
			}
		}
	}
	// log.Printf("Found max area between vectors %v and %v: %d\n", batch[candidate1], batch[candidate2], maxArea)

	return maxArea, nil
}

func printMask(mask [][]int, xRange, yRange int) {
	line := strings.Builder{}
	line.Grow(xRange)
	for j := range yRange {
		line.Reset()
		for i := range xRange {
			if mask[j][i] == 1 {
				line.WriteByte('#')
			} else {
				line.WriteByte('.')
			}
		}
		fmt.Println(line.String())
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
