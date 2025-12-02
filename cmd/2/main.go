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
	numWorkers       = 4
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

func (r Span) GetInvalidIds() []int {
	invalids := make([]int, 0)
	// This should be optimized by reducing the search space and automatically rejecting certain ranges.
	for i := r.Start; i <= r.End; i++ {
		if isInvalid(i) {
			invalids = append(invalids, i)
		}
	}
	return invalids
}

func isInvalid(id int) bool {
	// Option 1: strconv.Itoa: 0.01325s per full run
	// strId := strconv.Itoa(id)
	// numDigits := len(strId)
	// if numDigits%2 != 0 {
	// 	return false
	// }
	// return strId[:numDigits/2] == strId[numDigits/2:]

	// Option 2: math.Log10: 0.00662s per full run
	// numDigits := int(math.Log10(float64(id))) + 1
	// if numDigits%2 != 0 {
	// 	return false
	// }
	// left := id / int(math.Pow10(numDigits/2))
	// right := id % int(math.Pow10(numDigits/2))
	// return left == right

	// Option 3: single loop: 0.00238s per full run
	left := id
	right := 0
	multiplier := 1
	for {
		rightDigit := left % 10
		left /= 10
		right = rightDigit*multiplier + right
		multiplier *= 10
		if left == right {
			return true
		}
		if left < multiplier*10 {
			return false
		}
	}
}

func Part1(input io.Reader) (int, error) {
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
	for range numWorkers {
		wg.Go(func() {
			localSum := 0
			for span := range spans {
				invalids := span.GetInvalidIds()
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

	invalidIds, err := Part1(file)
	if err != nil {
		panic(err)
	}
	fmt.Println("Part 1:", invalidIds)
}
