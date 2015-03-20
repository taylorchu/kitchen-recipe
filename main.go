package main

import (
	"fmt"

	"github.com/taylorchu/blueprint/nl"
)

func main() {
	for i, s := range Instr1 {
		fmt.Printf("Step %d: %s\n", i, s.String())
	}
	fmt.Println("step count:", len(Instr1))
	fmt.Println("things:", Instr1.Thing())
	fmt.Println("---")

	for i, d := range Instr1.Ref() {
		fmt.Printf("Step %d directly depends on %v\n", i, d)
	}
	fmt.Println("---")

	for i, d := range Instr1.AllRef() {
		fmt.Printf("Step %d depends on %v\n", i, d)
	}
	fmt.Println("---")

	var steps []*ContinueStep
	for _, s := range Input1 {
		chunks, err := nl.Chunk(Expand(Grammar1, "SENT"), s)
		if err != nil {
			fmt.Println("!", err, s)
		} else {
			fmt.Println(len(chunks), chunks)
			steps = append(steps, ParseStep(chunks)...)
		}
	}
	instr := ParseInstruction(steps)
	fmt.Println("---")

	fmt.Println(instr)
	fmt.Println("---")

	for i, s := range instr {
		fmt.Printf("Step %d: %s\n", i, s.String())
	}
	fmt.Println("---")

	for i, d := range instr.Ref() {
		fmt.Printf("Step %d directly depends on %v\n", i, d)
	}
}
