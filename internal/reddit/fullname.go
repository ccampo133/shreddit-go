package reddit

const (
	// Reddit "things" (e.g. comments, posts) have "fullnames", which are
	// unique identifiers constructed as a prefix followed by some opaque
	// ID string. See https://www.reddit.com/dev/api/#fullnames.
	commentPrefix = "t1_"
	postPrefix    = "t3_"
)

// commentFullName returns the fullname of a comment given its ID.
func commentFullName(id string) string {
	return commentPrefix + id
}

// postFullName returns the fullname of a post given its ID.
func postFullName(id string) string {
	return postPrefix + id
}
