package main

import (
	"errors"
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
	bankBufferSize   = 8
	resultBufferSize = 64
	numWorkers       = 8
)

type Bank struct {
	v []int
}

func ParseBank(input string) (Bank, error) {
	v := make([]int, len(input))
	for i, ch := range input {
		e, err := strconv.Atoi(string(ch))
		if err != nil {
			return Bank{}, err
		}
		v[i] = e
	}
	if len(v) < 2 {
		return Bank{}, fmt.Errorf("input must have at least two elements: '%s'", input)
	}
	if v[0] == 0 {
		return Bank{}, errors.New("first element cannot be zero")
	}
	return Bank{v: v}, nil
}

func (b Bank) MaximumJoltage() int {
	leftMaxIndex := 0
	for i := 1; i < len(b.v)-1; i++ {
		if b.v[leftMaxIndex] == 9 {
			break
		}
		if b.v[i] > b.v[leftMaxIndex] {
			leftMaxIndex = i
		}
	}
	rightMaxIndex := len(b.v) - 1
	for i := len(b.v) - 2; i > leftMaxIndex; i-- {
		if b.v[rightMaxIndex] == 9 {
			break
		}
		if b.v[i] > b.v[rightMaxIndex] {
			rightMaxIndex = i
		}
	}
	return 10*b.v[leftMaxIndex] + b.v[rightMaxIndex]
}

func processBanksWithWorkers(input io.Reader, validator func(Bank) int, workers int) (int, error) {
	data, err := io.ReadAll(input)
	if err != nil {
		return 0, fmt.Errorf("error reading input: %w", err)
	}

	bankLines := strings.Split(strings.TrimSpace(string(data)), "\n")
	banks := make(chan Bank, bankBufferSize)
	results := make(chan int, resultBufferSize)

	go func() {
		for i, line := range bankLines {
			bank, err := ParseBank(line)
			if err != nil {
				panic(fmt.Errorf("error parsing bank on line %d: %w", i+1, err))
			}
			banks <- bank
		}
		close(banks)
	}()

	wg := sync.WaitGroup{}
	for range workers {
		wg.Go(func() {
			for bank := range banks {
				result := validator(bank)
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
	return processBanksWithWorkers(input, func(b Bank) int {
		return b.MaximumJoltage()
	}, numWorkers)
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
