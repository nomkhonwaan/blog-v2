package blog

// Engagement represents social network engagement of the object
type Engagement struct {
	// Total object shared counter
	ShareCount int `bson:"-" json:"shareCount" graphql:"shareCount"`
}
