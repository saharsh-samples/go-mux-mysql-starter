package utils

import (
	"errors"
	"testing"

	"github.com/saharsh-samples/go-mux-sql-starter/test"
)

func TestJoinErrors(t *testing.T) {
	test.AssertEquals(
		"",
		"First Error; Second Error; Third Error",
		JoinErrors(
			errors.New("First Error"),
			errors.New("Second Error"),
			nil,
			errors.New("Third Error"),
		),
		t)
}
