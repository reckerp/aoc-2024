package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

func readInput(pattern string) ([][]string, error) {
	// Open the file
	file, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the entire content of the file
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Compile the regex pattern
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	// Find all matches
	return re.FindAllStringSubmatch(string(content), -1), nil
}

func part1() ([][]int, error) {
	matches, err := readInput(`mul\(\s*(\d+)[^\d]+(\d+)\s*\)`)
	if err != nil {
		return nil, err
	}
	// Prepare the result slice
	var result [][]int

	// Iterate over the matches and convert them to integers
	for _, match := range matches {
		if len(match) == 3 {
			num1, err1 := strconv.Atoi(match[1])
			num2, err2 := strconv.Atoi(match[2])
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("error converting string to int: %v, %v", err1, err2)
			}
			result = append(result, []int{num1, num2})
		}
	}

	return result, nil
}

func part2() ([][]int, error) {
	matches, err := readInput(`(?:mul\(\s*(\d+)[^\d]+(\d+)\s*\)|don't\(\)|do\(\))`)
	if err != nil {
		return nil, err
	}

	// Prepare the result slice
	var result [][]int

	// mul enabled
	isEnabled := true

	// Iterate over the matches and convert them to integers
	for _, match := range matches {
		switch match[0] {
		case "do()":
			isEnabled = true
		case "don't()":
			isEnabled = false
		default:
			if len(match) == 3 && isEnabled {
				num1, err1 := strconv.Atoi(match[1])
				num2, err2 := strconv.Atoi(match[2])
				if err1 != nil || err2 != nil {
					return nil, fmt.Errorf("error converting string to int: %v, %v", err1, err2)
				}
				result = append(result, []int{num1, num2})
			}
		}
	}

	return result, nil
}

func main() {
	input1, err := part1()
	if err != nil {
		panic(err)
	}

	input2, err := part2()
	if err != nil {
		panic(err)
	}

	sumPart1 := 0
	for _, match := range input1 {
		sumPart1 += match[0] * match[1]
	}

	sumPart2 := 0
	for _, match := range input2 {
		sumPart2 += match[0] * match[1]
	}

	println("The sum in part 1 is: ", sumPart1)
	println("The sum in part 2 is: ", sumPart2)

}
