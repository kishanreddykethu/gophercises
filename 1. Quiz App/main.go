package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)


type problem struct {
	question string
	answer string
}


func main() {
	csvFileName := flag.String("csv", "problems.csv", "a csv file in the format of 'question, answer'")
	timeLimit := flag.Int("limit", 30, "The time limit for the quic in seconds")
	flag.Parse()
	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFileName))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the CSV file.")
	}
	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0
	for i, problem := range problems {
		fmt.Printf("Problem #%d: %s = \n", i + 1, problem.question)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
			case <-timer.C:
				fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
				return
			case answer := <-answerCh:
				if answer == problem.answer {
					correct++
				}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}


func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem {
			question: line[0],
			answer: line[1],
		}
	}
	return ret
}


func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}