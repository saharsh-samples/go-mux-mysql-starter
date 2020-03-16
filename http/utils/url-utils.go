package utils

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// URLUtils can be used to extract information from URLs
type URLUtils interface {
	GetPathParams(r *http.Request) map[string]string
	GetQueryParameterAsString(r *http.Request, queryParamName string, defaultValue *string) (string, error)
	GetQueryParameterAsInteger(r *http.Request, queryParamName string, defaultValue *int) (int, error)
	GetQueryParameterAsArrayOfString(r *http.Request, queryParamName string, minimumRequired int) ([]string, error)
}

type urlUtils struct{}

// GetPathParams returns path params as key value pairs
func (urlUtils *urlUtils) GetPathParams(r *http.Request) map[string]string {
	return mux.Vars(r)
}

// GetQueryParameterAsString extracts 'string' value for specified query param
func (urlUtils *urlUtils) GetQueryParameterAsString(r *http.Request, queryParamName string, defaultValue *string) (string, error) {

	values := r.URL.Query()[queryParamName]
	numOfValues := len(values)

	if numOfValues == 0 && defaultValue != nil {
		return *defaultValue, nil
	}

	if numOfValues != 1 {
		return "", fmt.Errorf("Expected ONE and only ONE value for '%v' Query Parameter", queryParamName)
	}

	return values[0], nil

}

// GetQueryParameterAsInteger extracts 'int' value for specified query param
func (urlUtils *urlUtils) GetQueryParameterAsInteger(r *http.Request, queryParamName string, defaultValue *int) (int, error) {

	values := r.URL.Query()[queryParamName]
	numOfValues := len(values)

	if numOfValues == 0 && defaultValue != nil {
		return *defaultValue, nil
	}

	if numOfValues != 1 {
		return 0, fmt.Errorf("Expected ONE and only ONE value for '%v' Query Parameter", queryParamName)
	}

	value, err := strconv.Atoi(values[0])
	if err != nil {
		return 0, fmt.Errorf("Error parsing query parameter '%v=%v' as an integer", queryParamName, values[0])
	}

	return value, nil

}

// GetQueryParameterAsArrayOfString extracts '[]string' value for specified query param
func (urlUtils *urlUtils) GetQueryParameterAsArrayOfString(r *http.Request, queryParamName string, minimumRequired int) ([]string, error) {

	values := r.URL.Query()[queryParamName]
	numOfValues := len(values)

	if numOfValues < minimumRequired {
		return nil, fmt.Errorf(
			"Expected %v value(s) but only %v value(s) found for query parameter '%v'",
			minimumRequired, numOfValues, queryParamName)
	}

	return values, nil
}
