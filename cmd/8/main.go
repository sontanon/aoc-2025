package main

import (
	"bufio"
	"cmp"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
)

type Vector [3]int

type Batch []Vector

type PairDistance struct {
	Distance float64
	Indices  [2]int
}

type Circuit map[int]struct{}

type Circuits []Circuit

func ParseInput(input io.Reader) (Batch, error) {
	lineScanner := bufio.NewScanner(input)

	result := make([]Vector, 0)
	numLine := 0
	for lineScanner.Scan() {
		numLine++
		line := lineScanner.Text()
		splits := strings.Split(line, ",")
		if len(splits) != 3 {
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

func CalculateDistances(b Batch) ([]PairDistance, error) {
	n := len(b)
	if n < 2 {
		return nil, errors.New("vector batch must contain at least one pair")
	}

	pds := make([]PairDistance, 0, n*(n-1)/2)
	for i := range n {
		for j := i + 1; j < n; j++ {
			pds = append(pds, PairDistance{
				Distance: Distance(b[i], b[j]),
				Indices:  [2]int{i, j},
			})
		}
	}
	slices.SortFunc(pds, func(a, b PairDistance) int {
		return cmp.Compare(a.Distance, b.Distance)
	})
	return pds, nil
}

func Distance(u, v Vector) float64 {
	return math.Sqrt(float64((u[0]-v[0])*(u[0]-v[0]) + (u[1]-v[1])*(u[1]-v[1]) + (u[2]-v[2])*(u[2]-v[2])))
}

func Connect(cs Circuits, a, b int) Circuits {
	aIdx := -1
	bIdx := -1

	for i := range cs {
		if aIdx == -1 {
			if _, exists := cs[i][a]; exists {
				aIdx = i
			}
		}
		if bIdx == -1 {
			if _, exists := cs[i][b]; exists {
				bIdx = i
			}
		}
		if aIdx != -1 && bIdx != -1 {
			break
		}
	}

	if aIdx == -1 && bIdx == -1 {
		return append(cs, map[int]struct{}{a: {}, b: {}})
	}
	if aIdx == bIdx {
		return cs
	}
	if aIdx != -1 && bIdx == -1 {
		cs[aIdx][b] = struct{}{}
		return cs
	}
	if bIdx != -1 && aIdx == -1 {
		cs[bIdx][a] = struct{}{}
		return cs
	}

	smallIdx := min(aIdx, bIdx)
	bigIdx := max(aIdx, bIdx)
	for node := range cs[bigIdx] {
		cs[smallIdx][node] = struct{}{}
	}
	return append(cs[:bigIdx], cs[bigIdx+1:]...)
}

func Part1(input io.Reader, numConnections int) (int, error) {
	batch, err := ParseInput(input)
	if err != nil {
		return 0, err
	}

	pds, err := CalculateDistances(batch)
	if err != nil {
		return 0, err
	}

	cs := make(Circuits, 0)
	limit := min(numConnections, len(pds))

	for i := range limit {
		a := pds[i].Indices[0]
		b := pds[i].Indices[1]
		cs = Connect(cs, a, b)
	}

	if len(cs) < 3 {
		return 0, errors.New("did not create at least three distinct circuit groups")
	}

	// Extract circuit sizes
	sizes := make([]int, len(cs))
	for i, circuit := range cs {
		sizes[i] = len(circuit)
	}
	slices.Sort(sizes)

	return sizes[len(sizes)-1] * sizes[len(sizes)-2] * sizes[len(sizes)-3], nil
}

func Part2(input io.Reader) (int, error) {
	batch, err := ParseInput(input)
	if err != nil {
		return 0, err
	}

	pds, err := CalculateDistances(batch)
	if err != nil {
		return 0, err
	}

	cs := make(Circuits, 0)
	targetSize := len(batch)

	for i := range pds {
		a := pds[i].Indices[0]
		b := pds[i].Indices[1]
		cs = Connect(cs, a, b)

		for _, circuit := range cs {
			if len(circuit) == targetSize {
				return batch[a][0] * batch[b][0], nil
			}
		}
	}

	return 0, errors.New("could not connect all vectors into a single circuit")
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

	result1, err := Part1(file, 1_000)
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
