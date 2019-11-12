package facebook

// URL represents a URL shared on a timeline or in a comment
type URL struct {
	// Engagement is struct which contains number of counts of different ways people interacted with the URL
	Engagement struct {
		// Number of comments on the URL
		CommentCount int `json:"comment_count"`

		// Number of comments on the plugin gathered using the Comments Plugin on your site
		CommentPluginCount int `json:"comment_plugin_count"`

		// Number of reactions to the URL
		ReactionCount int `json:"reaction_count"`

		// Number of times the URL was shared
		ShareCount int `json:"share_count"`
	} `json:"engagement"`
}
