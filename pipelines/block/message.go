package block

// Messager is any block that can process messages
type Messager interface {
	// Put message to pipeline block
	Put(msg *Message) error
	// Get messager's id
	Id() int
	// Get id for next messager or "" for none
	Next() int

	// Set id
	SetId(int)
}

// Message is the data chunk that each block gets as input
// Value is input value
// Returnvalue is expected to have outputted value
type Message struct {
	// Device id
	Device *string
	// Measurement name
	Measurement *string
	// Input value
	Value float64
	// Output value
	ReturnValue float64
	// Next block id
	NextBlock int
}

func (m *Message) SetNext(n int) {
	m.NextBlock = n
}
