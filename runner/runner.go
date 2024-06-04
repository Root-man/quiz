package runner

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"

	"github.com/root-man/quiz/config"
	"github.com/root-man/quiz/problem"
)

type QuizRunner struct {
	problems []problem.Problem
	config   *config.QuizConfig
}

func New(config *config.QuizConfig) QuizRunner {
	file, err := os.Open(config.GetProblemsFile())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	problems := []problem.Problem{}

	for {
		data, err := reader.Read()
		if err == io.EOF {
			break
		}

		if len(data) != 2 {
			fmt.Printf("Got invalid problem %s, skipping\n", data)
		}

		question := data[0]
		answer := data[1]

		problems = append(problems, problem.New(question, answer))
	}

	if config.GetShuffle() {
		fmt.Println("Shuffling the list of questions...")
		rand.Shuffle(len(problems), func(i, j int) { problems[i], problems[j] = problems[j], problems[i] })
	}

	return QuizRunner{
		problems: problems,
		config:   config,
	}
}

func (r QuizRunner) Run() {
	userResults := make(chan bool)
	var correctAnswers, wrongAnswers int

	stdinReader := bufio.NewReader(os.Stdin)
	fmt.Printf("Press ENTER to start. You have %v to complete the quiz\n", r.config.GetTimer())
	stdinReader.ReadLine()
	context, cancel := context.WithTimeout(context.Background(), r.config.GetTimer())
	defer cancel()

problemLoop:
	for _, problem := range r.problems {
		go func() {
			answer, _ := problem.AskAndWaitForAnswer()
			userResults <- answer
		}()

		select {
		case result := <-userResults:
			if result {
				correctAnswers++
			} else {
				wrongAnswers++
			}
		case <-context.Done():
			fmt.Println("Quiz timed out")
			break problemLoop
		}
	}
	printQuizResults(correctAnswers, wrongAnswers, len(r.problems))
}

func printQuizResults(correct, incorrect, total int) {
	fmt.Printf("Your results...\nTotal questions: %d\n%d correct answers\n%d incorrect answers\n", total, correct, incorrect)
}
