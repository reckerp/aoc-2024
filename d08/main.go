package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

type Point struct {
	row, col int
}

func (p Point) GetDelta(other Point) Point {
	return Point{row: other.row - p.row, col: other.col - p.col}
}

type AntennaPositions map[string][]Point

func findAntennas(input [][]string) AntennaPositions {
	positions := make(AntennaPositions)

	for rowIdx, row := range input {
		for colIdx, cell := range row {
			if cell != "." {
				positions[cell] = append(positions[cell], Point{row: rowIdx, col: colIdx})
			}
		}
	}

	return positions
}

func isValid(input [][]string, p Point) bool {
	return p.row >= 0 && p.row < len(input) && p.col >= 0 && p.col < len(input[0])
}

func generateCombinations(points []Point, size int) [][]Point {
	if size == 0 {
		return [][]Point{{}}
	}
	if len(points) == 0 {
		return nil
	}

	var combinations [][]Point
	head := points[0]
	tail := points[1:]

	for _, combo := range generateCombinations(tail, size-1) {
		combinations = append(combinations, append([]Point{head}, combo...))
	}
	combinations = append(combinations, generateCombinations(tail, size)...)

	return combinations
}

func calcMaxSteps(rows, cols, dx, dy int) int {
	maxDim := max(rows, cols)
	maxDelta := max(abs(dx), abs(dy))
	if maxDelta == 0 {
		return 1
	}
	return int(math.Ceil(float64(maxDim) / float64(maxDelta)))
}

func getAntinodes(input [][]string, p1, p2 Point, maxDistance bool) []Point {
	delta := p1.GetDelta(p2)
	var antinodes []Point

	if maxDistance {
		candidates := []Point{
			{row: p1.row - delta.row, col: p1.col - delta.col},
			{row: p2.row + delta.row, col: p2.col + delta.col},
		}
		for _, p := range candidates {
			if isValid(input, p) {
				antinodes = append(antinodes, p)
			}
		}
	} else {
		maxSteps := calcMaxSteps(len(input), len(input[0]), delta.row, delta.col)
		for i := -maxSteps; i <= maxSteps; i++ {
			candidates := []Point{
				{row: p1.row - i*delta.row, col: p1.col - i*delta.col},
				{row: p2.row + i*delta.row, col: p2.col + i*delta.col},
			}
			for _, p := range candidates {
				if isValid(input, p) {
					antinodes = append(antinodes, p)
				}
			}
		}
	}

	return antinodes
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func getAllAntinodes(input [][]string, maxDistance bool) map[Point]bool {
	antennaPositions := findAntennas(input)
	uniquePoints := make(map[Point]bool)

	for _, positions := range antennaPositions {
		if len(positions) <= 1 {
			continue
		}

		for _, combo := range generateCombinations(positions, 2) {
			antinodes := getAntinodes(input, combo[0], combo[1], maxDistance)
			for _, p := range antinodes {
				uniquePoints[p] = true
			}
		}
	}

	return uniquePoints
}

func getInput(filename string) ([][]string, error) {
	var input [][]string

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		input = append(input, strings.Split(line, ""))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return input, nil
}

func main() {
	input, err := getInput("input.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	part1 := getAllAntinodes(input, true)
	part2 := getAllAntinodes(input, false)

	fmt.Println("[PART1] Number of Antinodes:", len(part1))
	fmt.Println("[PART2] Number of Antinodes:", len(part2))
}
