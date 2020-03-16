package test

import (
	"errors"
	"testing"

	"github.com/saharsh-samples/go-mux-sql-starter/test"
)

func Test_MockResponseWriter_Header(t *testing.T) {

	// Arrange
	responseWriter := &MockResponseWriter{HeaderReturn: make(map[string][]string)}

	// Act
	responseWriter.Header().Set("Content-Type", "applicaton/json")

	// Assert
	test.AssertEquals("", "applicaton/json", responseWriter.HeaderReturn.Get("Content-Type"), t)

}

func Test_MockResponseWriter_Write(t *testing.T) {

	// Arrange
	responseWriter := &MockResponseWriter{WriteIntReturn: 11, WriteErrorReturn: errors.New("Some Error")}

	// Act
	intReturn, errorReturn := responseWriter.Write([]byte("cheese"))

	// Asssert
	test.AssertEquals("", 11, intReturn, t)
	test.AssertEquals("", "Some Error", errorReturn.Error(), t)
	test.AssertEquals("", "cheese", string(responseWriter.WriteBytesArg), t)
}

func Test_MockResponseWriter_WriteHeader(t *testing.T) {

	// Arrange
	responseWriter := &MockResponseWriter{}

	// Act
	responseWriter.WriteHeader(401)

	// Assert
	test.AssertEquals("", 401, responseWriter.WriteHeaderStatusCodeArg, t)
}
