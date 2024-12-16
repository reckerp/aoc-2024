package main

import (
	"container/heap"
	"fmt"
	"os"
	"strings"
)

type Point struct{ x, y int }

const (
	TurnCost    = 1000
	MoveCost    = 1
	StartMarker = 'S'
	EndMarker   = 'E'
	WallMarker  = '#'
)

// Directions: RIGHT, DOWN, LEFT, UP
var Directions = []Point{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

// State represents the position, direction, and score in the grid.
type State struct {
	score, x, y, dir int
}

// PriorityQueue for BFS with scores
type PriorityQueue []State

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].score < pq[j].score }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(State))
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func main() {
	grid, start, end := getInput("input.txt")
	part1, part2 := solve(grid, start, end)

	fmt.Println("[PART 1] The total cost is:", part1)
	fmt.Println("[PART 2] Total tiles:", part2)
}

func getInput(filename string) ([][]rune, Point, Point) {
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("Failed to read file: %v", err))
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	grid := make([][]rune, len(lines))

	var start, end Point
	for y, line := range lines {
		grid[y] = []rune(line)
		for x, char := range line {
			switch char {
			case StartMarker:
				start = Point{x, y}
			case EndMarker:
				end = Point{x, y}
			}
		}
	}
	return grid, start, end
}

func solve(grid [][]rune, start, end Point) (int, int) {
	width, height := len(grid[0]), len(grid)
	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, State{score: 0, x: start.x, y: start.y, dir: 0})

	seen := make(map[[3]int]bool) // State: (x, y, dir) -> visited
	origins := make(map[State][]State)
	var finalCost int

	// BFS loop to find the shortest path
	for pq.Len() > 0 {
		curr := heap.Pop(pq).(State)

		// Stop BFS when we reach the end
		if curr.x == end.x && curr.y == end.y {
			finalCost = curr.score
			break
		}

		// Skip if this state has already been visited
		key := [3]int{curr.x, curr.y, curr.dir}
		if seen[key] {
			continue
		}
		seen[key] = true

		// Turn to new directions
		for i := 1; i <= 3; i++ {
			newDir := (curr.dir + i) % 4
			heap.Push(pq, State{
				score: curr.score + TurnCost,
				x:     curr.x,
				y:     curr.y,
				dir:   newDir,
			})
			origins[State{curr.score + TurnCost, curr.x, curr.y, newDir}] = append(origins[State{curr.score + TurnCost, curr.x, curr.y, newDir}], curr)
		}

		// Move forward in the current direction
		dx, dy := Directions[curr.dir].x, Directions[curr.dir].y
		nx, ny := curr.x+dx, curr.y+dy
		if nx >= 0 && nx < width && ny >= 0 && ny < height && grid[ny][nx] != WallMarker {
			heap.Push(pq, State{
				score: curr.score + MoveCost,
				x:     nx,
				y:     ny,
				dir:   curr.dir,
			})
			origins[State{curr.score + MoveCost, nx, ny, curr.dir}] = append(origins[State{curr.score + MoveCost, nx, ny, curr.dir}], curr)
		}
	}

	// Trace back to count all unique tiles visited
	return finalCost, countVisitedTiles(origins, end, finalCost)
}

// countVisitedTiles traces back from the end point to count all unique tiles.
func countVisitedTiles(origins map[State][]State, end Point, finalCost int) int {
	visited := make(map[Point]bool)
	queue := []State{}

	// Start tracing back from all possible directions at the end point
	for d := 0; d < 4; d++ {
		queue = append(queue, State{score: finalCost, x: end.x, y: end.y, dir: d})
	}

	for len(queue) > 0 {
		curr := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		visited[Point{curr.x, curr.y}] = true
		for _, prev := range origins[curr] {
			queue = append(queue, prev)
		}
	}

	return len(visited)
}
