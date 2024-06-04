package config

import "time"

type QuizConfig struct {
	problemsFile string
	timer        time.Duration
	shuffle      bool
}

func (qc *QuizConfig) GetProblemsFile() string {
	return qc.problemsFile
}

func (qc *QuizConfig) GetTimer() time.Duration {
	return qc.timer
}

func (qc *QuizConfig) GetShuffle() bool {
	return qc.shuffle
}

func NewConfig(mods ...Option) *QuizConfig {
	config := QuizConfig{}

	for _, mod := range mods {
		mod(&config)
	}

	return &config
}

type Option func(o *QuizConfig)

func WithProblemsFile(p string) Option {
	return func(o *QuizConfig) {
		o.problemsFile = p
	}
}

func WithTimer(t time.Duration) Option {
	return func(o *QuizConfig) {
		o.timer = t
	}
}

func WithShuffle(s bool) Option {
	return func(o *QuizConfig) {
		o.shuffle = s
	}
}
