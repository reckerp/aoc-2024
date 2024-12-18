package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Register struct {
	name string
	data int64
}

type Program struct {
	ptr int
	ops []int64
}

func main() {
	a, b, c, prg := getInput("input.txt")

	prog := Program{
		ptr: 0,
		ops: prg,
	}

	regs := []Register{
		{name: "A", data: a},
		{name: "B", data: b},
		{name: "C", data: c},
	}

	p1 := part1(prog, regs)
	p2 := part2(prog, regs)

	fmt.Println("[PART 1]: ", p1)
	fmt.Println("[PART 2]: ", p2)

}

func getInput(filename string) (int64, int64, int64, []int64) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var A, B, C int64
	var program []int64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Register A:") {
			A, _ = strconv.ParseInt(strings.TrimSpace(strings.Split(line, ":")[1]), 10, 64)
		} else if strings.Contains(line, "Register B:") {
			B, _ = strconv.ParseInt(strings.TrimSpace(strings.Split(line, ":")[1]), 10, 64)
		} else if strings.Contains(line, "Register C:") {
			C, _ = strconv.ParseInt(strings.TrimSpace(strings.Split(line, ":")[1]), 10, 64)
		} else if strings.Contains(line, "Program:") {
			programStr := strings.TrimSpace(strings.Split(line, ":")[1])
			programStrs := strings.Split(programStr, ",")
			program = make([]int64, len(programStrs))
			for i, s := range programStrs {
				intVar, _ := strconv.Atoi(s)
				program[i] = int64(intVar)
			}
		}
	}

	return A, B, C, program
}

func part1(prog Program, regs []Register) string {
	vals := runProgram(prog, regs)
	s := ""
	for i := 0; i < len(vals); i++ {
		s += fmt.Sprintf("%d", vals[i])
		if i != len(vals)-1 {
			s += ","
		}
	}
	return s
}

func part2(prog Program, regs []Register) int64 {
	return findQuine(prog, regs)
}

type State struct {
	segs []int64
}

func findQuine(prog Program, regs []Register) int64 {
	queue := []State{}
	for i := 0; i < 8; i++ {
		queue = append(queue, State{[]int64{int64(i)}})
	}
	var final int64
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		var x int64
		for i := len(cur.segs) - 1; i >= 0; i-- {
			s := cur.segs[i] << (3 * i)
			x = x | s
		}

		regs[0].data = x
		vals := runProgram(prog, regs)
		vp := 0
		matched := true
		for p := len(prog.ops) - len(vals); p < len(prog.ops); p++ {
			if vals[vp] != prog.ops[p] {
				matched = false
				break
			}
			vp++
		}

		done := matched && len(prog.ops) == len(vals)
		if done {
			final = x
			break
		}

		if matched {
			for i := 0; i < 8; i++ {
				nseg := make([]int64, len(cur.segs))
				copy(nseg, cur.segs)
				nseg = append([]int64{int64(i)}, nseg...)
				queue = append(queue, State{nseg})
			}
		}
	}
	return final
}

func runProgram(prog Program, regs []Register) []int64 {
	output := []int64{}
	rm := make(map[string]Register)
	for _, r := range regs {
		rm[r.name] = r
	}

	getVal := func(i int64) int64 {
		switch i {
		case 0, 1, 2, 3, 7:
			return i
		case 4:
			return rm["A"].data
		case 5:
			return rm["B"].data
		case 6:
			return rm["C"].data
		}
		return -1
	}

	setVal := func(r string, i int64) {
		reg := rm[r]
		reg.data = i
		rm[r] = reg
	}

	for {
		if prog.ptr >= len(prog.ops) {
			break
		}
		cur := prog.ops[prog.ptr]
		operand := prog.ops[prog.ptr+1]

		dojump := true
		jmp := 2
		switch cur {
		case 0:
			n := rm["A"].data
			d := int64(math.Pow(2, float64(getVal(operand))))
			if d != 0 {
				v := n / d
				setVal("A", v)
			}
		case 1:
			n := rm["B"].data
			v := n ^ operand
			setVal("B", v)
		case 2:
			n := getVal(operand)
			n = n % 8
			setVal("B", n)
		case 3:
			jnz := rm["A"].data
			if jnz != 0 {
				dojump = false
				prog.ptr = int(operand)
			}
		case 4:
			b := rm["B"].data
			c := rm["C"].data
			v := b ^ c
			setVal("B", v)
		case 5:
			v := getVal(operand) % 8
			output = append(output, v)
		case 6:
			n := rm["A"].data
			d := int64(math.Pow(2, float64(getVal(operand))))
			if d != 0 {
				v := n / d
				setVal("B", v)
			}
		case 7:
			n := rm["A"].data
			d := int64(math.Pow(2, float64(getVal(operand))))
			if d != 0 {
				v := n / d
				setVal("C", v)
			}
		}
		if dojump {
			prog.ptr += jmp
		}

	}
	return output
}
