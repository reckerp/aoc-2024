package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var cache = make(map[string]int)

func applyRules(i int) []int {
	if i == 0 {
		return []int{1}
	}
	strNum := strconv.Itoa(i)
	if len(strNum)%2 == 0 {
		mid := len(strNum) / 2
		left, _ := strconv.Atoi(strNum[:mid])
		right, _ := strconv.Atoi(strNum[mid:])
		return []int{left, right}
	}
	return []int{i * 2024}
}

func countStones(num int, blinks int) int {
	// Base case: no blinks left
	if blinks == 0 {
		return 1
	}

	// Check cache
	key := fmt.Sprintf("%d,%d", num, blinks)
	if val, exists := cache[key]; exists {
		return val
	}

	// Apply rules and recursively count resulting stones
	result := 0
	for _, next := range applyRules(num) {
		result += countStones(next, blinks-1)
	}

	// Cache the result
	cache[key] = result
	return result
}

func calculateTotalStones(input []int, blinks int) int {
	total := 0
	for _, num := range input {
		total += countStones(num, blinks)
	}
	return total
}

func getInput(filename string) ([]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var integers []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		for _, field := range fields {
			num, err := strconv.Atoi(field)
			if err != nil {
				return nil, fmt.Errorf("failed to convert %s to integer: %w", field, err)
			}
			integers = append(integers, num)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return integers, nil
}

func main() {
	input, err := getInput("input.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	resultP1 := calculateTotalStones(input, 25)
	fmt.Println("[PART 1] After 25 blinks:", resultP1)
	resultP2 := calculateTotalStones(input, 75)
	fmt.Println("[PART 2] After 75 blinks:", resultP2)
}
