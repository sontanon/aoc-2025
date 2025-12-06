package main

import (
	"cmp"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
)

const (
	numWorkers = 8
)

type Range struct {
	Start int
	End   int
}

type SparseRange struct {
	SubRanges     []Range
	GlobalMinimum int
	GlobalMaximum int
}

func (sr SparseRange) Normalize() SparseRange {
	n := len(sr.SubRanges)
	if n <= 1 {
		return sr
	}

	newRanges := make([]Range, 0, n)
	left := sr.SubRanges[0]
	for _, right := range sr.SubRanges[1:] {
		if right.Start <= left.End+1 {
			left.End = max(left.End, right.End)
		} else {
			newRanges = append(newRanges, left)
			left = right
		}
	}
	newRanges = append(newRanges, left)
	return SparseRange{newRanges, newRanges[0].Start, newRanges[len(newRanges)-1].End}
}

func (sr SparseRange) Contains(id int) bool {
	if len(sr.SubRanges) == 0 || id > sr.GlobalMaximum || id < sr.GlobalMinimum {
		return false
	}
	k, found := slices.BinarySearchFunc(sr.SubRanges, id, func(r Range, target int) int {
		return cmp.Compare(r.Start, target)
	})
	if found {
		return true
	}
	if k == 0 {
		return false
	}
	return id <= sr.SubRanges[k-1].End
}

func ParseInput(input io.Reader) (SparseRange, []int, error) {
	data, err := io.ReadAll(input)
	if err != nil {
		return SparseRange{}, nil, err
	}

	rangesStr, idsStr, found := strings.Cut(strings.TrimSpace(string(data)), "\n\n")
	if !found {
		return SparseRange{}, nil, errors.New("invalid input format")
	}

	rangeLines := strings.Split(strings.TrimSpace(rangesStr), "\n")
	ranges := make([]Range, 0, len(rangeLines))
	for _, line := range rangeLines {
		startStr, endStr, found := strings.Cut(line, "-")
		if !found {
			return SparseRange{}, nil, errors.New("invalid range format")
		}
		start, err := strconv.Atoi(startStr)
		if err != nil {
			return SparseRange{}, nil, err
		}
		end, err := strconv.Atoi(endStr)
		if err != nil {
			return SparseRange{}, nil, err
		}
		if end < start {
			return SparseRange{}, nil, errors.New("range end less than start")
		}
		ranges = append(ranges, Range{Start: start, End: end})
	}
	slices.SortFunc(ranges, func(a, b Range) int {
		return cmp.Compare(a.Start, b.Start)
	})

	idLines := strings.Split(strings.TrimSpace(idsStr), "\n")
	ids := make([]int, 0, len(idLines))
	for _, line := range idLines {
		id, err := strconv.Atoi(line)
		if err != nil {
			return SparseRange{}, nil, err
		}
		ids = append(ids, id)
	}
	if len(ranges) == 0 {
		return SparseRange{}, nil, errors.New("no ranges provided")
	}
	if len(ids) == 0 {
		return SparseRange{}, nil, errors.New("no ids provided")
	}
	sr := SparseRange{
		SubRanges:     ranges,
		GlobalMinimum: ranges[0].Start,
		GlobalMaximum: ranges[len(ranges)-1].End,
	}
	sr = sr.Normalize()
	return sr, ids, nil
}

func Part1Sequential(input io.Reader) (int, error) {
	sr, ids, err := ParseInput(input)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, id := range ids {
		if sr.Contains(id) {
			count++
		}
	}
	return count, nil
}

func Part1Parallel(input io.Reader, numWorkers int) (int, error) {
	sr, ids, err := ParseInput(input)
	if err != nil {
		return 0, err
	}
	idsPerWorker := (len(ids) + numWorkers - 1) / numWorkers
	results := make(chan int, numWorkers)
	wg := sync.WaitGroup{}
	for w := range numWorkers {
		wg.Go(func() {
			start := w * idsPerWorker
			if start >= len(ids) {
				results <- 0
				return
			}
			end := min((w+1)*idsPerWorker, len(ids))
			localCount := 0
			for _, id := range ids[start:end] {
				if sr.Contains(id) {
					localCount++
				}
			}
			results <- localCount
		})
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	total := 0
	for r := range results {
		total += r
	}
	return total, nil
}

func Part2Sequential(input io.Reader) (int, error) {
	sr, _, err := ParseInput(input)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, r := range sr.SubRanges {
		count += r.End - r.Start + 1
	}
	return count, nil
}

func Part2Parallel(input io.Reader, numWorkers int) (int, error) {
	sr, _, err := ParseInput(input)
	if err != nil {
		return 0, err
	}
	rangesPerWorker := (len(sr.SubRanges) + numWorkers - 1) / numWorkers
	results := make(chan int, numWorkers)
	wg := sync.WaitGroup{}
	for w := range numWorkers {
		wg.Go(func() {
			start := w * rangesPerWorker
			if start >= len(sr.SubRanges) {
				results <- 0
				return
			}
			end := min((w+1)*rangesPerWorker, len(sr.SubRanges))
			localCount := 0
			for _, r := range sr.SubRanges[start:end] {
				localCount += r.End - r.Start + 1
			}
			results <- localCount
		})
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	total := 0
	for r := range results {
		total += r
	}
	return total, nil
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

	result, err := Part1Sequential(file)
	if err != nil {
		panic(err)
	}
	println("Part 1:", result)

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		panic(err)
	}
	result2, err := Part2Sequential(file)
	if err != nil {
		panic(err)
	}
	println("Part 2:", result2)
}
