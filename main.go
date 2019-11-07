package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question, answer'")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	check(err, fmt.Sprintf("failed to open the CSV file: %s\n", *csvFilename))

	// once file is opened, want to read it, then parse it
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	check(err, fmt.Sprintf("Failed to parse provided csv file"))

	// lines is a 2d array => [[5+5 10] [1+1 2]]
	problems := parseLines(lines)

	numCorrect := 0

	for i, prob := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, prob.q)
		var answer string
		fmt.Scanf("%s\n", &answer) // create pointer to answer
		if answer == prob.a {
			numCorrect++
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
