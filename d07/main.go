package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func evaluateExpression(numbers []int, operators []string) int {
	result := numbers[0]
	for i := 0; i < len(operators); i++ {
		switch operators[i] {
		case "+":
			result += numbers[i+1]
		case "*":
			result *= numbers[i+1]
		case "||":
			result = concatenate(result, numbers[i+1])
		}
	}
	return result
}

func concatenate(a, b int) int {
	aStr := strconv.Itoa(a)
	bStr := strconv.Itoa(b)
	concatStr := aStr + bStr
	concat, _ := strconv.Atoi(concatStr)
	return concat
}

func generateOperatorCombinations(numOperators int, operators []string) [][]string {
	var result [][]string
	result = append(result, []string{})

	for i := 0; i < numOperators; i++ {
		var newResult [][]string
		for _, combination := range result {
			for _, op := range operators {
				newCombination := append([]string(nil), combination...)
				newCombination = append(newCombination, op)
				newResult = append(newResult, newCombination)
			}
		}
		result = newResult
	}

	return result
}

func isValidEquation(testValue int, numbers []int, operators []string) bool {
	operatorCombinations := generateOperatorCombinations(len(numbers)-1, operators)

	for _, ops := range operatorCombinations {
		result := evaluateExpression(numbers, ops)
		if result == testValue {
			return true
		}
	}
	return false
}

func getInput(filename string) ([][]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result [][]int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line format: %s", line)
		}

		testValue, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid number: %s", parts[0])
		}

		numbersStr := strings.Fields(parts[1])
		var numbers []int
		for _, val := range numbersStr {
			num, err := strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("invalid number: %s", val)
			}
			numbers = append(numbers, num)
		}

		result = append(result, append([]int{testValue}, numbers...))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func main() {
	slices, err := getInput("input.txt")
	if err != nil {
		panic(err)
	}

	part1Operators := []string{"+", "*"}
	part2Operators := []string{"+", "*", "||"}

	part1Total := 0
	part2Total := 0

	for _, slice := range slices {
		testValue := slice[0]
		numbers := slice[1:]

		if isValidEquation(testValue, numbers, part1Operators) {
			part1Total += testValue
		}

		if isValidEquation(testValue, numbers, part2Operators) {
			part2Total += testValue
		}
	}

	fmt.Println("[PART 1] Total Calibration Result:", part1Total)
	fmt.Println("[PART 2] Total Calibration Result:", part2Total)
}
