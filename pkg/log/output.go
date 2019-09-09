package log

import (
	"log"
	"os"
)

// Outputer is a compatible interface for logging with format
type Outputer interface {
	// Log with format to the output
	Printf(format string, args ...interface{})
}

// DefaultOutputer implements Outputer interface which embeds `log.Logger` inside
type DefaultOutputer struct {
	*log.Logger
}

func (outputer DefaultOutputer) Printf(format string, args ...interface{}) {
	outputer.Logger.Printf(format, args...)
}

// NewDefaultOutputer returns new `DefaultOutputer`
func NewDefaultOutputer() DefaultOutputer {
	return DefaultOutputer{log.New(os.Stdout, "", 0)}
}
