package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type Operation string

const (
	OperationInvalid Operation = ""
	OperationAdd     Operation = "+"
	OperationMul     Operation = "*"
)

type Worksheet struct {
	Columns    [][]int
	Operations []Operation
}

func ParseInput(input io.Reader) (Worksheet, error) {
	lineScanner := bufio.NewScanner(input)

	var previousLine string
	hasContent := false
	ws := Worksheet{}

	for lineScanner.Scan() {
		currentLine := lineScanner.Text()
		if hasContent {
			if err := ws.addColumn(previousLine); err != nil {
				return Worksheet{}, err
			}
		}
		hasContent = true
		previousLine = currentLine
	}
	if err := lineScanner.Err(); err != nil {
		return Worksheet{}, err
	}
	if err := ws.addOperations(previousLine); err != nil {
		return Worksheet{}, err
	}
	return ws, nil
}

func (ws *Worksheet) addColumn(line string) error {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return errors.New("empty line found when parsing columns")
	}

	if ws.Columns == nil {
		ws.Columns = make([][]int, len(fields))
	}

	if len(fields) != len(ws.Columns) {
		return fmt.Errorf("line has %d columns, expected %d", len(fields), len(ws.Columns))
	}

	for i, field := range fields {
		value, err := strconv.Atoi(field)
		if err != nil {
			return fmt.Errorf("invalid integer '%s' at column %d: %w", field, i, err)
		}
		ws.Columns[i] = append(ws.Columns[i], value)
	}

	return nil
}

func (ws *Worksheet) addOperations(line string) error {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return errors.New("empty line found when parsing operations")
	}
	if len(fields) != len(ws.Columns) {
		return fmt.Errorf("number of operations %d does not match number of columns %d", len(fields), len(ws.Columns))
	}

	tempBuffer := make([]Operation, 0)
	for i, field := range fields {
		operation := Operation(field)
		if operation != OperationAdd && operation != OperationMul {
			return fmt.Errorf("invalid operation '%s' at column %d", string(operation), i)
		}
		tempBuffer = append(tempBuffer, operation)
	}
	ws.Operations = tempBuffer
	return nil
}

func (ws Worksheet) Calculate() int {
	results := make([]int, len(ws.Columns))

	for j := range ws.Columns {
		switch ws.Operations[j] {
		case OperationAdd:
			results[j] = ApplyAddition(ws.Columns[j])
		case OperationMul:
			results[j] = ApplyMultiplication(ws.Columns[j])
		default:
			panic(fmt.Errorf("invalid operation '%s' at position %d", string(ws.Operations[j]), j))
		}
	}

	result := 0
	for j := range results {
		result += results[j]
	}
	return result
}

func ApplyAddition(col []int) int {
	result := 0
	for j := range col {
		result += col[j]
	}
	return result
}

func ApplyMultiplication(col []int) int {
	result := 1
	for j := range col {
		result *= col[j]
	}
	return result
}

func Part1(input io.Reader) (int, error) {
	ws, err := ParseInput(input)
	if err != nil {
		return 0, err
	}
	return ws.Calculate(), nil
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

	result1, err := Part1(file)
	if err != nil {
		panic(err)
	}
	println("Part 1:", result1)
}
