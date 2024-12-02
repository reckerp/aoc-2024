package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

func readInputMatrix() [][]int {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var matrix [][]int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		row := make([]int, len(fields))
		for i, field := range fields {
			num, err := strconv.Atoi(field)
			if err != nil {
				log.Fatal(err)
			}
			row[i] = num
		}
		matrix = append(matrix, row)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return matrix
}

func isSafe(report []int) bool {
	isIncreasing := true
	isDecreasing := true

	for i := 1; i < len(report); i++ {
		diff := report[i] - report[i-1]

		if diff > 3 || diff < -3 || diff == 0 {
			return false
		}

		if diff > 0 {
			isDecreasing = false
		} else if diff < 0 {
			isIncreasing = false
		}
	}

	// Return true if the report is either fully increasing or fully decreasing
	return isIncreasing || isDecreasing
}

func canBeMadeValid(report []int) bool {
	for i := range report {
		newSlice := make([]int, len(report))
		copy(newSlice, report)

		newSlice = slices.Delete(newSlice, i, i+1)

		if isSafe(newSlice) {
			return true
		}
	}
	return false
}

func sumSafeReports(matrix [][]int) int {
	sum := 0
	for _, report := range matrix {
		if isSafe(report) {
			sum++
		}
	}

	return sum
}

func sumSafeReportsWithDampeners(matrix [][]int) int {
	sum := 0
	for _, report := range matrix {
		// Check if the report is safe without any removal
		if isSafe(report) {
			sum++
		} else if canBeMadeValid(report) {
			sum++
		}

	}
	return sum
}

func main() {
	input := readInputMatrix()

	sumSafe := sumSafeReports(input)
	sumSafe2 := sumSafeReportsWithDampeners(input)
	fmt.Println("Number of safe reports:", sumSafe)
	fmt.Println("Number of safe reports with dampeners:", sumSafe2)
}
