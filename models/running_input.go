package models

type RunningInput struct {
	Input  Input
	Config *InputConfig
}

type InputConfig struct {
	Name  string
	Async bool
}

func NewRunningInput(input Input, config *InputConfig) *RunningInput {
	return &RunningInput{
		Input:  input,
		Config: config,
	}
}
