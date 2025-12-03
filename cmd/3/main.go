package main

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
)

type Bank struct {
	buffer string
}

func ParseBank(input string) (Bank, error) {
	return Bank{buffer: input}, nil
}

func (b Bank) MaximumJoltage() int {
	return 0
}

func Part1(input io.Reader) (int, error) {
	return 0, nil
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
