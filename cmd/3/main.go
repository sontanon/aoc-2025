package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

const (
	bankBufferSize   = 8
	resultBufferSize = 64
	numWorkers       = 8
)

type Bank struct {
	v []int
}

func ParseBank(input string) (int, error) {
	if len(input) <= 2 {
		return 0, errors.New("input too short to be a valid bank")
	}
	lByte := input[0]
	if lByte < '0' || lByte > '9' {
		return 0, fmt.Errorf("invalid left digit: %q", lByte)
	}
	rByte := input[len(input)-1]
	if rByte < '0' || rByte > '9' {
		return 0, fmt.Errorf("invalid right digit: %q", rByte)
	}
	lIdx := 0
	for i := 1; i < len(input)-1; i++ {
		if lByte == '9' {
			break
		}
		if input[i] < '0' || input[i] > '9' {
			return 0, fmt.Errorf("invalid digit at index %d: %q", i, input[i])
		}
		if input[i] > lByte {
			lByte = input[i]
			lIdx = i
		}
	}
	for i := len(input) - 2; i > lIdx; i-- {
		if rByte == '9' {
			break
		}
		if input[i] < '0' || input[i] > '9' {
			return 0, fmt.Errorf("invalid digit at index %d: %q", i, input[i])
		}
		if input[i] > rByte {
			rByte = input[i]
		}
	}
	lDigit := int(lByte - '0')
	rDigit := int(rByte - '0')
	return 10*lDigit + rDigit, nil
}

func processBanksWithWorkers(input io.Reader, parser func(string) (int, error), workers int) (int, error) {
	data, err := io.ReadAll(input)
	if err != nil {
		return 0, fmt.Errorf("error reading input: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	banks := make(chan string, bankBufferSize)
	results := make(chan int, resultBufferSize)

	go func() {
		for _, line := range lines {
			banks <- line
		}
		close(banks)
	}()

	wg := sync.WaitGroup{}
	for range workers {
		wg.Go(func() {
			for bank := range banks {
				result, err := parser(bank)
				if err != nil {
					panic(fmt.Sprintf("error parsing bank %q: %v", bank, err))
				}
				results <- result
			}
		})
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	total := 0
	for result := range results {
		total += result
	}
	return total, nil
}

func Part1(input io.Reader) (int, error) {
	return processBanksWithWorkers(input,
		ParseBank, numWorkers)
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

	result, err := Part1(file)
	if err != nil {
		panic(err)
	}
	println("Part 1:", result)

}
