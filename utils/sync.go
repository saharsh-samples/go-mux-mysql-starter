package utils

import (
	"errors"
	"time"
)

// ErrTimeout is the error text of timeout error from WaitTill
const ErrTimeout = "Timed out waiting"

// WaitTill provided condition returns true or timeout occurs
func WaitTill(condition func() bool, timeoutInSeconds int) error {

	waitTill := time.Now().Add(time.Duration(timeoutInSeconds) * time.Second)

	for !condition() {
		if time.Now().After(waitTill) {
			return errors.New(ErrTimeout)
		}
		time.Sleep(10 * time.Millisecond)
	}

	return nil
}
