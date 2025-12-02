package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

const (
	spanBufferSize   = 8
	resultBufferSize = 64
	numWorkers       = 8
)

type Span struct {
	Start int
	End   int
}

func ParseSpan(s string) (Span, error) {
	before, after, found := strings.Cut(s, "-")
	if !found {
		return Span{}, fmt.Errorf("invalid span string: %s", s)
	}
	start, err := strconv.Atoi(before)
	if err != nil {
		return Span{}, fmt.Errorf("invalid start in span string: %s", s)
	}
	end, err := strconv.Atoi(after)
	if err != nil {
		return Span{}, fmt.Errorf("invalid end in span string: %s", s)
	}
	if start <= 0 || end <= 0 || end < start {
		return Span{}, fmt.Errorf("invalid span range in span string: %s", s)
	}
	return Span{Start: start, End: end}, nil
}

var allInvalidsPart1 []int

func init() {
	allInvalidsPart1 = generateAllInvalidsPart1()
}

func generateAllInvalidsPart1() []int {
	invalids := make([]int, 0, 100000)

	for ab := 1; ab <= 9; ab++ {
		invalids = append(invalids, ab*11)
	}

	for ab := 10; ab <= 99; ab++ {
		invalids = append(invalids, ab*100+ab)
	}

	for abc := 100; abc <= 999; abc++ {
		invalids = append(invalids, abc*1000+abc)
	}

	for abcd := 1000; abcd <= 9999; abcd++ {
		invalids = append(invalids, abcd*10000+abcd)
	}

	for abcde := 10000; abcde <= 99999; abcde++ {
		invalids = append(invalids, abcde*100000+abcde)
	}

	return invalids
}

func findInvalidsInRange(start, end int) []int {
	left := 0
	right := len(allInvalidsPart1)
	for left < right {
		mid := (left + right) / 2
		if allInvalidsPart1[mid] < start {
			left = mid + 1
		} else {
			right = mid
		}
	}
	startIdx := left

	left = startIdx
	right = len(allInvalidsPart1)
	for left < right {
		mid := (left + right) / 2
		if allInvalidsPart1[mid] <= end {
			left = mid + 1
		} else {
			right = mid
		}
	}
	endIdx := left

	if startIdx >= endIdx {
		return nil
	}
	return allInvalidsPart1[startIdx:endIdx]
}

func (r Span) GetInvalidIdsPart1() []int {
	return findInvalidsInRange(r.Start, r.End)
}

func (r Span) GetInvalidIdsPart2() []int {
	invalids := make([]int, 0)
	for i := r.Start; i <= r.End; i++ {
		if isInvalidPart2(i) {
			invalids = append(invalids, i)
		}
	}
	return invalids
}

func isInvalidPart2(id int) bool {
	strId := strconv.Itoa(id)
	numDigits := len(strId)
	if numDigits < 2 {
		return false
	}
	for numCandidates := 2; numCandidates <= numDigits; numCandidates++ {
		if numDigits%numCandidates != 0 {
			continue
		}
		candidateLength := numDigits / numCandidates
		candidate := strId[:candidateLength]
		matched := true
		for i := 1; i < numCandidates; i++ {
			if strId[i*candidateLength:(i+1)*candidateLength] != candidate {
				matched = false
				break
			}
		}
		if matched {
			return true
		}
	}
	return false
}

func processSpans(input io.Reader, validator func(Span) []int) (int, error) {
	return processSpansWithWorkers(input, validator, numWorkers)
}

func processSpansWithWorkers(input io.Reader, validator func(Span) []int, workers int) (int, error) {
	data, err := io.ReadAll(input)
	if err != nil {
		return 0, fmt.Errorf("error reading input: %w", err)
	}

	spanStrs := strings.Split(strings.TrimSpace(string(data)), ",")

	spans := make(chan Span, spanBufferSize)
	results := make(chan int, resultBufferSize)

	go func() {
		for i, spanStr := range spanStrs {
			span, err := ParseSpan(spanStr)
			if err != nil {
				panic(fmt.Errorf("error parsing span on index %d: %w", i, err))
			}
			spans <- span
		}
		close(spans)
	}()

	wg := sync.WaitGroup{}
	for range workers {
		wg.Go(func() {
			localSum := 0
			for span := range spans {
				invalids := validator(span)
				for _, id := range invalids {
					localSum += id
				}
			}
			results <- localSum
		})
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	totalSum := 0
	for sum := range results {
		totalSum += sum
	}

	return totalSum, nil
}

func Part1(input io.Reader) (int, error) {
	return processSpans(input, Span.GetInvalidIdsPart1)
}

func Part2(input io.Reader) (int, error) {
	return processSpans(input, Span.GetInvalidIdsPart2)
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

	idSum, err := Part1(file)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 1:", idSum)

	file.Seek(0, io.SeekStart)
	idSum, err = Part2(file)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 2:", idSum)
}
