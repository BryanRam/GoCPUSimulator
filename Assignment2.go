//	Assignment 2.go
// 	7873 MAPS
// 	Doesn't do much!
// 	Bradley Ramsay
//	26004360

package main

// Imported packages

import (
	"fmt"       // for console I/O
	"math/rand" // for randomly creating opcodes
	"os"
	"time" // for the random number generator and 'executing' opcodes
)

//////////////////////////////////////////////////////////////////////////////////
// Function definitions
//////////////////////////////////////////////////////////////////////////////////

//----------------------------------------------------------------------------------
// Randomly generate an instruction 'opcode' between 1 and 5 and send to the retire function
//----------------------------------------------------------------------------------

func generateInstructions(instruction chan<- int) {

	for i := 0; i < 15; i++ { // do a limited number

		opcode := (rand.Intn(5) + 1) // Randomly generate a new opcode (between 1 and 5)
		//opcode := i // testing that all 15 instructions are generated

		fmt.Printf("Instruction: %d\n", opcode) // Report this to console display

		instruction <- opcode // Send the instruction for retirement
	}
}

func pipeline(id int, toPipeline chan int, fromPipeline chan int, readyForNext chan<- int) {
	//Delay execution for the duration of id (using the pipeline's index, but must be replaced with opcode)
	/* time.Sleep(time.Duration(id) * time.Second)
	fmt.Printf("Duration for: %d\n", id) */
	//instruction := <-toPipeline
	for {
		//fmt.Println("Ready for next instruction")
		readyForNext <- id
		instruction := <-toPipeline
		tag := instruction / 10
		opcode := instruction % 10
		fmt.Printf("Tag: %d\n", tag)
		//fmt.Printf("CURRENTLY IN PIPE: %d\n", instruction)
		time.Sleep(time.Duration(opcode) * time.Second)
		//fmt.Printf("Duration for: %d\t %d\n", id, instruction)
		fromPipeline <- opcode
		//default:

	}

}

func dispatcher(fromGenerateToDispatcher <-chan int, toPipeline [numberOfPipelines]chan int, readyForNext [numberOfPipelines]chan int) {
	i := 10
	for {
		instruction := <-fromGenerateToDispatcher
		instruction += i
		/* r1 := <-readyForNext[0]
		fmt.Println("Down pipe 0:: %d", r1)
		toPipeline[0] <- instruction

		r2 := <-readyForNext[1]
		fmt.Println("Down pipe 1:: %d", r2)
		toPipeline[1] <- instruction

		r3 := <-readyForNext[2]
		fmt.Println("Down pipe 2:: %d", r3)
		toPipeline[2] <- instruction */

		select {
		case <-readyForNext[0]:
			toPipeline[0] <- instruction
		case <-readyForNext[1]:
			toPipeline[1] <- instruction
		case <-readyForNext[2]:
			toPipeline[2] <- instruction
			//default:
			//fmt.Println("No pipe chosen")
			//time.Sleep(time.Duration(3) * time.Second)
		}
		i += 10
		//fmt.Println("Something")
	}

}

//--------------------------------------------------------------------------------
// Retires instructions by writing them to the console
//--------------------------------------------------------------------------------
func retireInstruction(fromPipeline [numberOfPipelines]chan int) {

	for { // do forever
		opcode1 := <-fromPipeline[0]

		fmt.Printf("Retired: %d\n", opcode1) // Report to console

		opcode2 := <-fromPipeline[1]
		fmt.Printf("Retired 2: %d\n", opcode2) // Report to console

		opcode3 := <-fromPipeline[2]
		fmt.Printf("Retired 3: %d\n", opcode3) // Report to console
		// Receive an instruction from the generator
		// opcode := <-fromPipeline

		// fmt.Printf("Retired: %d \n", opcode) // Report to console

	}
}

func sortInstructions(current chan int, incoming chan int, opcode int) {
	i := <-current
	j := <-incoming
	if j > i {
		incoming <- i
		current <- j
	}
}

//Takes input from stdin
func readInput() {
	var button string
	for {
		//read keyboard input
		fmt.Scan(&button)

		if button == "Q" || button == "q" {
			os.Exit(3)
		}
	}
}

//////////////////////////////////////////////////////////////////////////////////
//  Main program, create required channels, then start goroutines in parallel.
//////////////////////////////////////////////////////////////////////////////////
const numberOfPipelines = 3

func main() {
	rand.Seed(time.Now().Unix()) // Seed the random number generator

	// Set up required channel
	fromGenerateToDispatcher := make(chan int) // channel for flow of generated opcodes

	// current := make(chan int)
	// incoming := make(chan int)

	//toPipeline := make([]chan int, numberOfPipelines)
	var toPipeline [numberOfPipelines]chan int
	for i := range toPipeline {
		toPipeline[i] = make(chan int)
	}

	//readyForNext := make([]chan int, numberOfPipelines)
	var readyForNext [numberOfPipelines]chan int
	for i := range readyForNext {
		readyForNext[i] = make(chan int)

	}

	var fromPipeline [numberOfPipelines]chan int
	for i := range fromPipeline {
		fromPipeline[i] = make(chan int)
	}

	// Now start the goroutines in parallel.
	fmt.Printf("Start Go routines...\n")

	//create 3 pipelines
	for i := 0; i < numberOfPipelines; i++ {
		go pipeline(i, toPipeline[i], fromPipeline[i], readyForNext[i])
	}

	go generateInstructions(fromGenerateToDispatcher)

	go dispatcher(fromGenerateToDispatcher, toPipeline, readyForNext)
	go readInput()
	go retireInstruction(fromPipeline)

	for { // Needed to keep the 'main' process alive!
	}

} // end of main /////////////////////////////////////////////////////////////////
