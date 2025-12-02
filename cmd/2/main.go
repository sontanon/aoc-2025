package main

import (
	"flag"
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

func Part1(input io.Reader, methodChoice string, allInvalidsPart1 []int) (int, error) {
	var validator func(Span) []int
	if methodChoice == "arithmetic" {
		validator = Span.GetInvalidIdsPart1
	} else {
		validator = func(s Span) []int {
			return s.GetInvalidIdsPart1Direct(allInvalidsPart1)
		}
	}
	return processSpans(input, validator)
}

func Part2(input io.Reader, methodChoice string, allInvalidsPart2 []int) (int, error) {
	var validator func(Span) []int
	if methodChoice == "arithmetic" {
		validator = Span.GetInvalidIdsPart2
	} else {
		validator = func(s Span) []int {
			return s.GetInvalidIdsPart2Direct(allInvalidsPart2)
		}
	}
	return processSpans(input, validator)
}

func getInputPath() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	return filepath.Join(dir, "input.txt")
}

func main() {
	methodFlag := flag.String("method", "arithmetic", "validation method: arithmetic|direct")
	flag.Parse()
	methodChoice := *methodFlag
	if methodChoice != "arithmetic" && methodChoice != "direct" {
		panic(fmt.Errorf("invalid method choice: %s", methodChoice))
	}

	var allInvalidsPart1 []int
	var allInvalidsPart2 []int

	if methodChoice == "arithmetic" {
		allInvalidsPart1 = nil
		allInvalidsPart2 = nil
	} else {
		allInvalidsPart1 = generateAllInvalidsPart1()
		allInvalidsPart2 = generateAllInvalidsPart2()
	}

	file, err := os.Open(getInputPath())
	if err != nil {
		panic(err)
	}
	defer file.Close()

	idSum, err := Part1(file, methodChoice, allInvalidsPart1)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 1:", idSum)

	file.Seek(0, io.SeekStart)
	idSum, err = Part2(file, methodChoice, allInvalidsPart2)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 2:", idSum)
}
