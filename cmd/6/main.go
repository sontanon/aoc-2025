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

func Part2(input io.Reader) (int, error) {
	data, err := io.ReadAll(input)
	if err != nil {
		return 0, err
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	numLines := len(lines)
	if numLines < 2 {
		return 0, fmt.Errorf("input must contain a t least two lines, received %d", numLines)
	}
	ops := strings.Fields(lines[numLines-1])
	for i, op := range ops {
		if op != string(OperationAdd) && op != string(OperationMul) {
			return 0, fmt.Errorf("invalid operation at position %d: '%s'", i, op)
		}
	}
	lines = lines[0 : numLines-1]
	lineLength := len(lines[0])
	for i := range lines {
		if len(lines[i]) != lineLength {
			return 0, fmt.Errorf("inconsistent line length for line %d, received length %d but expected %d", i+1, len(lines[i]), lineLength)
		}
	}

	numCols := len(ops)
	right := lineLength - 1
	result := 0
	for j := numCols - 1; j >= 0; j-- {
		op := Operation(ops[j])
		columnResult, newRight, err := operateColumn(lines, right, op)
		if err != nil {
			return 0, fmt.Errorf("error processing column %d: %w", j, err)
		}
		result += columnResult
		right = newRight
	}
	return result, nil
}

func operateColumn(lines []string, right int, op Operation) (int, int, error) {
	if right < 0 {
		return 0, 0, fmt.Errorf("initial right offset cannot be less than zero, received %d", right)
	}
	n := len(lines)
	if n == 0 {
		return 0, 0, errors.New("received empty lines")
	}
	buffer := make([]byte, n)

	result := 0
	if op == OperationMul {
		result = 1
	}

	allSpaces := false
	for ; right >= 0 && !allSpaces; right-- {
		countSpaces := 0
		for i := range lines {
			b := lines[i][right]
			if b != ' ' && (b < '0' || b > '9') {
				return 0, 0, fmt.Errorf("received unexpected byte at line %d, offset %d: %q", i+1, right, b)
			}
			if b == ' ' {
				countSpaces++
			}
			buffer[i] = b
		}
		if countSpaces == n {
			allSpaces = true
			continue
		}
		subResult := 0
		multiplier := 1
		for i := len(lines) - 1; i >= 0; i-- {
			if buffer[i] == ' ' {
				continue
			}
			subResult += multiplier * int(buffer[i]-'0')
			multiplier *= 10
		}
		switch op {
		case OperationAdd:
			result += subResult
		case OperationMul:
			result *= subResult
		}
	}
	return result, right, nil
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

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		panic(err)
	}
	result2, err := Part2(file)
	println("Part 2:", result2)
}
