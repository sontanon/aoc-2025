package main

var evenDigitBoundaries = []Span{
	{10, 99},
	{1_000, 9_999},
	{100_000, 999_999},
	{10_000_000, 99_999_999},
	{1_000_000_000, 9_999_999_999},
	{100_000_000_000, 999_999_999_999},
	{10_000_000_000_000, 99_999_999_999_999},
	{1_000_000_000_000_000, 9_999_999_999_999_999},
	{100_000_000_000_000_000, 999_999_999_999_999_999},
	// 10^20 does not fit in int64
}

func (r Span) GetInvalidIdsPart1() []int {
	estimatedCapacity := max(4, (r.End-r.Start)/1_000)
	invalids := make([]int, 0, estimatedCapacity)
	for _, boundary := range evenDigitBoundaries {
		overlapStart := max(r.Start, boundary.Start)
		overlapEnd := min(r.End, boundary.End)
		if overlapStart > overlapEnd {
			continue
		}
		for i := overlapStart; i <= overlapEnd; i++ {
			if isInvalidPart1(i) {
				invalids = append(invalids, i)
			}
		}
	}
	return invalids
}

func (r Span) GetInvalidIdsPart2() []int {
	invalids := make([]int, 0)
	for i := r.Start; i <= r.End; i++ {
		if isInvalidPart2(i) {
			invalids = append(invalids, i)
		}
	}
	return invalids
}

func isInvalidPart1(id int) bool {
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

func isInvalidPart2(id int) bool {
	numDigits := countDigits(id)
	if numDigits < 2 {
		return false
	}

	divisors, ok := divisorsTable[numDigits]
	if !ok {
		return false
	}

	for _, numRepeats := range divisors {
		patternLen := numDigits / numRepeats
		divisor := pow10Table[patternLen]
		pattern := id % divisor // Extract the rightmost pattern

		// Verify all segments match the pattern
		temp := id
		matched := true
		for range numRepeats {
			if temp%divisor != pattern {
				matched = false
				break
			}
			temp /= divisor
		}
		if matched {
			return true
		}
	}
	return false
}

func countDigits(n int) int {
	count := 0
	for n > 0 {
		n /= 10
		count++
	}
	return count
}
