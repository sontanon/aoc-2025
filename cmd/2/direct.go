package main

import (
	"slices"
)

var pow10Table = []int{
	1,
	10,
	100,
	1_000,
	10_000,
	100_000,
	1_000_000,
	10_000_000,
	100_000_000,
	1_000_000_000,
	// 10_000_000_000,
	// 100_000_000_000,
	// 1_000_000_000_000,
	// 10_000_000_000_000,
	// 100_000_000_000_000,
	// 1_000_000_000_000_000,
	// 10_000_000_000_000_000,
	// 100_000_000_000_000_000,
	// 1_000_000_000_000_000_000,
	// 10^20 does not fit in int64
}

var divisorsTable = map[int][]int{
	2:  {2},
	3:  {3},
	4:  {2, 4},
	5:  {5},
	6:  {2, 3, 6},
	7:  {7},
	8:  {2, 4, 8},
	9:  {3, 9},
	10: {2, 5, 10},
	// 11: {11},
	// 12: {2, 3, 4, 6, 12},
	// 13: {13},
	// 14: {2, 7, 14},
	// 15: {3, 5, 15},
	// 16: {2, 4, 8, 16},
	// 17: {17},
	// 18: {2, 3, 6, 9, 18},
	// 19: {19},
	// 20: {2, 4, 5, 10, 20}, // does not fit in int64
}

func generateAllInvalidsPart1() []int {
	invalids := make([]int, 0, 100000)

	for ab := 1; ab <= 9; ab++ {
		invalids = append(invalids, ab*11)
	}
	for ab := 10; ab <= 99; ab++ {
		invalids = append(invalids, ab*100+ab)
	}
	for abc := 100; abc <= 999; abc++ {
		invalids = append(invalids, abc*1000+abc)
	}
	for abcd := 1000; abcd <= 9999; abcd++ {
		invalids = append(invalids, abcd*10000+abcd)
	}
	for abcde := 10000; abcde <= 99999; abcde++ {
		invalids = append(invalids, abcde*100000+abcde)
	}
	// for abcdef := 100000; abcdef <= 999999; abcdef++ {
	// 	invalids = append(invalids, abcdef*1000000+abcdef)
	// }
	// for abcdefg := 1000000; abcdefg <= 9999999; abcdefg++ {
	// 	invalids = append(invalids, abcdefg*10000000+abcdefg)
	// }
	// for abcdefgh := 10000000; abcdefgh <= 99999999; abcdefgh++ {
	// 	invalids = append(invalids, abcdefgh*100000000+abcdefgh)
	// }
	// for abcdefghi := 100000000; abcdefghi <= 999999999; abcdefghi++ {
	// 	invalids = append(invalids, abcdefghi*1000000000+abcdefghi)
	// }

	return invalids
}

func generateAllInvalidsPart2() []int {
	invalidSet := make(map[int]struct{})

	for numDigits := 2; numDigits <= 10; numDigits++ {
		divisors := divisorsTable[numDigits]
		if divisors == nil {
			continue
		}

		for _, numRepeats := range divisors {
			patternLen := numDigits / numRepeats
			patternMin := 1
			if patternLen > 1 {
				patternMin = pow10Table[patternLen-1] // e.g., 10 for 2-digit patterns
			}
			patternMax := pow10Table[patternLen] - 1 // e.g., 99 for 2-digit patterns

			for pattern := patternMin; pattern <= patternMax; pattern++ {
				num := 0
				multiplier := 1
				for range numRepeats {
					num += pattern * multiplier
					multiplier *= pow10Table[patternLen]
				}
				invalidSet[num] = struct{}{}
			}
		}
	}

	invalids := make([]int, 0, len(invalidSet))
	for num := range invalidSet {
		invalids = append(invalids, num)
	}
	slices.Sort(invalids)
	return invalids
}

func findInvalidsInRange(start, end int, invalids []int) []int {
	left := 0
	right := len(invalids)
	for left < right {
		mid := (left + right) / 2
		if invalids[mid] < start {
			left = mid + 1
		} else {
			right = mid
		}
	}
	startIdx := left

	left = startIdx
	right = len(invalids)
	for left < right {
		mid := (left + right) / 2
		if invalids[mid] <= end {
			left = mid + 1
		} else {
			right = mid
		}
	}
	endIdx := left

	if startIdx >= endIdx {
		return nil
	}
	return invalids[startIdx:endIdx]
}

func (r Span) GetInvalidIdsPart1Direct(allInvalidsPart1 []int) []int {
	return findInvalidsInRange(r.Start, r.End, allInvalidsPart1)
}

func (r Span) GetInvalidIdsPart2Direct(allInvalidsPart2 []int) []int {
	return findInvalidsInRange(r.Start, r.End, allInvalidsPart2)
}
