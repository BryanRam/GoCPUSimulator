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
// Randomly generate an instruction 'opcode' between 1 and 5 and send to the dispatcher function
//----------------------------------------------------------------------------------

func generateInstructions(instruction chan<- int) {

	for i := 0; i < numberOfInstructions; i++ { // do a limited number

		opcode := (rand.Intn(5) + 1) // Randomly generate a new opcode (between 1 and 5)
		//opcode := i // testing that all 15 instructions are generated

		//fmt.Printf("Instruction: %d\n", opcode) // Report this to console display

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
		i += 10
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
func pipeline(id int, toPipeline <-chan int, fromPipeline chan<- int, readyForNext chan<- int, outputCompletedInstructions []int, outputAssignedCompletedPipeline []int) {

	for {

		readyForNext <- id
		instruction := <-toPipeline
		//tag := instruction / 10
		opcode := instruction % 10
		//fmt.Printf("Instruction: %d Tag: %d\n", opcode, tag)
		//Delay for opcode seconds
		time.Sleep(time.Duration(opcode) * time.Second)

		/*select {
		case id == 0:
			outputCompletedInstructions = append(outputCompletedInstructions, opcode)
			outputAssignedCompletedPipeline = append(outputAssignedCompletedPipeline, id)
		case id == 1:
			outputCompletedInstructions = append(outputCompletedInstructions, opcode)
			outputAssignedCompletedPipeline = append(outputAssignedCompletedPipeline, id)
		case id == 2:
			outputCompletedInstructions = append(outputCompletedInstructions, opcode)
			outputAssignedCompletedPipeline = append(outputAssignedCompletedPipeline, id)

		}*/
		//outputCompletedInstructions = append(outputCompletedInstructions, opcode)
		//outputAssignedCompletedPipeline = append(outputAssignedCompletedPipeline, id)

		fromPipeline <- instruction

		/*
			if (instruction / 10) == numberOfInstructions {
				fmt.Printf("\nCompleted: %v\n", outputCompletedInstructions)
				fmt.Printf("\nPipelines: %v\n", outputAssignedCompletedPipeline)
			}*/

		//go formatCompletedInstructionString(outputCompletedInstructions, outputAssignedCompletedPipeline, id, instruction)

		//fmt.Printf("pipeline in\n")

	}

}

func formatCompletedInstructionString(outputCompletedInstructions []int, outputAssignedCompletedPipeline []int, id int, instruction int) {
	for {
		outputCompletedInstructions = append(outputCompletedInstructions, instruction%10)
		outputAssignedCompletedPipeline = append(outputAssignedCompletedPipeline, id)

		//if (instruction / 10) == numberOfInstructions {
		fmt.Printf("\nCompleted: %v\n", outputCompletedInstructions)
		fmt.Printf("\nPipelines: %v\n", outputAssignedCompletedPipeline)
		//}
	}
}

//--------------------------------------------------------------------------------
// Retires instructions by writing them to the console
//--------------------------------------------------------------------------------
func retireInstruction(fromPipeline [numberOfPipelines]chan int, sortedPipeInstructions [4]chan int,
	outputRetiredInstructions []int, outputAssignedRetiredPipeline []int,
	outputCompletedInstructions []int, outputAssignedCompletedPipeline []int) {

	for { // do forever

		//sortInstructions(fromPipeline[0], fromPipeline[1])
		//sortInstructions(fromPipeline[0], fromPipeline[2])
		//sortInstructions(fromPipeline[1], fromPipeline[2])
		for i := 0; i < 3; i++ {
			go sortInstructions(sortedPipeInstructions[i], sortedPipeInstructions[i+1])
		}

		// tag1 := <-sortedPipeInstructions[0]
		// tag2 := <-sortedPipeInstructions[1]
		// tag3 := <-sortedPipeInstructions[2]

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
				//fmt.Printf("\nRetired Tag: %d\n\n", retired/10)
			}

			if len(outputRetiredInstructions) == (numberOfInstructions - 3) {
				fmt.Printf("\nCompleted: %v\n", outputCompletedInstructions)
				fmt.Printf("Pipelines: %v\n", outputAssignedCompletedPipeline)
				fmt.Printf("\nRetired:   %v\n", outputRetiredInstructions)
				//fmt.Printf("Pipelines: %v\n", outputAssignedRetiredPipeline)
			}
		}

		/* fmt.Printf("Retired Tag: %d\n", tag1) // Report to console
		fmt.Printf("Retired Tag: %d\n", tag2) // Report to console
		fmt.Printf("Retired Tag: %d\n", tag3) // Report to console */
	}
}

//---------------------------------------------------------------------------
//goRoutine is just a place to put all the sortInstructions function calls
//---------------------------------------------------------------------------
func goRoutine(fromPipeline [numberOfPipelines]chan int, sortedPipeInstructions [numberOfPipelines]chan int) {
	/*sortInstructions(fromPipeline[0], fromPipeline[1], sortedPipeInstructions[0])
	sortInstructions(fromPipeline[1], fromPipeline[2], sortedPipeInstructions[1])
	sortInstructions(fromPipeline[0], fromPipeline[2], sortedPipeInstructions[0])

		pipe1 := <-fromPipeline[0]
		pipe2 := <-fromPipeline[1]
		pipe3 := <-fromPipeline[2]
		opcode1 := pipe1 % 10
		tag1 := pipe1 / 10
		opcode2 := pipe2 % 10
		tag2 := pipe2 / 10
		opcode3 := pipe3 % 10
		tag3 := pipe3 / 10

		fmt.Printf("Retired: %d Tag: %d\n", opcode1, tag1)   // Report to console
		fmt.Printf("Retired 2: %d Tag: %d\n", opcode2, tag2) // Report to console
		fmt.Printf("Retired 3: %d Tag: %d\n", opcode3, tag3) // Report to console
	*/
	// Receive an instruction from the generator
	// opcode := <-fromPipeline

}

func sortInstructions(incoming <-chan int, current chan<- int) { //, sortedSecond chan<- int incoming <-chan int, /*, sortedFirst chan<- int*/

	//fmt.Printf("In sort\n")

	i := <-incoming

	for {
		j := <-incoming

		//if j's tag is higher than i's tag, swap them around
		if (j) < (i) {
			current <- j
			//sortedSecond <- j / 10

			//fmt.Printf("Sorted current: %d incoming: %d \n", j/10, i/10)
		} else {
			current <- i
			i = j

			//fmt.Printf("Sorted current: %d incoming: %d \n", i/10, j/10)
		}

	}
	// j := <-incoming

	// //if j's tag is higher than i's tag, swap them around
	// if (j / 10) < (i / 10) {
	// 	sortedFirst <- i / 10
	// 	sortedSecond <- j / 10

	// 	//fmt.Printf("Sorted current: %d incoming: %d \n", j/10, i/10)
	// } else {
	// 	sortedFirst <- j / 10
	// 	sortedSecond <- i / 10

	// 	//fmt.Printf("Sorted current: %d incoming: %d \n", i/10, j/10)
	// }

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

	// current := make(chan int)
	// incoming := make(chan int)

	//toPipeline := make([]chan int, numberOfPipelines)
	var toPipeline [numberOfPipelines]chan int //channel array for sending a generated opcode to each pipeline
	for i := range toPipeline {
		toPipeline[i] = make(chan int)
	}

	//readyForNext := make([]chan int, numberOfPipelines)
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
	outputAssignedRetiredPipeline := make([]int, 0)

	// Now start the goroutines in parallel.
	fmt.Printf("Start Go routines...\n")

	//create 3 pipelines
	for i := 0; i < numberOfPipelines; i++ {
		go pipeline(i, toPipeline[i], fromPipeline[i], readyForNext[i], outputCompletedInstructions, outputAssignedCompletedPipeline)
	}

	//generate 15 instructions
	go generateInstructions(fromGenerateToDispatcher)

	//get those instructions and send them to each pipeline
	go dispatcher(fromGenerateToDispatcher, toPipeline, readyForNext, outputGeneratedInstructions, outputAssignedGeneratedPipeline)

	//check for user input
	go readInput()

	//receive retired instructions and display them to the screen
	go retireInstruction(fromPipeline, sortedPipeInstructions, outputRetiredInstructions, outputAssignedRetiredPipeline, outputCompletedInstructions, outputAssignedCompletedPipeline)

	for { // Needed to keep the 'main' process alive!
	}

} // end of main /////////////////////////////////////////////////////////////////
