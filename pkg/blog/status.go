package blog

// Status for indicating the access level of the post
type Status string

// Published indicates that post is public, accessible to everyone
const Published Status = "PUBLISHED"

// Draft indicates that the post is private, can only be accessed by the author
const Draft Status = "DRAFT"

// PendingReview indicates that the post has been submitted to check and wait for approval from the system administrator
const PendingReview Status = "PENDING_REVIEW"
