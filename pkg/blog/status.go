package blog

// Status for indicating the access level of the post
type Status string

func (s Status) String() string {
	return string(s)
}

// IsPublished returns "true" if status is Published
func (s Status) IsPublished() bool {
	return s == StatusPublished
}

// IsDraft returns "true" if status is Draft
func (s Status) IsDraft() bool {
	return s == StatusDraft
}

// StatusPublished indicates that post is public, accessible to everyone
const StatusPublished Status = "PUBLISHED"

// StatusDraft indicates that the post is private, can only be accessed by the author
const StatusDraft Status = "DRAFT"
