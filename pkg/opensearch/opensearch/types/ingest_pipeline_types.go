package types

type IngestPipeline struct {
	Description string      `json:"description,omitempty"`
	Processors  []Processor `json:"processors,omitempty"`
}

type Processor struct {
	MontyLoggingProcessor *MontyProcessorConfig `json:"monty-logging-processor,omitempty"`
	MontyPreProcessor     *MontyProcessorConfig `json:"montypre,omitempty"`
}

type MontyProcessorConfig struct {
}
