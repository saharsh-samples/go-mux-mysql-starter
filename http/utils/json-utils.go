package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/saharsh-samples/go-mux-sql-starter/db"
)

// JSONBody is base type for types that represent JSON data
type JSONBody interface {
	Validate() error
}

// AlwaysValidJSON can be embedded by JSONBody types
// that don't need to validate
type AlwaysValidJSON struct{}

// Validate always returns no error for AlwaysValidJSON
func (body *AlwaysValidJSON) Validate() error {
	return nil
}

// PagedResponse is the container for paginated responses
type PagedResponse struct {
	Limit   int
	Offset  int
	Total   int64
	Payload []interface{}

	// response only struct, so validation is unnecessary
	AlwaysValidJSON
}

// JSONUtils can be used to delegate JSON specific concerns
type JSONUtils interface {

	// serialization/deserialization
	SetJSONResponse(w http.ResponseWriter, statusCode int, body interface{})
	ParseJSONRequest(r *http.Request, value JSONBody, w http.ResponseWriter) error
	Unmarshal(r *http.Request, value interface{}, w http.ResponseWriter) error

	// error handling support
	BadRequest(w http.ResponseWriter, detail string)
	Unauthorized(w http.ResponseWriter, detail string)
	Forbidden(w http.ResponseWriter, detail string)
	NotFound(w http.ResponseWriter, detail string)
	InternalError(w http.ResponseWriter, detail string)
	HandleDatabaseError(w http.ResponseWriter, err db.Error)
}

type jsonUtils struct{}

//---------------
// Error messages
//---------------

// ErrorMessage model
type ErrorMessage struct {
	StatusCode int
	Message    string
	Detail     string
}

func (jsonUtils *jsonUtils) setErrorResponse(w http.ResponseWriter, errorMsg *ErrorMessage) {
	fmt.Printf("ErrorMessage: '%v'\n", *errorMsg)
	jsonUtils.SetJSONResponse(w, errorMsg.StatusCode, errorMsg)
}

// BadRequest will set response header and body to indicate Bad Request error
func (jsonUtils *jsonUtils) BadRequest(w http.ResponseWriter, detail string) {
	jsonUtils.setErrorResponse(w, &ErrorMessage{http.StatusBadRequest, "Bad Request", detail})
}

// Unauthorized will set response header and body to indicate Unauthorized error
func (jsonUtils *jsonUtils) Unauthorized(w http.ResponseWriter, detail string) {
	jsonUtils.setErrorResponse(w, &ErrorMessage{http.StatusUnauthorized, "Unauthorized", detail})
}

// Forbidden will set response header and body to indicate Forbidden error
func (jsonUtils *jsonUtils) Forbidden(w http.ResponseWriter, detail string) {
	jsonUtils.setErrorResponse(w, &ErrorMessage{http.StatusForbidden, "Forbidden", detail})
}

// NotFound will set response header and body to indicate Not Found error
func (jsonUtils *jsonUtils) NotFound(w http.ResponseWriter, detail string) {
	jsonUtils.setErrorResponse(w, &ErrorMessage{http.StatusNotFound, "Not Found", detail})
}

// InternalError will set response header and body to indicate ISE
func (jsonUtils *jsonUtils) InternalError(w http.ResponseWriter, detail string) {
	jsonUtils.setErrorResponse(w, &ErrorMessage{http.StatusInternalServerError, "Internal Server Error", detail})
}

// HandleDatabaseError cetralizes logic to process database errors
func (jsonUtils *jsonUtils) HandleDatabaseError(w http.ResponseWriter, err db.Error) {
	if err.Type() == db.BadRequest {
		jsonUtils.BadRequest(w, err.Error())
	} else if err.Type() == db.NotFound {
		jsonUtils.NotFound(w, err.Error())
	} else {
		jsonUtils.InternalError(w, err.Error())
	}
}

// ----------------------------------
// JSON serialization/deserialization
// ----------------------------------

// SetJSONResponse is used to serialize given struct to a JSON object
func (jsonUtils *jsonUtils) SetJSONResponse(w http.ResponseWriter, statusCode int, body interface{}) {

	bodyJSON, marshalError := json.Marshal(body)
	if marshalError != nil {
		jsonUtils.InternalError(w, marshalError.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(bodyJSON)
}

// ParseJSONRequest is used to parse request body as a JSON object
// the JSON object is deserialized into provided struct
func (jsonUtils *jsonUtils) ParseJSONRequest(r *http.Request, value JSONBody, w http.ResponseWriter) error {

	decoder := json.NewDecoder(r.Body)

	unmarshalError := decoder.Decode(value)
	if unmarshalError != nil {
		jsonUtils.BadRequest(w, "Malformed JSON body")
		return unmarshalError
	}

	validationError := value.Validate()
	if validationError != nil {
		jsonUtils.BadRequest(w, validationError.Error())
		return validationError
	}

	return nil

}

// Unmarshal is a WIP function to allow unmarshalling bodies not compliant yet with JSONBody
func (jsonUtils *jsonUtils) Unmarshal(r *http.Request, value interface{}, w http.ResponseWriter) error {

	decoder := json.NewDecoder(r.Body)

	unmarshalError := decoder.Decode(value)
	if unmarshalError != nil {
		jsonUtils.BadRequest(w, "Malformed JSON body")
		return unmarshalError
	}

	return nil
}
