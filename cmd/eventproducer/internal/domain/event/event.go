package event

// An Event contains the information passed to the message queue to publish an event
type Event struct {
	Type    Type
	Payload interface{}
}

// Type specifies the type of event to pass to the message queue
type Type int

const (
	// UserCreation is an event type
	UserCreation Type = iota
	// TweetCreation is an event type
	TweetCreation
	// FollowerCreation is an event type
	FollowerCreation
)

func (t Type) String() string {
	types := [...]string{
		"UserCreation",
		"TweetCreation",
		"FollowerCreation",
	}

	return types[t]
}

// Producer is the event producer interface
type Producer interface {
	Produce(Event) error
}
