package event

// Event is a generic struct that contains a Type and Payload for the consumer to read and execute with
type Event struct {
	Type    string // use event type enum
	Payload interface{}
}

// TO DO - create EventType enum type (e.g., "Create User", "Create Tweet")
