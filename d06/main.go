package main

import (
	"bufio"
	"fmt"
	"os"
)

type Position struct {
	row, col int
}

type State struct {
	pos       Position
	direction int
}

var directions = [][2]int{
	{-1, 0}, // Up
	{0, 1},  // Right
	{1, 0},  // Down
	{0, -1}, // Left
}

func nextDirection(currentIndex int) int {
	return (currentIndex + 1) % len(directions)
}

func main() {
	matrix := readInput()

	// Part 1
	sumPositions := distinctGuardPositions(matrix)
	fmt.Println("[Part 1] Sum of distinct guard positions:", sumPositions)

	// Part 2
	matrix = readInput()
	loopInducingCount := findLoopInducingObstructions(matrix)
	fmt.Println("[Part 2] Number of positions causing a loop:", loopInducingCount)
}

func readInput() [][]string {
	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var matrix [][]string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		row := []string{}
		for _, char := range line {
			row = append(row, string(char))
		}
		matrix = append(matrix, row)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return matrix
}

func distinctGuardPositions(input [][]string) int {
	x, y, currentDirection := getGuardStartCoords(input)
	input[y][x] = "X"
	sumPositions := 1

	for {
		posx, posy := x+directions[currentDirection][1], y+directions[currentDirection][0]
		if posy < 0 || posy >= len(input) || posx < 0 || posx >= len(input[0]) {
			break
		}

		if input[posy][posx] == "#" {
			currentDirection = nextDirection(currentDirection)
			continue
		} else if input[posy][posx] != "X" {
			sumPositions++
			input[posy][posx] = "X"
		}
		x, y = posx, posy
	}
	return sumPositions
}

func getGuardStartCoords(input [][]string) (int, int, int) {
	directionsMap := map[string]int{
		"^": 0, // Up
		">": 1, // Right
		"v": 2, // Down
		"<": 3, // Left
	}
	for y, row := range input {
		for x, char := range row {
			if dir, ok := directionsMap[char]; ok {
				return x, y, dir
			}
		}
	}
	return -1, -1, -1
}

func duplicateMatrix(matrix [][]string) [][]string {
	duplicate := make([][]string, len(matrix))
	for i := range matrix {
		duplicate[i] = make([]string, len(matrix[i]))
		copy(duplicate[i], matrix[i])
	}
	return duplicate
}

func findLoopInducingObstructions(input [][]string) int {
	startX, startY, startDirection := getGuardStartCoords(input)
	validObstructions := 0

	for y := range input {
		for x := range input[y] {
			if input[y][x] != "." || (x == startX && y == startY) {
				continue
			}

			// Simulate obstruction placement
			tempMatrix := duplicateMatrix(input)
			tempMatrix[y][x] = "#"

			visitedStates := make(map[State]bool)
			initialState := State{pos: Position{startY, startX}, direction: startDirection}

			if createsLoop(tempMatrix, initialState, visitedStates) {
				validObstructions++
			}
		}
	}

	return validObstructions
}

func createsLoop(matrix [][]string, initialState State, visitedStates map[State]bool) bool {
	currentState := initialState

	for {
		// Check if we've visited this state before
		if visitedStates[currentState] {
			return true
		}

		// Mark the current state as visited
		visitedStates[currentState] = true

		// Calculate the next position
		nextRow := currentState.pos.row + directions[currentState.direction][0]
		nextCol := currentState.pos.col + directions[currentState.direction][1]

		// Check bounds
		if nextRow < 0 || nextRow >= len(matrix) || nextCol < 0 || nextCol >= len(matrix[0]) {
			return false
		}

		// If the next position is an obstacle, change direction
		if matrix[nextRow][nextCol] == "#" {
			currentState.direction = nextDirection(currentState.direction)
			continue
		}

		// Move to the next position
		currentState.pos = Position{nextRow, nextCol}
	}
}
