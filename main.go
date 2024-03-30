package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	score     int
	lineCount int
)

type problem struct {
	q string
	a string
}

func createReader(fileName string) (*csv.Reader, error) {
	contents, err := os.ReadFile(fileName) // For read access.
	if err != nil {
		return nil, err
	}

	return csv.NewReader(bytes.NewReader(contents)), nil
}

func parseCsv(r *csv.Reader) ([]problem, error) {
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	ret := make([]problem, len(records))

	for i, record := range records {
		ret[i] = problem{
			q: record[0],
			a: record[1],
		}
	}

	return ret, nil

}
func quiz(problems []problem, timer *time.Timer, score, lineCount *int) {
	for _, p := range problems {
		*lineCount++
		fmt.Printf("Problem: %s\n", p.q)

		// Taking input from user
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			return
		case answer := <-answerCh:
			if answer == p.a {
				*score++
			}
		}
	}
}

func main() {
	csvPtr := flag.String("csv", "problems.csv", "Path to the CSV file containing the problems in a 'question,answer' format.")
	timerPtr := flag.Int("timer", 30, "The time allowed to answer all questions.")

	flag.Parse()

	r, err := createReader(*csvPtr)
	if err != nil {
		log.Fatalf("Unable to find CSV file.")
	}

	problems, err := parseCsv(r)
	if err != nil {
		log.Fatalf("Unable to read CSV file.")
	}

	// fmt.Println("Go!")
	timer := time.NewTimer(time.Duration(*timerPtr) * time.Second)

	quiz(problems, timer, &score, &lineCount)

	fmt.Printf("Your score was %d out of %d\n", score, lineCount)
}
