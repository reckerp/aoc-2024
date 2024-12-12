package main

import (
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	X, Y int
}

var directions = []Point{
	{0, 1}, {1, 0}, {0, -1}, {-1, 0},
}

func main() {
	grid := getInput("input.txt")

	totalPricePart1 := calculateTotalPrice(grid, calculatePart1Price)
	fmt.Println("[PART 1] The total price is:", totalPricePart1)

	totalPricePart2 := calculateTotalPrice(grid, calculatePart2Price)
	fmt.Println("[PART 2] The total price is:", totalPricePart2)
}

func getInput(filename string) [][]rune {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var grid [][]rune
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		grid = append(grid, []rune(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return grid
}

// calculateTotalPrice calculates the total price using pricing strategy
func calculateTotalPrice(grid [][]rune, pricingFunc func([][]rune, Point, map[Point]bool) int) int {
	visited := make(map[Point]bool)
	totalPrice := 0

	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[0]); x++ {
			point := Point{x, y}
			if !visited[point] {
				totalPrice += pricingFunc(grid, point, visited)
			}
		}
	}
	return totalPrice
}

// calculatePart1Price calculates price for Part 1 (area * perimeter)
func calculatePart1Price(grid [][]rune, start Point, visited map[Point]bool) int {
	area, perimeter := exploreRegion(grid, start, visited)
	return area * perimeter
}

// calculatePart2Price calculates price for Part 2 (region size * region perimeter)
func calculatePart2Price(grid [][]rune, start Point, visited map[Point]bool) int {
	_, region := findContiguousRegion(grid, start, visited)
	return len(region) * calculateRegionPerimeter(region)
}

// exploreRegion explores a region and calculates its area and perimeter
func exploreRegion(grid [][]rune, start Point, visited map[Point]bool) (int, int) {
	queue := []Point{start}
	typeRune := grid[start.Y][start.X]
	area, perimeter := 0, 0

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}
		visited[current] = true
		area++

		for _, dir := range directions {
			neighbor := Point{current.X + dir.X, current.Y + dir.Y}
			if isInBounds(grid, neighbor) {
				if grid[neighbor.Y][neighbor.X] == typeRune {
					if !visited[neighbor] {
						queue = append(queue, neighbor)
					}
				} else {
					perimeter++
				}
			} else {
				perimeter++
			}
		}
	}
	return area, perimeter
}

// findContiguousRegion finds a contiguous region of the same type
func findContiguousRegion(grid [][]rune, start Point, visited map[Point]bool) (rune, map[Point]bool) {
	cellType := grid[start.Y][start.X]
	queue := []Point{start}
	region := make(map[Point]bool)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}
		visited[current] = true
		region[current] = true

		for _, dir := range directions {
			neighbor := Point{current.X + dir.X, current.Y + dir.Y}
			if isInBounds(grid, neighbor) &&
				grid[neighbor.Y][neighbor.X] == cellType &&
				!visited[neighbor] {
				queue = append(queue, neighbor)
			}
		}
	}

	return cellType, region
}

func calculateRegionPerimeter(region map[Point]bool) int {
	perimeter := 0
	for point := range region {
		x, y := point.X, point.Y
		checks := []struct{ nx, ny, x1, y1, x2, y2 int }{
			{x + 1, y, x, y - 1, x + 1, y - 1},
			{x - 1, y, x, y - 1, x - 1, y - 1},
			{x, y + 1, x - 1, y, x - 1, y + 1},
			{x, y - 1, x - 1, y, x - 1, y - 1},
		}

		for _, check := range checks {
			neighbor := Point{check.nx, check.ny}
			corner1 := Point{check.x1, check.y1}
			corner2 := Point{check.x2, check.y2}

			if !region[neighbor] && !(region[corner1] && !region[corner2]) {
				perimeter++
			}
		}
	}
	return perimeter
}

func isInBounds(grid [][]rune, p Point) bool {
	return p.Y >= 0 && p.Y < len(grid) && p.X >= 0 && p.X < len(grid[0])
}
