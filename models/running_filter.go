package models

type RunningFilter struct {
	Filter
	Config *FilterConfig
}

type FilterConfig struct {
	Name   string     `yaml:"name"`
	Values RawMessage `yaml:"values"`
}

func NewRunningFilter(filter Filter, config *FilterConfig) *RunningFilter {
	return &RunningFilter{
		Filter: filter,
		Config: config,
	}
}
