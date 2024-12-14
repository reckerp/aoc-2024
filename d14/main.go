package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	FIELD_WIDTH  = 101
	FIELD_HEIGHT = 103
)

type Robot struct {
	X, Y, VX, VY int
}

func getInput() ([]Robot, error) {
	file, err := os.Open("input.txt")
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var robots []Robot
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line format: %s", line)
		}

		positionPart := strings.TrimPrefix(parts[0], "p=")
		velocityPart := strings.TrimPrefix(parts[1], "v=")

		position := strings.Split(positionPart, ",")
		velocity := strings.Split(velocityPart, ",")

		x, _ := strconv.Atoi(position[0])
		y, _ := strconv.Atoi(position[1])
		vx, _ := strconv.Atoi(velocity[0])
		vy, _ := strconv.Atoi(velocity[1])

		robots = append(robots, Robot{X: x, Y: y, VX: vx, VY: vy})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return robots, nil
}

func moveRobot(robot Robot) Robot {
	return Robot{
		X:  (robot.X + robot.VX + FIELD_WIDTH) % FIELD_WIDTH,
		Y:  (robot.Y + robot.VY + FIELD_HEIGHT) % FIELD_HEIGHT,
		VX: robot.VX,
		VY: robot.VY,
	}
}

func simulateRobotIterations(robots []Robot, times int) []Robot {
	for i := 0; i < times; i++ {
		for j := range robots {
			robots[j] = moveRobot(robots[j])
		}
	}
	return robots
}

func calcSecurityLevel(robots []Robot) int {
	quadrants := make([]int, 4)
	for _, robot := range robots {
		if robot.X != FIELD_WIDTH/2 && robot.Y != FIELD_HEIGHT/2 {
			if robot.X < FIELD_WIDTH/2 && robot.Y < FIELD_HEIGHT/2 {
				quadrants[0]++
			} else if robot.X > FIELD_WIDTH/2 && robot.Y < FIELD_HEIGHT/2 {
				quadrants[1]++
			} else if robot.X < FIELD_WIDTH/2 && robot.Y > FIELD_HEIGHT/2 {
				quadrants[2]++
			} else if robot.X > FIELD_WIDTH/2 && robot.Y > FIELD_HEIGHT/2 {
				quadrants[3]++
			}
		}
	}
	return quadrants[0] * quadrants[1] * quadrants[2] * quadrants[3]
}

func robotDensity(robots []Robot) float64 {
	sum := 0.0
	count := 0
	for i := 0; i < len(robots)-1; i++ {
		for j := i + 1; j < len(robots); j++ {
			dx := float64(robots[i].X - robots[j].X)
			dy := float64(robots[i].Y - robots[j].Y)
			sum += math.Sqrt(dx*dx + dy*dy)
			count++
		}
	}
	return sum / float64(count)
}

func part1(robots []Robot) int {
	simulatedRobots := simulateRobotIterations(robots, 100)
	return calcSecurityLevel(simulatedRobots)
}

func part2(robots []Robot) int {
	minDensity := math.Inf(1)
	minTime := 0

	for t := 0; t < 20000; t++ {
		density := robotDensity(robots)
		if density < minDensity {
			minDensity = density
			minTime = t
		}
		robots = simulateRobotIterations(robots, 1) // Simulate one step at a time
	}

	return minTime
}

func main() {
	part1Input, err := getInput()
	if err != nil {
		panic(err)
	}

	part2Input := make([]Robot, len(part1Input))
	copy(part2Input, part1Input)
	fmt.Println("[PART 1] Security level after 100 iterations:", part1(part1Input))
	fmt.Printf("[PART 2] The Easter egg appears after %d seconds\n", part2(part2Input))
}
