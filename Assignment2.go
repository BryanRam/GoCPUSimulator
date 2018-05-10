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
	//Initialise tag with a value of 10. This will be incremented over the duration of the forever loop
	tag := 10
	for {
		//get instruction from channel
		instruction := <-fromGenerateToDispatcher
		instruction += tag //append unique tag to it

		//declare a select clause that will listen from each pipeline to see if they are ready to receive an instruction
		select {
		//Listen to the readyForNext channels. If they have values, then the corresponding pipeline is ready to receive an instruction
		case <-readyForNext[0]:
			toPipeline[0] <- instruction
			outputGeneratedInstructions = append(outputGeneratedInstructions, instruction%10) //extract instruction and append it to appropriate output
			outputAssignedGeneratedPipeline = append(outputAssignedGeneratedPipeline, 0)      //do the same for the pipeline used
		case <-readyForNext[1]:
			toPipeline[1] <- instruction
			outputGeneratedInstructions = append(outputGeneratedInstructions, instruction%10)
			outputAssignedGeneratedPipeline = append(outputAssignedGeneratedPipeline, 1)
		case <-readyForNext[2]:
			toPipeline[2] <- instruction
			outputGeneratedInstructions = append(outputGeneratedInstructions, instruction%10)
			outputAssignedGeneratedPipeline = append(outputAssignedGeneratedPipeline, 2)

		}
		if (tag / 10) < numberOfInstructions+1 {
			tag += 10 //increment the tag as long as there are still instructions left to generate
		}

		//Display all generated instructions once no more instructions are pending to be generated
		if (tag / 10) == numberOfInstructions+1 {
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
		readyForNext <- id //indicate that the pipeline is ready for another instruction

		instruction := <-toPipeline //take instruction from toPipeline channel
		opcode := instruction % 10  //extract instruction without tag to execute delay

		time.Sleep(time.Duration(opcode) * time.Second) //Delay for opcode seconds

		fromPipeline <- instruction //retire instruction by sending it down the fromPipeline channel

	}

}

//--------------------------------------------------------------------------------
// Retires instructions by writing them to the console
//--------------------------------------------------------------------------------
func retireInstruction(fromPipeline [numberOfPipelines]chan int, sortedPipeInstructions [4]chan int,
	outputRetiredInstructions []int, outputCompletedInstructions []int, outputAssignedCompletedPipeline []int) {

	for { // do forever

		for i := 0; i < 3; i++ { //go through the length of sortedPipeInstructions
			go sortInstructions(sortedPipeInstructions[i], sortedPipeInstructions[i+1]) //call sortInstructions routine
		}

		for {
			//Get retired instructions from each pipeline, and transfer them to the first member of sortedPipeInstructions
			select {
			case x := <-fromPipeline[0]:
				sortedPipeInstructions[0] <- x
				outputCompletedInstructions = append(outputCompletedInstructions, x%10)      //extract instruction and append it to appropriate output
				outputAssignedCompletedPipeline = append(outputAssignedCompletedPipeline, 0) //do the same for the pipeline used

			case y := <-fromPipeline[1]:
				sortedPipeInstructions[0] <- y
				outputCompletedInstructions = append(outputCompletedInstructions, y%10)
				outputAssignedCompletedPipeline = append(outputAssignedCompletedPipeline, 1)

			case z := <-fromPipeline[2]:
				sortedPipeInstructions[0] <- z
				outputCompletedInstructions = append(outputCompletedInstructions, z%10)
				outputAssignedCompletedPipeline = append(outputAssignedCompletedPipeline, 2)

			case retired := <-sortedPipeInstructions[3]: //extract the sorted instruction
				outputRetiredInstructions = append(outputRetiredInstructions, retired%10) //append the retired instruction to the array

			}

			/*Three sorted retired instructions will always be left in the other members of the sortedPipeInstructions array.
			  So display the instructions in order of completion, the pipelines used, as well as the other retired instructions,
			  in the order that they were originally generated
			*/
			if len(outputRetiredInstructions) == (numberOfInstructions - 3) {
				fmt.Printf("\nCompleted: %v\n", outputCompletedInstructions)
				fmt.Printf("Pipelines: %v\n", outputAssignedCompletedPipeline)
				fmt.Printf("\nRetired:   %v\n", outputRetiredInstructions)

			}
		}

	}
}

//-------------------------------------------------------------
// sortInstructions is a sorting algorithm that swaps the relies on the nature
// of channels to get its values
//-------------------------------------------------------------
func sortInstructions(incoming <-chan int, current chan<- int) {
	/*Both i and j get the value of the current instruction, however i's value stays
	  the same as it is declared outside the forever loop (that is, until certain conditions within the loop are met).
	*/
	i := <-incoming

	for {
		j := <-incoming

		//if j is less than i
		if (j) < (i) {
			current <- j //change the value in the incoming channel with j

		} else {
			current <- i //otherwise change the value in the incoming channel with i
			i = j        //then give i j's value

		}

	}
	/*
	  Due to the nature of this loop, the sorted value would always end up in the incoming
	  channel. Since this go routine is called within a for loop on the sortedPipeInstructions array,
	  this means that the last member of sortedPipeInstructions gets the correctly sorted value upon termination
	  of the for loop.
	*/
}

//Takes input from stdin
func readInput() {
	var button string
	for {
		//read keyboard input
		fmt.Scan(&button)

		//exit the program if q or Q is entered
		//If run through VS Code, this form of keyboard input
		//will only be recognised if running an .exe of the program through the terminal
		if button == "Q" || button == "q" {
			fmt.Printf("Program terminated.\n")
			os.Exit(3)
		}
	}
}

/*
Constant Declarations:
Define the number of Pipelines and Instructions needed for the program
*/
const numberOfPipelines = 3
const numberOfInstructions = 15

//////////////////////////////////////////////////////////////////////////////////
//  Main program, create required channels, then start goroutines in parallel.
//////////////////////////////////////////////////////////////////////////////////

func main() {
	rand.Seed(time.Now().Unix()) // Seed the random number generator

	// Set up required channels
	fromGenerateToDispatcher := make(chan int) // channel for flow of generated opcodes

	var toPipeline [numberOfPipelines]chan int //channel array for sending a generated opcode to each pipeline
	for i := range toPipeline {
		toPipeline[i] = make(chan int) //create int channels for each member of the channel array
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

	/*
	  The following "output" variables are int arrays
	  used to display instructions to the screen at the various states
	  they take during the operation (Generated, Completed, Retired)
	  as well as the pipelines used
	*/
	outputGeneratedInstructions := make([]int, 0) //int array to hold generat
	outputAssignedGeneratedPipeline := make([]int, 0)

	outputCompletedInstructions := make([]int, 0)
	outputAssignedCompletedPipeline := make([]int, 0)

	outputRetiredInstructions := make([]int, 0)

	// Now start the goroutines in parallel.
	fmt.Printf("Start Go routines...\n")
	fmt.Printf("Enter 'q' or 'Q' to quit.\n")

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
