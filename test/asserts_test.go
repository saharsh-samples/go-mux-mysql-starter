package test

import (
	"fmt"
	"testing"
)

func TestCompliance(t *testing.T) {

	AssertTrue("", true, t)

	AssertFalse("", false, t)

	AssertEquals("", 1, 1, t)

	defer AssertPanic("", t)
	panic("Intentionally triggered panic")

}

func TestAssertTrue(t *testing.T) {

	// success case
	d := &dummyT{}
	AssertTrue("Condition not true", true, d)

	if d.fail {
		t.Error("Expected dummyT.fail to be false")
	}
	if d.errorMsg != "" {
		t.Error("Expected dummyT.errorMsg to be blank")
	}

	// error case
	d = &dummyT{}
	AssertTrue("Condition not true", false, d)

	if !d.fail {
		t.Error("Expected dummyT.fail to be true")
	}
	if d.errorMsg != "Condition not true" {
		t.Error("Expected dummyT.errorMsg to be 'Condition not true'")
	}

}

func TestAssertFalse(t *testing.T) {

	// success case
	d := &dummyT{}
	AssertFalse("Condition not false", false, d)

	if d.fail {
		t.Error("Expected dummyT.fail to be false")
	}
	if d.errorMsg != "" {
		t.Error("Expected dummyT.errorMsg to be blank")
	}

	// error case
	d = &dummyT{}
	AssertFalse("Condition not false", true, d)

	if !d.fail {
		t.Error("Expected dummyT.fail to be true")
	}
	if d.errorMsg != "Condition not false" {
		t.Error("Expected dummyT.errorMsg to be 'Condition not false'")
	}

}

func TestAssertEquals(t *testing.T) {

	// success case
	d := &dummyT{}
	AssertEquals("Equality failed", 1, 1, d)

	if d.fail {
		t.Error("Expected dummyT.fail to be false")
	}
	if d.errorMsg != "" {
		t.Error("Expected dummyT.errorMsg to be blank")
	}

	// error case
	d = &dummyT{}
	AssertEquals("Equality failed", 1, 2, d)

	if !d.fail {
		t.Error("Expected dummyT.fail to be true")
	}

	expectedErrorMsg := "'Equality failed'. Expected: '1'. Actual: '2'"
	if d.errorMsg != expectedErrorMsg {
		t.Error("Expected dummyT.errorMsg to be '" + expectedErrorMsg + "'")
	}

}

func TestAssertPanic_for_success_case(t *testing.T) {

	d := &dummyT{}

	defer func() {
		if d.fail {
			t.Error("Expected dummyT.fail to be false")
		}
		if d.errorMsg != "" {
			t.Error("Expected dummyT.errorMsg to be blank")
		}
	}()

	defer AssertPanic("Expected panic", d)

	panic("Intentionally triggered panic")
}

func TestAssertPanic_for_error_case(t *testing.T) {

	d := &dummyT{}

	defer func() {
		if !d.fail {
			t.Error("Expected dummyT.fail to be true")
		}
		if d.errorMsg != "Expected panic" {
			t.Error("Expected dummyT.errorMsg to be 'Expected panic'")
		}
	}()

	defer AssertPanic("Expected panic", d)

	// do nothing. assure no panics
}

// dummy tester
type dummyT struct {
	fail     bool
	errorMsg string
}

func (t *dummyT) Error(args ...interface{}) {
	errorMsg := ""
	for idx, arg := range args {
		if idx > 0 {
			errorMsg += "; "
		}
		errorMsg += fmt.Sprintf("%v", arg)
	}
	t.fail = true
	t.errorMsg = errorMsg
}

func (t *dummyT) Errorf(format string, args ...interface{}) {
	t.fail = true
	t.errorMsg = fmt.Sprintf(format, args...)
}
