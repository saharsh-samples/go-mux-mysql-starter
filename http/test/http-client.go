package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/saharsh-samples/go-mux-sql-starter/test"
)

// Get request expecting JSON response
func Get(t test.T, url string, respBody interface{}, apiToken string) int {
	return doRequestWithoutBody(t, url, http.MethodGet, respBody, apiToken)
}

// GetWithoutBody does GET request only expecting status code
func GetWithoutBody(t test.T, url string, apiToken string) int {
	return doRequestWithoutBody(t, url, http.MethodGet, nil, apiToken)
}

// Put request with JSON body expecting JSON response
func Put(t test.T, url string, reqBody interface{}, respBody interface{}, apiToken string) int {
	return doRequest(t, url, http.MethodPut, reqBody, respBody, apiToken)
}

// Post request with JSON body expecting JSON response
func Post(t test.T, url string, reqBody interface{}, respBody interface{}, apiToken string) int {
	return doRequest(t, url, http.MethodPost, reqBody, respBody, apiToken)
}

// Delete request expecting JSON response
func Delete(t test.T, url string, respBody interface{}, apiToken string) int {
	return doRequestWithoutBody(t, url, http.MethodDelete, respBody, apiToken)
}

// DeleteWithoutBody does DELETE request only expecting status code
func DeleteWithoutBody(t test.T, url string, apiToken string) int {
	return doRequestWithoutBody(t, url, http.MethodDelete, nil, apiToken)
}

func doRequestWithoutBody(t test.T, url string, httpMethod string, respBody interface{}, apiToken string) int {
	return doRequest(t, url, httpMethod, nil, respBody, apiToken)
}

func doRequest(t test.T, url string, httpMethod string, reqBody interface{}, respBody interface{}, apiToken string) int {

	// marshal request body
	var reqBodyReader io.Reader
	if reqBody != nil {
		marshalled, jsonErr := json.Marshal(reqBody)
		if jsonErr != nil {
			t.Error(jsonErr)
		}
		reqBodyReader = bytes.NewReader(marshalled)
	}

	client := &http.Client{}
	request, requestErr := http.NewRequest(httpMethod, url, reqBodyReader)
	if requestErr != nil {
		t.Error(requestErr)
	}

	request.Header.Add("Content-Type", "application/json")
	if apiToken != "" {
		request.Header.Add("Authorization", "Bearer "+apiToken)
	}

	response, callErr := client.Do(request)
	if callErr != nil {
		t.Error(callErr)
	}
	defer response.Body.Close()

	if respBody != nil {
		unmarshallErr := json.NewDecoder(response.Body).Decode(respBody)
		if unmarshallErr != nil {
			t.Error(unmarshallErr)
		}
	}

	return response.StatusCode
}
