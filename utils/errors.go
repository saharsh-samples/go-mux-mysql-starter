package utils

import (
	"fmt"
)

// JoinErrors concatenates error text of all errors into one semi-colon delimited string
func JoinErrors(errors ...error) string {
	joinedErrorString := ""
	for _, error := range errors {
		if error != nil {
			if joinedErrorString == "" {
				joinedErrorString = error.Error()
			} else {
				joinedErrorString = fmt.Sprintf("%v; %v", joinedErrorString, error.Error())
			}
		}
	}
	return joinedErrorString
}
