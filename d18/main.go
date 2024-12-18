package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x, y int
}

type State struct {
	point    Point
	steps    int
	priority int
}

type PriorityQueue []*State

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].priority < pq[j].priority }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*State))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func getInput(filename string) map[Point]bool {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	corruptedSpaces := make(map[Point]bool)
	scanner := bufio.NewScanner(file)
	count := 0

	for scanner.Scan() && count < 1024 {
		line := scanner.Text()
		if line == "" {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			continue
		}

		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])

		corruptedSpaces[Point{x: x, y: y}] = true
		count++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	return corruptedSpaces
}

func findShortestPath(corruptedSpaces map[Point]bool, gridSize int) int {
	start := Point{x: 0, y: 0}
	end := Point{x: gridSize, y: gridSize}

	directions := []Point{{x: 1, y: 0}, {x: -1, y: 0}, {x: 0, y: 1}, {x: 0, y: -1}}

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	heap.Push(&pq, &State{
		point:    start,
		steps:    0,
		priority: manhattanDistance(start, end),
	})

	visited := make(map[Point]bool)

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*State)

		if current.point == end {
			return current.steps
		}

		if visited[current.point] {
			continue
		}
		visited[current.point] = true

		for _, dir := range directions {
			next := Point{
				x: current.point.x + dir.x,
				y: current.point.y + dir.y,
			}

			if next.x < 0 || next.x > gridSize || next.y < 0 || next.y > gridSize {
				continue
			}

			if corruptedSpaces[next] {
				continue
			}

			if !visited[next] {
				heap.Push(&pq, &State{
					point:    next,
					steps:    current.steps + 1,
					priority: current.steps + 1 + manhattanDistance(next, end),
				})
			}
		}
	}

	return -1 // No path found
}

func findFirstBlockingByte(coordinates []Point, gridSize int) Point {
	corruptedSpaces := make(map[Point]bool)

	for _, coord := range coordinates {
		// Add current coordinate to corrupted spaces
		corruptedSpaces[coord] = true

		// Check if path is blocked
		if findShortestPath(corruptedSpaces, gridSize) == -1 {
			return coord
		}
	}

	// This should not happen based on problem description
	return Point{x: -1, y: -1}
}

func readAllCoordinates(filename string) []Point {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	var coordinates []Point
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			continue
		}

		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])

		coordinates = append(coordinates, Point{x: x, y: y})
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	return coordinates
}

func manhattanDistance(a, b Point) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	gridSize := 70

	input := getInput("input.txt")
	shortestPath := findShortestPath(input, gridSize)
	fmt.Printf("[PART 1]: Shortest Path: %d steps\n", shortestPath)

	coordinates := readAllCoordinates("input.txt")
	blockingPoint := findFirstBlockingByte(coordinates, gridSize)
	fmt.Printf("[PART 2]: First blocking point: %d,%d\n", blockingPoint.x, blockingPoint.y)
}
