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
	numWorkers          = 8
	maxSafetyIterations = 1_000_000
	maxNeighbors        = 4
)

type Matrix struct {
	data [][]int
	rows int
	cols int
}

func (m Matrix) CalculateMask(mask [][]int, workers int) int {
	rowsPerWorker := (m.rows - 2 + workers - 1) / workers
	resultsChan := make(chan int, workers)

	var wg sync.WaitGroup
	for w := range workers {
		wg.Go(func() {
			start := 1 + w*rowsPerWorker
			end := min(start+rowsPerWorker, m.rows-1)

			canBeRemoved := 0
			for row := start; row < end; row++ {
				for col := 1; col < m.cols-1; col++ {
					mask[row][col] = ParseElement(m.data, row, col)
					canBeRemoved += mask[row][col]
				}
			}
			resultsChan <- canBeRemoved
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

func (m *Matrix) ApplyMask(mask [][]int, workers int) {
	rowsPerWorker := (m.rows - 2 + workers - 1) / workers

	var wg sync.WaitGroup
	for w := range workers {
		wg.Go(func() {
			start := 1 + w*rowsPerWorker
			end := min(start+rowsPerWorker, m.rows-1)

			for row := start; row < end; row++ {
				for col := 1; col < m.cols-1; col++ {
					m.data[row][col] -= mask[row][col]
					mask[row][col] = 0
				}
			}
		})
	}
	wg.Wait()
}

func (m *Matrix) RemoveRolls(workers, maxIterations int) int {
	if workers < 1 || maxIterations < 1 {
		return 0
	}

	mask := make([][]int, m.rows)
	for i := range mask {
		mask[i] = make([]int, m.cols)
	}

	removed := 0
	for range maxIterations {
		removedThisRound := m.CalculateMask(mask, workers)
		if removedThisRound == 0 {
			break
		}
		removed += removedThisRound
		m.ApplyMask(mask, workers)
	}
	return removed
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

	if neighbors >= maxNeighbors {
		return 0
	}
	return 1
}

func Part1(input io.Reader) (int, error) {
	matrix, err := ParseInput(input)
	if err != nil {
		return 0, fmt.Errorf("error parsing input: %w", err)
	}
	result := matrix.RemoveRolls(numWorkers, 1)
	return result, nil
}

func Part2(input io.Reader) (int, error) {
	matrix, err := ParseInput(input)
	if err != nil {
		return 0, fmt.Errorf("error parsing input: %w", err)
	}
	result := matrix.RemoveRolls(numWorkers, maxSafetyIterations)
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

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		panic(err)
	}
	result, err = Part2(file)
	if err != nil {
		panic(err)
	}
	println("Part 2:", result)
}
