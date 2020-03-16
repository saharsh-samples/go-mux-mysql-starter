package utils

import (
	"testing"
	"time"

	"github.com/saharsh-samples/go-mux-sql-starter/test"
)

func TestWaitTill_success_case(t *testing.T) {

	// Arrange
	till := time.Now().Add(500 * time.Millisecond)
	condition := func() bool { return time.Now().After(till) }

	// Act and assert
	test.AssertTrue("expected no timeout", nil == WaitTill(condition, 10), t)

}

func TestWaitTill_error_case(t *testing.T) {

	// Arrange
	till := time.Now().Add(24 * time.Hour)
	condition := func() bool { return time.Now().After(till) }

	// Act
	timeout := WaitTill(condition, 1)

	// Assert
	test.AssertFalse("expected timeout", nil == timeout, t)
	test.AssertEquals("", ErrTimeout, timeout.Error(), t)

}
