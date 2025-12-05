package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

const (
	numWorkers = 8
)

type Matrix struct {
	data [][]int
	rows int
	cols int
}

func ParseInput(input io.Reader) (Matrix, error) {
	data, err := io.ReadAll(input)
	if err != nil {
		return Matrix{}, fmt.Errorf("error reading input: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	rows := len(lines)
	paddedRows := rows + 2
	if rows == 0 {
		return Matrix{}, fmt.Errorf("input is empty")
	}
	matrix := make([][]int, paddedRows)

	cols := len(lines[0])
	if cols == 0 {
		return Matrix{}, fmt.Errorf("invalid row length: %d", cols)
	}
	paddedCols := cols + 2

	matrix[0] = make([]int, paddedCols)
	matrix[len(lines)+1] = make([]int, paddedCols)
	for i, line := range lines {
		if len(line) != cols {
			return Matrix{}, fmt.Errorf("inconsistent row lengths: expected %d, got %d", cols, len(line))
		}
		matrix[i+1] = make([]int, paddedCols)
		for j, char := range line {
			switch char {
			case '@':
				matrix[i+1][j+1] = 1
			case '.':
				// matrix[i+1][j+1] = 0
			default:
				return Matrix{}, fmt.Errorf("invalid character '%c' in input", char)
			}
		}
	}
	return Matrix{data: matrix, rows: paddedRows, cols: paddedCols}, nil
}

func ParseElement(m [][]int, row, col int) int {
	if m[row][col] == 0 {
		return 0
	}
	neighbors := m[row-1][col-1] + m[row-1][col] + m[row-1][col+1] + m[row][col-1] + m[row][col+1] + m[row+1][col-1] + m[row+1][col] + m[row+1][col+1]

	if neighbors >= 4 {
		return 0
	}
	return 1
}

func ParseRow(m [][]int, row, startCol, endCol int) int {
	sum := 0
	for j := startCol; j < endCol; j++ {
		sum += ParseElement(m, row, j)
	}
	return sum
}

func ParseWithWorkersChannels(m Matrix, workers int) int {
	rows := make(chan int, workers)
	results := make(chan int, workers)

	go func() {
		for i := 1; i < m.rows-1; i++ {
			rows <- i
		}
		close(rows)
	}()

	wg := sync.WaitGroup{}
	for range workers {
		wg.Go(func() {
			for row := range rows {
				rowSum := ParseRow(m.data, row, 1, m.cols-1)
				results <- rowSum
			}
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
	return total
}

func ParseWithWorkersStatic(m Matrix, workers int) int {
	rowsPerWorker := (m.rows - 2 + workers - 1) / workers

	var wg sync.WaitGroup
	resultsChan := make(chan int, workers)

	for w := range workers {
		wg.Go(func() {
			start := 1 + w*rowsPerWorker
			end := min(start+rowsPerWorker, m.rows-1)

			sum := 0
			for row := start; row < end; row++ {
				sum += ParseRow(m.data, row, 1, m.cols-1)
			}
			resultsChan <- sum
		})
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	total := 0
	for r := range resultsChan {
		total += r
	}
	return total
}

func ParseWithWorkersMutex(m Matrix, workers int) int {
	rows := make(chan int, workers)

	go func() {
		for i := 1; i < m.rows-1; i++ {
			rows <- i
		}
		close(rows)
	}()

	var wg sync.WaitGroup
	var mu sync.Mutex
	total := 0

	for range workers {
		wg.Go(func() {
			for row := range rows {
				rowSum := ParseRow(m.data, row, 1, m.cols-1)
				mu.Lock()
				total += rowSum
				mu.Unlock()
			}
		})
	}

	wg.Wait()
	return total
}

func Part1(input io.Reader) (int, error) {
	matrix, err := ParseInput(input)
	if err != nil {
		return 0, fmt.Errorf("error parsing input: %w", err)
	}
	result := ParseWithWorkersStatic(matrix, numWorkers)
	return result, nil
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
