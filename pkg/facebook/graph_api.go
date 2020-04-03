package facebook

// URLNode is a shared URLNode on the timeline or in a comment
type URLNode struct {
	// The URLNode itself
	ID string `json:"id"`

	// Engagement is struct which contains number of counts of different ways people interacted with the URLNode
	Engagement struct {
		// Number of comments on the URLNode
		CommentCount int `json:"comment_count"`

		// Number of comments on the plugin gathered using the Comments Plugin on your site
		CommentPluginCount int `json:"comment_plugin_count"`

		// Number of reactions to the URLNode
		ReactionCount int `json:"reaction_count"`

		// Number of times the URLNode was shared
		ShareCount int `json:"share_count"`
	} `json:"engagement"`
}
