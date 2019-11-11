package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question, answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	check(err, fmt.Sprintf("failed to open the CSV file: %s\n", *csvFilename))

	// once file is opened, want to read it, then parse it
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	check(err, fmt.Sprintf("Failed to parse provided csv file"))

	// lines is a 2d array => [[5+5 10] [1+1 2]]
	// want to turn them into array of problem types => [{q:5+5, a:10} {q:1+1, a:2]]
	problems := parseLines(lines)

	numCorrect := 0

	// start timer
	// NewTimer creates a new Timer that will send a message over its channel once time has expired
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	for i, prob := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, prob.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer) // create pointer to answer
			answerCh <- answer
		}()

		select {
		// if we get a message from the timer channel, stop the program
		case <-timer.C:
			fmt.Printf("\nYou scored %d out of %d.\n", numCorrect, len(problems))
			return
		// if we get a message from the answer channel, check it
		case answer := <-answerCh:
			if answer == prob.a {
				numCorrect++
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", numCorrect, len(problems))
}

// takes in the csv lines (2d array) and returns a slice of problems
func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))

	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			// trim space in case csv file isn't formatted correctly
			a: strings.TrimSpace(line[1]),
		}
	}

	return ret
}

// creating a type makes it easier to adjust in the future
// eg if wanted to use json instead of csv
type problem struct {
	q string
	a string
}

func check(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		os.Exit(1)
	}
}

// NOTES

// Package flag implements command-line flag parsing
// flag.String(name string, value string, usage string) *string
// After all flags are defined, call flag.Parse() to parse the
// command line into the defined flags.
// Flags may then be used directly. If you're using the flags themselves,
// they are all pointers; if you bind to variables, they're values.

// os.Exit(code int)
// Exit causes the current program to exit with the given status code.

// Scanf scans text read from standard input, storing successive
// space-separated values into successive arguments as determined by the format.
// gets rid of all leading and trailing spaces, so good for this program
// but not for ones that would take in a sentence
// https://ukiahsmith.com/blog/fmt-scanf-introduction/

// https://www.golang-book.com/books/intro/10
// A goroutine is a function that is capable of running concurrently with other functions.
// To create a goroutine we use the keyword go followed by a function invocation.
// Channels provide a way for two goroutines to communicate with one another and synchronize their execution.
// The <- (left arrow) operator is used to send and receive messages on the channel.
// c <- "ping" means send "ping"
// msg := <- c means receive a message and store it in msg.
// Go has a special statement called select which works like a switch but for channels.
