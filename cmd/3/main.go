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

func ParseBankPart1(input string) (int, error) {
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

const part2Length = 12

func ParseBankPart2(input string) (int, error) {
	n := len(input)
	if n < part2Length {
		return 0, errors.New("input too short to be a valid bank for part 2")
	}
	var buffer [part2Length]byte
	searchSpace := []byte(input)
	for i := range part2Length {
		leaveRoom := part2Length - i - 1
		lIdx := searchMaxByte(searchSpace[:len(searchSpace)-leaveRoom])
		if lIdx == -1 {
			return 0, fmt.Errorf("invalid digit found during left search at iteration %d", i)
		}
		buffer[i] = searchSpace[lIdx]
		searchSpace = searchSpace[lIdx+1:]
	}
	result := 0
	multiplier := 1
	for i := part2Length - 1; i >= 0; i-- {
		digit := int(buffer[i] - '0')
		result += digit * multiplier
		multiplier *= 10
	}
	return result, nil
}

func validByte(b byte) bool {
	return b >= '0' && b <= '9'
}

func searchMaxByte(inputSlice []byte) int {
	if len(inputSlice) == 0 {
		return -1
	}
	if len(inputSlice) == 1 {
		return 0
	}
	lByte := inputSlice[0]
	if !validByte(lByte) {
		return -1
	}
	lIdx := 0
	for i := 1; i < len(inputSlice); i++ {
		if lByte == '9' {
			return lIdx
		}
		if !validByte(inputSlice[i]) {
			return -1
		}
		if inputSlice[i] > lByte {
			lByte = inputSlice[i]
			lIdx = i
		}
	}
	return lIdx
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
		ParseBankPart1, numWorkers)
}

func Part2(input io.Reader) (int, error) {
	return processBanksWithWorkers(input,
		ParseBankPart2, numWorkers)
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

	file.Seek(0, io.SeekStart)
	result, err = Part2(file)
	if err != nil {
		panic(err)
	}
	println("Part 2:", result)
}
