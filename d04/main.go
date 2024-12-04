package main

import (
	"bufio"
	"fmt"
	"os"
)

type Direction struct {
	dx, dy int
}

func main() {
	part1()
	part2()
}

func part1() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var grid []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		grid = append(grid, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Define the word to search
	word := "XMAS"

	// Define all 8 directions
	directions := []Direction{
		{0, 1},   // Right
		{0, -1},  // Left
		{1, 0},   // Down
		{-1, 0},  // Up
		{1, 1},   // Diagonal Down-Right
		{-1, -1}, // Diagonal Up-Left
		{1, -1},  // Diagonal Down-Left
		{-1, 1},  // Diagonal Up-Right
	}

	// Function to check if a position is valid
	isValid := func(x, y int) bool {
		return x >= 0 && x < len(grid) && y >= 0 && y < len(grid[0])
	}

	// Function to search for the word starting at (x, y) in a given direction
	search := func(x, y int, direction Direction) bool {
		for i := 0; i < len(word); i++ {
			nx, ny := x+i*direction.dx, y+i*direction.dy
			if !isValid(nx, ny) || grid[nx][ny] != word[i] {
				return false
			}
		}
		return true
	}

	// Find all occurrences of the word
	var occurrences []string
	for x := 0; x < len(grid); x++ {
		for y := 0; y < len(grid[0]); y++ {
			for _, direction := range directions {
				if search(x, y, direction) {
					occurrences = append(occurrences, fmt.Sprintf("Start at (%d, %d), Direction: (%d, %d)", x, y, direction.dx, direction.dy))
				}
			}
		}
	}

	fmt.Printf("[PART1] Number of '%s': %d\n", word, len(occurrences))
}

func part2() {
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var grid []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		grid = append(grid, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	validPatterns := []string{"MAS", "SAM"}

	// Function to check if a position is valid
	isValid := func(x, y int) bool {
		return x >= 0 && x < len(grid) && y >= 0 && y < len(grid[0])
	}

	// Function to check for "MAS" or "SAM" in a diagonal direction
	checkDiagonal := func(startX, startY, dx, dy int) bool {
		word := ""
		for i := 0; i < 3; i++ {
			nx, ny := startX+i*dx, startY+i*dy
			if !isValid(nx, ny) {
				return false
			}
			word += string(grid[nx][ny])
		}
		for _, pattern := range validPatterns {
			if word == pattern {
				return true
			}
		}
		return false
	}

	// Find all X-MAS patterns
	count := 0
	for x := 1; x < len(grid)-1; x++ {
		for y := 1; y < len(grid[0])-1; y++ {
			// Check diagonals for X-MAS pattern
			if checkDiagonal(x-1, y-1, 1, 1) && checkDiagonal(x-1, y+1, 1, -1) {
				count++
			}
		}
	}

	fmt.Printf("[PART2] Number of X-MAS patterns: %d\n", count)
}
