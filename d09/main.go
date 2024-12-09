package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func getInput(filename string) ([]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil, fmt.Errorf("file is empty or could not read line")
	}

	line := scanner.Text()
	var result []int
	for _, char := range line {
		digit, err := strconv.Atoi(string(char))
		if err != nil {
			return nil, fmt.Errorf("invalid digit '%c': %w", char, err)
		}
		result = append(result, digit)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error while reading file: %w", err)
	}

	return result, nil
}

func createLongFormat(digits []int) []rune {
	var blocks []rune
	fileID := 0
	for i, d := range digits {
		if i%2 == 0 {
			// Append file blocks with current file ID
			for j := 0; j < d; j++ {
				blocks = append(blocks, rune('0'+fileID))
			}
			fileID++
		} else {
			// Append free space blocks
			for j := 0; j < d; j++ {
				blocks = append(blocks, '.')
			}
		}
	}
	return blocks
}

func compressPart1(input []rune) []rune {
	for i := len(input) - 1; i >= 0; i-- {
		if input[i] != '.' {
			// Find the leftmost free space
			for j := 0; j < i; j++ {
				if input[j] == '.' {
					// Move the block
					input[j] = input[i]
					input[i] = '.'
					break
				}
			}
		}
	}
	return input
}

func calcCheckSum(input []rune) int {
	sum := 0
	for i, val := range input {
		if val != '.' {
			fileID := int(val - '0')
			sum += fileID * i
		}
	}
	return sum
}

func compressPart2(input []rune) []rune {
	// Find the highest file ID
	maxID := rune('0')
	for _, r := range input {
		if r != '.' && r > maxID {
			maxID = r
		}
	}

	// Iterate through file IDs in descending order
	for id := maxID; id >= '0'; id-- {
		// Find the file
		start, end := -1, -1
		for i, r := range input {
			if r == id {
				if start == -1 {
					start = i
				}
				end = i
			} else if start != -1 {
				break
			}
		}

		if start == -1 {
			continue
		}

		fileSize := end - start + 1

		// Find the leftmost free space that can fit the file
		freeStart := -1
		freeSize := 0
		for i := 0; i < start; i++ { // Only look at space to the left of the file
			if input[i] == '.' {
				if freeStart == -1 {
					freeStart = i
				}
				freeSize++
				if freeSize == fileSize {
					break
				}
			} else {
				freeStart = -1
				freeSize = 0
			}
		}

		// Move the file if a suitable free space was found
		if freeSize == fileSize {
			copy(input[freeStart:freeStart+fileSize], input[start:end+1])
			for i := start; i <= end; i++ {
				input[i] = '.'
			}
		}
	}

	return input
}

func main() {
	integers, err := getInput("input.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Part 1: Compress by moving individual blocks
	lfPart1 := createLongFormat(integers)
	compressedPart1 := compressPart1(lfPart1)
	checkSumPart1 := calcCheckSum(compressedPart1)
	fmt.Println("[PART 1]: The checksum is:", checkSumPart1)

	// Part 2: Compress by moving whole files
	lfPart2 := createLongFormat(integers)
	compressedPart2 := compressPart2(lfPart2)
	checkSumPart2 := calcCheckSum(compressedPart2)
	fmt.Println("[PART 2]: The checksum is:", checkSumPart2)
}
