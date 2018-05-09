//	Assignment 2.go
// 	7873 MAPS
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
// Randomly generate an instruction 'opcode' between 1 and 5 and send to the dispatcher function
//----------------------------------------------------------------------------------

func generateInstructions(instruction chan<- int) {

	for i := 0; i < numberOfInstructions; i++ { // do a limited number

		opcode := (rand.Intn(5) + 1) // Randomly generate a new opcode (between 1 and 5)

		instruction <- opcode // Send the instruction for retirement

	}

}

//------------------------------------------------------------------------------------
// Gets a generated instruction, then checks each pipeline to see if they are ready to
// receive an instruction. If yes, then that instruction is sent to the pipeline
// -----------------------------------------------------------------------------------
func dispatcher(fromGenerateToDispatcher <-chan int, toPipeline [numberOfPipelines]chan int, readyForNext [numberOfPipelines]chan int,
	outputGeneratedInstructions []int, outputAssignedGeneratedPipeline []int) {
	i := 10
	for {
		instruction := <-fromGenerateToDispatcher
		instruction += i

		select {
		case <-readyForNext[0]:
			toPipeline[0] <- instruction
			outputGeneratedInstructions = append(outputGeneratedInstructions, instruction%10)
			outputAssignedGeneratedPipeline = append(outputAssignedGeneratedPipeline, 0)
		case <-readyForNext[1]:
			toPipeline[1] <- instruction
			outputGeneratedInstructions = append(outputGeneratedInstructions, instruction%10)
			outputAssignedGeneratedPipeline = append(outputAssignedGeneratedPipeline, 1)
		case <-readyForNext[2]:
			toPipeline[2] <- instruction
			outputGeneratedInstructions = append(outputGeneratedInstructions, instruction%10)
			outputAssignedGeneratedPipeline = append(outputAssignedGeneratedPipeline, 2)

		}
		if (i / 10) < numberOfInstructions+1 {
			i += 10
		}

		//Display all generated instructions once that point has been reached
		if (i / 10) == numberOfInstructions+1 {
			fmt.Printf("Opcodes:   %v\n", outputGeneratedInstructions)
			fmt.Printf("Pipelines: %v\n", outputAssignedGeneratedPipeline)
		}

	}

}

//----------------------------------------------------------------------------------------------
// Pipeline function tells dispatcher when it is ready to receive an instruction, delays operation for
// as long as the instruction dictates, then sends that instruction to be retired
//----------------------------------------------------------------------------------------------
func pipeline(id int, toPipeline <-chan int, fromPipeline chan<- int, readyForNext chan<- int) {

	for {

		readyForNext <- id

		instruction := <-toPipeline
		opcode := instruction % 10

		//Delay for opcode seconds
		time.Sleep(time.Duration(opcode) * time.Second)

		fromPipeline <- instruction

	}

}

//--------------------------------------------------------------------------------
// Retires instructions by writing them to the console
//--------------------------------------------------------------------------------
func retireInstruction(fromPipeline [numberOfPipelines]chan int, sortedPipeInstructions [4]chan int,
	outputRetiredInstructions []int, outputCompletedInstructions []int, outputAssignedCompletedPipeline []int) {

	for { // do forever

		for i := 0; i < 3; i++ {
			go sortInstructions(sortedPipeInstructions[i], sortedPipeInstructions[i+1])
		}

		for {
			select {
			case x := <-fromPipeline[0]:
				x2 := x
				sortedPipeInstructions[0] <- x
				outputCompletedInstructions = append(outputCompletedInstructions, x2%10)
				outputAssignedCompletedPipeline = append(outputAssignedCompletedPipeline, 0)

			case y := <-fromPipeline[1]:
				y2 := y
				sortedPipeInstructions[0] <- y
				outputCompletedInstructions = append(outputCompletedInstructions, y2%10)
				outputAssignedCompletedPipeline = append(outputAssignedCompletedPipeline, 1)

			case z := <-fromPipeline[2]:
				z2 := z
				sortedPipeInstructions[0] <- z
				outputCompletedInstructions = append(outputCompletedInstructions, z2%10)
				outputAssignedCompletedPipeline = append(outputAssignedCompletedPipeline, 2)

			case retired := <-sortedPipeInstructions[3]:
				outputRetiredInstructions = append(outputRetiredInstructions, retired%10)

			}

			if len(outputRetiredInstructions) == (numberOfInstructions - 3) {
				fmt.Printf("\nCompleted: %v\n", outputCompletedInstructions)
				fmt.Printf("Pipelines: %v\n", outputAssignedCompletedPipeline)
				fmt.Printf("\nRetired:   %v\n", outputRetiredInstructions)

			}
		}

	}
}

func sortInstructions(incoming <-chan int, current chan<- int) {

	i := <-incoming

	for {
		j := <-incoming

		//if j's tag is higher than i's tag, swap them around
		if (j) < (i) {
			current <- j

		} else {
			current <- i
			i = j

		}

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
const numberOfInstructions = 15

func main() {
	rand.Seed(time.Now().Unix()) // Seed the random number generator

	// Set up required channels
	fromGenerateToDispatcher := make(chan int) // channel for flow of generated opcodes

	var toPipeline [numberOfPipelines]chan int //channel array for sending a generated opcode to each pipeline
	for i := range toPipeline {
		toPipeline[i] = make(chan int)
	}

	var readyForNext [numberOfPipelines]chan int //channel array to indicate to the dispatcher that a pipeline is free to receive an opcode
	for i := range readyForNext {
		readyForNext[i] = make(chan int)

	}

	var fromPipeline [numberOfPipelines]chan int //channel array for sending retired opcodes to retireInstructions
	for i := range fromPipeline {
		fromPipeline[i] = make(chan int)
	}

	var sortedPipeInstructions [4]chan int //channel array for sorting retired opcodes before they are displayed
	for i := range sortedPipeInstructions {
		sortedPipeInstructions[i] = make(chan int)
	}

	outputGeneratedInstructions := make([]int, 0)
	outputAssignedGeneratedPipeline := make([]int, 0)

	outputCompletedInstructions := make([]int, 0)
	outputAssignedCompletedPipeline := make([]int, 0)

	outputRetiredInstructions := make([]int, 0)

	// Now start the goroutines in parallel.
	fmt.Printf("Start Go routines...\n")

	//create 3 pipelines
	for i := 0; i < numberOfPipelines; i++ {
		go pipeline(i, toPipeline[i], fromPipeline[i], readyForNext[i])
	}

	//generate 15 instructions
	go generateInstructions(fromGenerateToDispatcher)

	//get those instructions and send them to each pipeline
	go dispatcher(fromGenerateToDispatcher, toPipeline, readyForNext, outputGeneratedInstructions, outputAssignedGeneratedPipeline)

	//check for user input
	go readInput()

	//receive retired instructions and display them to the screen
	go retireInstruction(fromPipeline, sortedPipeInstructions, outputRetiredInstructions, outputCompletedInstructions, outputAssignedCompletedPipeline)

	for { // Needed to keep the 'main' process alive!
	}

} // end of main /////////////////////////////////////////////////////////////////
