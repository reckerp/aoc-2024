package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	left, right, err := readInput("input.txt")
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return
	}

	totalDistance := calculateTotalDistance(left, right)
	similarityScore := calculateSimilarityScore(left, right)

	fmt.Printf("The total distance: %d\n", totalDistance)
	fmt.Printf("The similarity score between the lists is: %d\n", similarityScore)
}

// PART 1
func calculateTotalDistance(left, right []int) int {
	sort.Ints(left)
	sort.Ints(right)

	totalDistance := 0
	for i := range left {
		totalDistance += int(math.Abs(float64(left[i] - right[i])))
	}

	return totalDistance
}

// PART 2
func calculateSimilarityScore(left, right []int) int {
	// Build a frequency map for the right list
	rightFrequency := make(map[int]int)
	for _, num := range right {
		rightFrequency[num]++
	}

	// Calculate the similarity score
	similarityScore := 0
	for _, num := range left {
		similarityScore += num * rightFrequency[num]
	}

	return similarityScore
}

func readInput(filename string) ([]int, []int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	var left, right []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, nil, fmt.Errorf("invalid line format: %s", line)
		}

		leftNum, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, nil, err
		}
		rightNum, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, nil, err
		}

		left = append(left, leftNum)
		right = append(right, rightNum)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return left, right, nil
}
