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
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	limit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	problems := readProblems(*csvFilename)
	correct := runQuiz(problems, *limit)

	fmt.Printf("\nYou scored %d out of %d\n", correct, len(problems))
}

type Problem struct {
	question string
	answer   string
}

func readProblems(fileName string) []Problem {
	file, err := os.Open(fileName)
	if err != nil {
		msg := fmt.Sprintf("Failed to open the CSV file: %s", fileName)
		exit(msg)
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provied CSV file")
	}

	return parseLines(lines)
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func parseLines(lines [][]string) []Problem {
	res := make([]Problem, len(lines))
	for i, line := range lines {
		res[i] = Problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return res
}

func runQuiz(problems []Problem, limit int) int {
	timer := time.NewTimer(time.Duration(limit) * time.Second)
	correct := 0
	for i, p := range problems {
		answerCh := make(chan bool)
		go func() {
			answerCh <- askUser(i+1, p)
		}()

		select {
		case <-timer.C:
			return correct
		case res := <-answerCh:
			if res {
				correct++
			}
		}
	}
	return correct
}

func askUser(id int, problem Problem) bool {
	fmt.Printf("Problem %d: %s = ", id, problem.question)
	var answer string
	fmt.Scanf("%s", &answer)
	return answer == problem.answer
}
