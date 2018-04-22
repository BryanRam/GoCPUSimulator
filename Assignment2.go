// Example CPU execute pipeline - SimpleCPU.go
// 7873 MAPS
// Doesn't do much!
// A.Oram 2017

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

		//opcode := (rand.Intn(5) + 1) // Randomly generate a new opcode (between 1 and 5)
		opcode := i // Randomly generate a new opcode (between 1 and 5)

		fmt.Printf("Instruction: %d\n", opcode) // Report this to console display

		instruction <- opcode // Send the instruction for retirement
	}
}

func pipeline(id int, toPipeline <-chan int, fromPipeline <-chan int, readyForNext <-chan int) {
	time.Sleep(time.Duration(id*2) * time.Second)
	fmt.Printf("Duration for: %d\n", id)

}

func dispatcher(opcodes chan<- int, toPipeline <-chan int, readyForNext <-chan int) {

}

//--------------------------------------------------------------------------------
// Retires instructions by writing them to the console
//--------------------------------------------------------------------------------
func retireInstruction(retired <-chan int) {

	for { // do forever
		// Receive an instruction from the generator
		opcode := <-retired

		fmt.Printf("Retired: %d \n", opcode) // Report to console
	}
}

//Takes input from stdin
func readInput() {
	var button string
	for {
		//fmt.Printf("In readInput\n")
		fmt.Scan(&button)
		//reader := bufio.NewReader(os.Stdin)
		//input, _ := reader.ReadString('\n')
		//buf.Read([]byte(button))

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

	opcodes := make(chan int) // channel for flow of generated opcodes
	toPipeline := make([]chan int, numberOfPipelines)

	readyForNext := make([]chan int, numberOfPipelines)
	fromPipeline := make([]chan int, numberOfPipelines)

	// Now start the goroutines in parallel.
	fmt.Printf("Start Go routines...\n")

	for i := 0; i < numberOfPipelines; i++ {
		go pipeline(i, toPipeline[i], fromPipeline[i], readyForNext[i])
	}
	go generateInstructions(opcodes)
	go readInput()
	go retireInstruction(opcodes)

	for { // Needed to keep the 'main' process alive!
	}

} // end of main /////////////////////////////////////////////////////////////////
