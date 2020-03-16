package utils

import (
	"fmt"
	"strings"
	"time"
)

const (

	// DateFormat is 'YYYY-MM-DD'
	DateFormat = "2006-01-02"
)

// JSONDate is marshalled as string in 'YYYY-MM-DD' format and
// unmarshalled as time.Time instance
type JSONDate time.Time

// UnmarshalJSON string in 'YYYY-MM-DD' format into JSONDate
func (j *JSONDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(DateFormat, s)
	if err != nil {
		return err
	}
	*j = JSONDate(t)
	return nil
}

// MarshalJSON string in 'YYYY-MM-DD' format from JSONDate
func (j *JSONDate) MarshalJSON() ([]byte, error) {
	asTime := time.Time(*j)
	return []byte(fmt.Sprint("\"", asTime.Format(DateFormat), "\"")), nil
}
