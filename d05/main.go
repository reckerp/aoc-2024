package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Input struct {
	Rules   map[int]map[int]bool
	Updates [][]int
}

func main() {
	input, err := readInput("input.txt")
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	validUpdates, invalidUpdates := validateUpdates(input.Rules, input.Updates)

	sumMedian := 0
	for _, update := range validUpdates {
		sumMedian += findMedian(update)
	}

	fmt.Println("[PART1] Sum of medians:", sumMedian)

	fixedUpdates := fixInvalid(invalidUpdates, input.Rules)

	sumMedian = 0
	for _, update := range fixedUpdates {
		sumMedian += findMedian(update)
	}

	fmt.Println("[PART2] Sum of medians:", sumMedian)
}

func readInput(filename string) (Input, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Input{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	input := Input{
		Rules: make(map[int]map[int]bool),
	}

	// Read rules
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			continue
		}
		a, _ := strconv.Atoi(parts[0])
		b, _ := strconv.Atoi(parts[1])
		if input.Rules[a] == nil {
			input.Rules[a] = make(map[int]bool)
		}
		input.Rules[a][b] = true
	}

	// Read updates
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		update := make([]int, len(parts))
		for i, part := range parts {
			update[i], _ = strconv.Atoi(part)
		}
		input.Updates = append(input.Updates, update)
	}

	return input, scanner.Err()
}

func validateUpdates(rules map[int]map[int]bool, updates [][]int) ([][]int, [][]int) {
	validUpdates := [][]int{}
	invalidUpdates := [][]int{}
	for _, update := range updates {
		if isValidUpdate(update, rules) {
			validUpdates = append(validUpdates, update)
		} else {
			invalidUpdates = append(invalidUpdates, update)
		}
	}
	return validUpdates, invalidUpdates
}

func isValidUpdate(update []int, rules map[int]map[int]bool) bool {
	for i := 0; i < len(update); i++ {
		for j := i + 1; j < len(update); j++ {
			a, b := update[i], update[j]
			if rules[b][a] {
				return false
			}
		}
	}
	return true
}

func findMedian(arr []int) int {
	if len(arr) == 0 {
		return 0
	}
	if len(arr)%2 == 1 {
		return arr[len(arr)/2]
	}
	return (arr[len(arr)/2-1] + arr[len(arr)/2]) / 2
}

func fixInvalid(invalidUpdates [][]int, rules map[int]map[int]bool) [][]int {
	fixedUpdates := make([][]int, len(invalidUpdates))
	for i, update := range invalidUpdates {
		fixedUpdates[i] = fixUpdate(update, rules)
	}
	return fixedUpdates
}

func fixUpdate(update []int, rules map[int]map[int]bool) []int {
	fixed := make([]int, len(update))
	copy(fixed, update)

	for i := 0; i < len(fixed); i++ {
		for j := i + 1; j < len(fixed); j++ {
			if rules[fixed[j]][fixed[i]] {
				// Swap elements if they violate a rule
				fixed[i], fixed[j] = fixed[j], fixed[i]
				// Start over from the beginning after a swap
				i = -1
				break
			}
		}
	}

	return fixed
}
