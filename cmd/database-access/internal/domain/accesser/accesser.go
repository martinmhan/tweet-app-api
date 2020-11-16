package accesser

// DBAccesser is the Data Access Object interface
type DBAccesser interface {
	Connect() error
	InsertUser(UserFields) User
	InsertTweet(TweetFields) Tweet
	GetUser(UserID string) User
	GetTweets(UserID string) []Tweet
}

// User is a struct type containing all fields of an Item record
type User struct {
	ID       string
	Username string
	Password string
}

// UserFields is a struct type containing all fields of Item excluding the ID
type UserFields struct {
	Username string
	Password string
}

// Tweet is a struct type containing all fields of a Tweet record
type Tweet struct {
	ID   string
	Text string
}

// TweetFields is a struct type containing all fields of a Tweet excluding the ID
type TweetFields struct {
	Text string
}
