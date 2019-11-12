package blog

// Engagement represents social network engagement of the object
type Engagement struct {
	ShareCount int `bson:"-" json:"shareCount" graphql:"shareCount"`
}
