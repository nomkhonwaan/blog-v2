package log_test

import (
	. "github.com/nomkhonwaan/myblog/pkg/log"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDefaultTimer_Now(t *testing.T) {
	// Given
	now := time.Now()
	timer := DefaultTimer(now)

	// When
	result := timer.Now()

	// Then
	assert.EqualValues(t, now.Format("2006-01-02 15:04:05"), result.Format("2006-01-02 15:04:05"))
}
