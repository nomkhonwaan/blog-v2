package facebook

// URL represents a URL shared on a timeline or in a comment
type URL struct {
	Engagement Engagement `json:"engagement" graphql:"-"`
}

// Engagement is struct which contains number of counts of different ways people interacted with the URL
type Engagement struct {
	// Number of comments on the URL
	CommentCount int `json:"commentCount" graphql:"commentCount"`

	// Number of comments on the plugin gathered using the Comments Plugin on your site
	CommentPluginCount int `json:"commentPluginCount" graphql:"commentPluginCount"`

	// Number of reactions to the URL
	ReactionCount int `json:"reactionCount" graphql:"reactionCount"`

	// Number of times the URL was shared
	ShareCount int `json:"shareCount" graphql:"shareCount"`
}
