package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

var (
	wg        sync.WaitGroup
	score     int
	lineCount int
	answer    string
)

func createReader(fileName string) (*csv.Reader, error) {
	contents, err := os.ReadFile(fileName) // For read access.
	if err != nil {
		return nil, err
	}

	return csv.NewReader(bytes.NewReader(contents)), nil
}
func quiz(r *csv.Reader, wg *sync.WaitGroup, score, lineCount *int) {
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		*lineCount++

		fmt.Printf("%s\n", record[0])

		// Taking input from user
		fmt.Print("Answer: ")
		fmt.Scanln(&answer)
		if answer == record[1] {
			*score++
		}
	}

	wg.Done()
}

func timer(d *int, score, lineCount *int) {
	time.Sleep(time.Duration(*d) * time.Second)
	fmt.Printf("\nYour %d second timer has passed, and your score so far was %d correct out of %d\n", *d, *score, *lineCount)
	os.Exit(0)
}

func main() {
	timerPtr := flag.Int("timer", 30, "The time allowed to answer all questions.")

	flag.Parse()

	fileName := "problems.csv"
	r, err := createReader(fileName)
	if err != nil {
		log.Panic("Unable to read CSV file.")
	}

	fmt.Println("Go!")
	wg.Add(1)
	go quiz(r, &wg, &score, &lineCount)
	go timer(timerPtr, &score, &lineCount)
	wg.Wait()
	fmt.Printf("Your score was %d out of %d\n", score, lineCount)
}
