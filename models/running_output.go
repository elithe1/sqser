package models

type RunningOutput struct {
	Output
	Config *OutputConfig
}

type OutputConfig struct {
}

func NewRunningOutput(output Output, config *OutputConfig) *RunningOutput {
	return &RunningOutput{
		Output: output,
		Config: config,
	}
}
