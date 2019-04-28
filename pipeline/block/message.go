package block

// Messager is any block that can process messages
type Messager interface {
	// Put message to pipeline block
	Put(msg *Message) error
}

// Message is the data chunk that each block gets as input
// Value is input value
// Returnvalue is expected to have outputted value
type Message struct {
	Device      *string
	Measurement *string
	Value       float64
	ReturnValue float64
}
