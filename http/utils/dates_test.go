package utils

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/saharsh-samples/go-mux-sql-starter/test"
)

type container struct {
	Date JSONDate
}

func TestUnmarshal_of_JSONDate(t *testing.T) {

	// Arrange
	jsonBody := json.RawMessage(`{"Date":"2020-04-01"}`)

	// Act
	unmarshalled := &container{}
	err := json.Unmarshal(jsonBody, unmarshalled)

	// Assert
	test.AssertTrue("Expected err to be nil", err == nil, t)
	test.AssertEquals("", "2020-04-01", time.Time(unmarshalled.Date).Format(DateFormat), t)

}

func TestUnmarshal_of_invalid_JSONDate(t *testing.T) {

	unmarshalled := &container{}
	err := json.Unmarshal(json.RawMessage(`{"Date":"2020/04/01"}`), unmarshalled)
	test.AssertFalse("Expected err to be non-nil", err == nil, t)
	test.AssertTrue("Expected 'cannot parse' in error message", strings.Contains(err.Error(), "cannot parse"), t)

	unmarshalled = &container{}
	err = json.Unmarshal(json.RawMessage(`{"Date":"2020-04-01T12:30:07.123Z"}`), unmarshalled)
	test.AssertFalse("Expected err to be non-nil", err == nil, t)
	test.AssertTrue("Expected 'extra text' in error message", strings.Contains(err.Error(), "extra text"), t)

	unmarshalled = &container{}
	err = json.Unmarshal(json.RawMessage(`{"Date":"something-strange"}`), unmarshalled)
	test.AssertFalse("Expected err to be non-nil", err == nil, t)
	test.AssertTrue("Expected 'cannot parse' in error message", strings.Contains(err.Error(), "cannot parse"), t)

}

func TestMarshal_of_JSONDate(t *testing.T) {

	// Arrange
	asTime, _ := time.Parse(DateFormat, "2020-04-01")
	unmarshalled := &container{Date: JSONDate(asTime)}

	// Act
	jsonBody, err := json.Marshal(unmarshalled)

	// Assert
	test.AssertTrue("Expected err to be nil", err == nil, t)
	test.AssertEquals("", `{"Date":"2020-04-01"}`, string(jsonBody), t)

}
