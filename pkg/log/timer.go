package log

import "time"

// Timer is a compatible interface for retrieving current system date-time
type Timer interface {
	// Return current system date-time
	Now() time.Time
}

// DefaultTimer implements Timer interface which returns current system date-time from `time.Now()` function
type DefaultTimer time.Time

func (timer *DefaultTimer) Now() time.Time {
	*timer = DefaultTimer(time.Now())
	return time.Time(*timer)
}

// NewDefaultTimer returns new `DefaultTimer` pointer
func NewDefaultTimer() *DefaultTimer {
	return new(DefaultTimer)
}
