package problem

import (
	"bufio"
	"fmt"
	"os"
)

type Problem struct {
	question string
	answer   string
}

func New(question, answer string) Problem {
	return Problem{question, answer}
}

func (p Problem) AskAndWaitForAnswer() (bool, error) {
	stdinReader := bufio.NewReader(os.Stdin)
	fmt.Printf("Question:\n%s\nYour answer?\n", p.question)

	line, _, err := stdinReader.ReadLine()
	if err != nil {
		return false, err
	}

	userAnswer := string(line)
	return userAnswer == p.answer, nil
}
