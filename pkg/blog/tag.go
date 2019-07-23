package blog

// Tag is a label attached to the post for the purpose of identification
type Tag struct {
	// Identifier of the tag
	ID string
	
	// Name of the tag
	Name string
	
	// Valid URL string composes with name and ID
	Slug string
}
