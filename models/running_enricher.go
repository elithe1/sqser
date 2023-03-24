package models

type RunningEnricher struct {
	Enricher
	Config *EnricherConfig
}

type EnricherConfig struct {
	Name   string     `yaml:"name"`
	Values RawMessage `yaml:"values"`
}

func NewRunningEnricher(enricher Enricher, config *EnricherConfig) *RunningEnricher {
	return &RunningEnricher{
		Enricher: enricher,
		Config:   config,
	}
}

type RawMessage struct {
	unmarshal func(interface{}) error
}

func (msg *RawMessage) UnmarshalYAML(unmarshal func(interface{}) error) error {
	msg.unmarshal = unmarshal
	return nil
}

// call this method later - when we know what concrete type to use
func (msg *RawMessage) Unmarshal(v interface{}) error {
	return msg.unmarshal(v)
}
