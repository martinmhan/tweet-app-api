package event

// An Event contains the information passed to the message queue to publish an event
type Event struct {
	Type    Type
	Payload interface{}
}

// Type specifies the type of event to pass to the message queue
type Type int

const (
	// UserCreation is an event type that creates a User
	UserCreation Type = iota
	// TweetCreation is an event type that creates a Tweet
	TweetCreation
	// FollowCreation is an event type that creates a Follow
	FollowCreation
)

func (t Type) String() string {
	types := [...]string{
		"UserCreation",
		"TweetCreation",
		"FollowCreation",
	}

	return types[t]
}

// Producer is the event producer interface
type Producer interface {
	Produce(Event) error
}
