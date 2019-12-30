package blog

// Status for indicating the access level of the post
type Status string

func (s Status) String() string {
	return string(s)
}

// IsPublished returns "true" if status is Published
func (s Status) IsPublished() bool {
	return s == Published
}

// IsDraft returns "true" if status is Draft
func (s Status) IsDraft() bool {
	return s == Draft
}

// Published indicates that post is public, accessible to everyone
const Published Status = "PUBLISHED"

// Draft indicates that the post is private, can only be accessed by the author
const Draft Status = "DRAFT"

// PendingReview indicates that the post has been submitted to check and wait for approval from the system administrator
const PendingReview Status = "PENDING_REVIEW"
