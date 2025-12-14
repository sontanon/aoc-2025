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

type Vector [3]float64

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
		x, err := strconv.Atoi(splits[0])
		if err != nil {
			return nil, fmt.Errorf("line %d contains invalid integer: %s", numLine, splits[0])
		}
		y, err := strconv.Atoi(splits[1])
		if err != nil {
			return nil, fmt.Errorf("line %d contains invalid integer: %s", numLine, splits[1])
		}
		z, err := strconv.Atoi(splits[2])
		if err != nil {
			return nil, fmt.Errorf("line %d contains invalid integer: %s", numLine, splits[2])
		}
		result = append(result, Vector{float64(x), float64(y), float64(z)})
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
			pds = append(pds,
				PairDistance{
					Distance: Distance(b[i], b[j]),
					Indices: [2]int{
						i, j,
					},
				},
			)
		}
	}
	slices.SortFunc(pds, func(a, b PairDistance) int {
		return cmp.Compare(a.Distance, b.Distance)
	})
	return pds, nil
}

func Distance(u, v Vector) float64 {
	return math.Sqrt((u[0]-v[0])*(u[0]-v[0]) + (u[1]-v[1])*(u[1]-v[1]) + (u[2]-v[2])*(u[2]-v[2]))
}

func Connect(cs Circuits, a, b int) Circuits {
	aIdx := -1
	bIdx := -1

	for i := 0; i < len(cs) && (aIdx == -1 || bIdx == -1); i++ {
		if _, aExists := cs[i][a]; aExists {
			aIdx = i
		}
		if _, bExists := cs[i][b]; bExists {
			bIdx = i
		}
	}

	if aIdx == -1 && bIdx == -1 {
		sc := map[int]struct{}{
			a: {},
			b: {},
		}
		cs = append(cs, sc)
		return cs
	}
	if aIdx != -1 && bIdx != -1 && aIdx == bIdx {
		return cs
	}
	if aIdx == -1 && bIdx != -1 {
		cs[bIdx][a] = struct{}{}
		return cs
	}
	if bIdx == -1 && aIdx != -1 {
		cs[aIdx][b] = struct{}{}
		return cs
	}
	smallIdx := min(aIdx, bIdx)
	bigIdx := max(aIdx, bIdx)
	for node := range cs[bigIdx] {
		cs[smallIdx][node] = struct{}{}
	}
	cs = append(cs[:bigIdx], cs[bigIdx+1:]...)
	return cs
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

	for i := range min(numConnections, len(pds)) {
		a := pds[i].Indices[0]
		b := pds[i].Indices[1]
		cs = Connect(cs, a, b)
	}

	if len(cs) < 3 {
		return 0, errors.New("did not create at least three distinct circuit groups")
	}

	ls := make([]int, len(cs))
	for i := range ls {
		ls[i] = len(cs[i])
	}

	slices.Sort(ls)
	return ls[len(ls)-1] * ls[len(ls)-2] * ls[len(ls)-3], nil
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
	println("Part 1:", result1)

	// if _, err := file.Seek(0, io.SeekStart); err != nil {
	// 	panic(err)
	// }
	// result2, err := Part2(file)
	// println("Part 2:", result2)
}
