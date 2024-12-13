package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Coordinate struct {
	X int
	Y int
}

type ClawMachine struct {
	ButtonA Coordinate
	ButtonB Coordinate
	Prize   Coordinate
}

func getInput(filename string) ([]ClawMachine, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var machines []ClawMachine
	scanner := bufio.NewScanner(file)
	var currentMachine ClawMachine

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if currentMachine != (ClawMachine{}) {
				machines = append(machines, currentMachine)
				currentMachine = ClawMachine{}
			}
			continue
		}

		parts := strings.Split(line, ": ")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line format: %s", line)
		}

		key := parts[0]
		value := parts[1]

		coord, err := parseCoordinate(value)
		if err != nil {
			return nil, fmt.Errorf("error parsing coordinate: %v", err)
		}

		switch key {
		case "Button A":
			currentMachine.ButtonA = coord
		case "Button B":
			currentMachine.ButtonB = coord
		case "Prize":
			currentMachine.Prize = coord
		default:
			return nil, fmt.Errorf("unknown key: %s", key)
		}
	}

	if currentMachine != (ClawMachine{}) {
		machines = append(machines, currentMachine)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return machines, nil
}

func parseCoordinate(s string) (Coordinate, error) {
	parts := strings.Split(s, ", ")
	if len(parts) != 2 {
		return Coordinate{}, fmt.Errorf("invalid coordinate format: %s", s)
	}

	x, err := parseValue(parts[0])
	if err != nil {
		return Coordinate{}, err
	}

	y, err := parseValue(parts[1])
	if err != nil {
		return Coordinate{}, err
	}

	return Coordinate{X: x, Y: y}, nil
}

func parseValue(s string) (int, error) {
	s = strings.TrimPrefix(s, "X")
	s = strings.TrimPrefix(s, "Y")
	s = strings.TrimPrefix(s, "=")
	s = strings.TrimPrefix(s, "+")

	return strconv.Atoi(s)
}

func solveClawMachine(machine ClawMachine, prizeOffset int) int {
	// Set up the linear system
	a11 := float64(machine.ButtonA.X)
	a12 := float64(machine.ButtonB.X)
	b1 := float64(machine.Prize.X + prizeOffset)
	a21 := float64(machine.ButtonA.Y)
	a22 := float64(machine.ButtonB.Y)
	b2 := float64(machine.Prize.Y + prizeOffset)

	// Calculate determinants
	det := a11*a22 - a12*a21
	detX := b1*a22 - b2*a12
	detY := a11*b2 - b1*a21

	// Check if the system has a solution
	if det == 0 {
		return -1 // No solution
	}

	// Solve the system
	x := detX / det
	y := detY / det

	// Check if the solution is non-negative and integral
	if x < 0 || y < 0 || math.Floor(x) != x || math.Floor(y) != y {
		return -1 // No non-negative integral solution
	}

	// Calculate the total tokens needed
	tokens := int(3*x + y)

	return tokens
}

func sumFewestTokens(machines []ClawMachine, prizeOffset int) int {
	const costA = 3
	const costB = 1
	var total int
	for _, machine := range machines {
		tokens := solveClawMachine(machine, prizeOffset)
		if tokens >= 0 {
			total += tokens
		}

	}
	return total
}

func main() {
	machines, err := getInput("input.txt")
	if err != nil {
		panic(err)
	}

	totalPart1 := sumFewestTokens(machines, 0)
	totalPart2 := sumFewestTokens(machines, 10000000000000)
	fmt.Println("[PART 1] Total tokens needed: ", totalPart1)
	fmt.Println("[PART 2] Total tokens needed: ", totalPart2)
}
