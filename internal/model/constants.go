package model

const (
	// MaxACC and MinACC define the allowed value range for the accumulator.
	MaxACC                = 999
	MinACC                = -999
	StreamTypesNumber     = 2  // Defines how many stream types exist (e.g., INPUT, OUTPUT).
	IOPositionsNumber     = 4  // Defines how many input/output positions exist.
	MaxStreamValuesLength = 30 // Defines the maximum number of values a stream can hold.
	NodesNumber           = 12 // Defines the total number of nodes in the puzzle grid.
	NodeTypesNumber       = 2  // Defines how many node types exist (e.g., COMPUTE, DAMAGED)
)
