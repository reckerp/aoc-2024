package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Position struct {
	row, col int
}

var (
	directions = map[rune]Position{
		'^': {-1, 0},
		'>': {0, 1},
		'v': {1, 0},
		'<': {0, -1},
	}
)

func main() {
	gridP1, instructions := getInput("input.txt")
	gridP2 := make([][]rune, len(gridP1))
	for i := range gridP1 {
		gridP2[i] = make([]rune, len(gridP1[i]))
		copy(gridP2[i], gridP1[i])
	}

	resultP1 := solvePart1(gridP1, instructions)
	fmt.Printf("[PART 1] GPS sum: %d\n", resultP1)

	resultP2 := solvePart2(gridP2, instructions)
	fmt.Printf("[PART 2] GPS sum: %d\n", resultP2)
}

func getInput(filename string) ([][]rune, string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var grid [][]rune
	var instructions string
	readingGrid := true

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			readingGrid = false
			continue
		}

		if readingGrid {
			grid = append(grid, []rune(line))
		} else {
			instructions += strings.TrimSpace(line)
		}
	}

	return grid, instructions
}

func findRobot(grid [][]rune) Position {
	for row := range grid {
		for col := range grid[row] {
			if grid[row][col] == '@' {
				grid[row][col] = '.'
				return Position{row, col}
			}
		}
	}
	panic("Robot not found in the grid")
}

func solvePart1(grid [][]rune, instructions string) int {
	robotPos := findRobot(grid)
	moveRobot(grid, robotPos, instructions, true)
	return calculateGPSSum(grid, 'O')
}

func solvePart2(grid [][]rune, instructions string) int {
	expandedGrid := expandGrid(grid)
	robotPos := findRobot(expandedGrid)
	moveRobot(expandedGrid, robotPos, instructions, false)
	return calculateGPSSum(expandedGrid, ']')
}

func moveRobot(grid [][]rune, startPos Position, instructions string, part1 bool) {
	currentPos := startPos
	for _, instruction := range instructions {
		dir := directions[instruction]
		nextPos := Position{currentPos.row + dir.row, currentPos.col + dir.col}

		if grid[nextPos.row][nextPos.col] == '#' {
			continue
		} else if grid[nextPos.row][nextPos.col] == '.' {
			currentPos = nextPos
		} else if (part1 && grid[nextPos.row][nextPos.col] == 'O') ||
			(!part1 && (grid[nextPos.row][nextPos.col] == '[' || grid[nextPos.row][nextPos.col] == ']')) {
			if pushBoxes(grid, currentPos, dir, part1) {
				currentPos = nextPos
			}
		}
	}
}

func pushBoxes(grid [][]rune, robotPos Position, dir Position, part1 bool) bool {
	queue := []Position{robotPos}
	seen := make(map[Position]bool)

	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]

		if seen[pos] {
			continue
		}
		seen[pos] = true

		nextPos := Position{pos.row + dir.row, pos.col + dir.col}
		if grid[nextPos.row][nextPos.col] == '#' {
			return false // Stop if blocked by an obstacle
		}

		if part1 {
			// For part1, handle boxes ('O') and add them to the queue
			if grid[nextPos.row][nextPos.col] == 'O' {
				queue = append(queue, nextPos)
			}
		} else {
			// For part2, handle special cases '[' and ']' only
			switch grid[nextPos.row][nextPos.col] {
			case '[':
				queue = append(queue, nextPos, Position{nextPos.row, nextPos.col + 1})
			case ']':
				queue = append(queue, nextPos, Position{nextPos.row, nextPos.col - 1})
			}
		}
	}

	// Move the boxes after the traversal is done
	moveBoxes(grid, seen, dir)

	return true
}

func moveBoxes(grid [][]rune, boxPositions map[Position]bool, dir Position) {
	for len(boxPositions) > 0 {
		var toMove Position
		for pos := range boxPositions {
			nextPos := Position{pos.row + dir.row, pos.col + dir.col}
			if !boxPositions[nextPos] {
				// Move the box to the new position and clear the old position
				grid[nextPos.row][nextPos.col] = grid[pos.row][pos.col]
				grid[pos.row][pos.col] = '.'
				toMove = pos
				break
			}
		}
		delete(boxPositions, toMove)
	}
}

func calculateGPSSum(grid [][]rune, boxRune rune) int {
	sum := 0
	for row := range grid {
		for col := range grid[row] {
			if grid[row][col] == boxRune {
				sum += 100*row + col
			}
		}
	}
	return sum
}

func expandGrid(grid [][]rune) [][]rune {
	rows, cols := len(grid), len(grid[0])
	expandedGrid := make([][]rune, rows)
	for r := range grid {
		expandedGrid[r] = make([]rune, cols*2)
		for c, ch := range grid[r] {
			switch ch {
			case '#':
				expandedGrid[r][c*2], expandedGrid[r][c*2+1] = '#', '#'
			case 'O':
				expandedGrid[r][c*2], expandedGrid[r][c*2+1] = '[', ']'
			case '.':
				expandedGrid[r][c*2], expandedGrid[r][c*2+1] = '.', '.'
			case '@':
				expandedGrid[r][c*2], expandedGrid[r][c*2+1] = '@', '.'
			}
		}
	}
	return expandedGrid
}
