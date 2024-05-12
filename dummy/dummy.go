package dummy

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

type QuizOpts struct {
	problemsFile string
	timer        time.Duration
	shuffle      bool
}

func NewOpts(mods ...Option) *QuizOpts {
	opts := QuizOpts{}

	for _, mod := range mods {
		mod(&opts)
	}

	return &opts
}

type Option func(o *QuizOpts)

func WithProblemsFile(p string) Option {
	return func(o *QuizOpts) {
		o.problemsFile = p
	}
}

func WithTimer(t time.Duration) Option {
	return func(o *QuizOpts) {
		o.timer = t
	}
}

func WithShuffle(s bool) Option {
	return func(o *QuizOpts) {
		o.shuffle = s
	}
}

type Problem struct {
	question string
	answer   string
}

func RunQuiz(cfg *QuizOpts) {
	file, err := os.Open(cfg.problemsFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	problems := []Problem{}
	var correctAnswers, wrongAnswers int

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

		problems = append(problems, Problem{question, answer})
	}

	if cfg.shuffle {
		fmt.Println("doing shuffle...")
		rand.Shuffle(len(problems), func(i, j int) { problems[i], problems[j] = problems[j], problems[i] })
	}

	userResults := make(chan bool)

	stdinReader := bufio.NewReader(os.Stdin)
	fmt.Printf("Press ENTER to start. You have %v to complete the quiz\n", cfg.timer)
	stdinReader.ReadLine()
	context, cancel := context.WithTimeout(context.Background(), cfg.timer)
	defer cancel()

problemLoop:
	for _, v := range problems {
		go runProblem(v.question, v.answer, userResults)

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
	printQuizResults(correctAnswers, wrongAnswers, len(problems))
}

func printQuizResults(correct, incorrect, total int) {
	fmt.Printf("Your results...\nTotal questions: %d\n%d correct answers\n%d incorrect answers\n", total, correct, incorrect)
}

func runProblem(problem, solution string, resultsChan chan bool) {
	stdinReader := bufio.NewReader(os.Stdin)
	fmt.Printf("Question:\n%s\nYour answer?\n", problem)

	line, _, err := stdinReader.ReadLine()
	if err != nil {
		resultsChan <- false
	}

	answer := string(line)
	resultsChan <- answer == solution
}
