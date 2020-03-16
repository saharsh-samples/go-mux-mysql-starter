package test

// T encapsulates relevant testing.T operations used by this package
type T interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

// AssertTrue asserts given condition is true
func AssertTrue(message string, condition bool, t T) {
	if !condition {
		t.Error(message)
	}
}

// AssertFalse asserts given condition is false
func AssertFalse(message string, condition bool, t T) {
	AssertTrue(message, !condition, t)
}

// AssertEquals asserts expected and actual are equal using '!='
func AssertEquals(message string, expected interface{}, actual interface{}, t T) {
	if expected != actual {
		t.Errorf("'%v'. Expected: '%v'. Actual: '%v'", message, expected, actual)
	}
}

// AssertPanic asserts panic occurred before this function call. Expected to be called via 'defer'
func AssertPanic(message string, t T) {
	if r := recover(); r == nil {
		t.Error(message)
	}
}
