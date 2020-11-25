package dbaccess

// DBAccesser is the Data Access Object interface
type DBAccesser interface {
	Connect() error
	Disconnect() error
	InsertUser(UserConfig) (InsertID, error)
	InsertFollower(Follower) (InsertID, error)
	InsertTweet(TweetConfig) (InsertID, error)
	GetUser(UserID string) (User, error)
	GetFollowers(UserID string) ([]Follower, error)
	GetTweets(UserID string) ([]Tweet, error)
	GetAllUsers() ([]User, error)
	GetAllFollowers() ([]Follower, error)
	GetAllTweets() ([]Tweet, error)
}

// InsertID is ID of a newly created item
type InsertID string

// User is a struct type containing all fields of a User record
type User struct {
	ID       string
	Username string
	Password string
}

// UserConfig is a struct type containing all fields of User excluding the ID
type UserConfig struct {
	Username string
	Password string
}

// Tweet is a struct type containing all fields of a Tweet record
type Tweet struct {
	ID       string
	UserID   string
	Username string
	Text     string
}

// TweetConfig is a struct type containing all fields of a Tweet excluding the ID
type TweetConfig struct {
	UserID   string
	Username string
	Text     string
}

// Follower is a unique pair of a follower's and followee's UserIDs
type Follower struct {
	FollowerUserID string
	FolloweeUserID string
}
