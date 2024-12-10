package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Direction struct {
	dx, dy int
}

var directions = []Direction{
	{0, -1}, {0, 1}, {-1, 0}, {1, 0},
}

func getInput() [][]int {
	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var map2D [][]int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]int, len(line))
		for i, ch := range line {
			row[i], _ = strconv.Atoi(string(ch))
		}
		map2D = append(map2D, row)
	}
	return map2D
}

func isValid(x, y int, m [][]int) bool {
	return x >= 0 && x < len(m[0]) && y >= 0 && y < len(m)
}

func dfs(x, y, height int, m [][]int, visited [][]bool) int {
	if !isValid(x, y, m) || visited[y][x] || m[y][x] != height {
		return 0
	}

	visited[y][x] = true
	if height == 9 {
		return 1
	}

	count := 0
	for _, dir := range directions {
		count += dfs(x+dir.dx, y+dir.dy, height+1, m, visited)
	}
	return count
}

func calculateTrailheadScores(m [][]int) int {
	totalScore := 0
	visited := make([][]bool, len(m))
	for i := range visited {
		visited[i] = make([]bool, len(m[0]))
	}

	for y := 0; y < len(m); y++ {
		for x := 0; x < len(m[0]); x++ {
			if m[y][x] == 0 {
				score := dfs(x, y, 0, m, visited)
				totalScore += score
				for i := range visited {
					for j := range visited[i] {
						visited[i][j] = false
					}
				}
			}
		}
	}
	return totalScore
}

func dfsCount(x, y, height int, m [][]int) int {
	if !isValid(x, y, m) || m[y][x] != height {
		return 0
	}

	if height == 9 {
		return 1
	}

	count := 0
	for _, dir := range directions {
		count += dfsCount(x+dir.dx, y+dir.dy, height+1, m)
	}
	return count
}

func calculateTrailheadRatings(m [][]int) int {
	totalRating := 0

	for y := 0; y < len(m); y++ {
		for x := 0; x < len(m[0]); x++ {
			if m[y][x] == 0 {
				rating := dfsCount(x, y, 0, m)
				totalRating += rating
			}
		}
	}
	return totalRating
}

func main() {
	input := getInput()
	totalScore := calculateTrailheadScores(input)
	fmt.Printf("[PART 1] Sum of scores of all trailheads: %d\n", totalScore)

	totalRating := calculateTrailheadRatings(input)
	fmt.Printf("[PART 2] Total rating of trailheads: %d\n", totalRating)
}
