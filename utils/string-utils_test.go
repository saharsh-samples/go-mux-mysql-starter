package utils

import (
	"fmt"
	"testing"

	"github.com/saharsh-samples/go-mux-sql-starter/test"
)

func TestIsNonEmptyString(t *testing.T) {
	test.AssertTrue("", IsNonEmptyString("blah"), t)
	test.AssertFalse("", IsNonEmptyString("  "), t)
	test.AssertFalse("", IsNonEmptyString(""), t)
	test.AssertFalse("", IsNonEmptyString("\t\n\n\t"), t)
}

func TestIsEmptyString(t *testing.T) {
	test.AssertFalse("", IsEmptyString("blah"), t)
	test.AssertTrue("", IsEmptyString("  "), t)
	test.AssertTrue("", IsEmptyString(""), t)
	test.AssertTrue("", IsEmptyString("\t\n\n\t"), t)
}

func TestIsStringMissingInSlice(t *testing.T) {
	test.AssertFalse("", IsStringMissingInSlice("boom", []string{"dhoom", "boom", "vroom"}), t)
	test.AssertTrue("", IsStringMissingInSlice("boom", []string{}), t)
	test.AssertTrue("", IsStringMissingInSlice("boom", nil), t)

	var doesNotContainBoom []string
	for i := 0; i < 10000; i++ {
		doesNotContainBoom = append(doesNotContainBoom, fmt.Sprintf("%d", i))
	}
	test.AssertTrue("", IsStringMissingInSlice("boom", doesNotContainBoom), t)
}

func TestIsStringValidEmail(t *testing.T) {
	test.AssertFalse("Expected string to fail email validation", IsStringValidEmail("i-am-not-an-email"), t)
	test.AssertTrue("Expected string to pass email validation", IsStringValidEmail("some-address@email.com"), t)
}
