package db

import "fmt"

// ---
// Error Type Definition
// ---

// Error creates Database specific errors
type Error interface {
	error
	Type() string
}

type databaseError struct {
	errorType   string
	errorDetail string
}

// Type returns error type
func (e *databaseError) Type() string {
	return e.errorType
}

// Error satisfies Go's built in error interface
func (e *databaseError) Error() string {
	return fmt.Sprintf("%s", e.errorDetail)
}

// ---
// Possible Error Types
// ---

// NotFound - errors where provided ID does not exist
var NotFound = "NotFound"

// NewNotFoundError from detail
func NewNotFoundError(detail string) Error {
	return &databaseError{errorType: NotFound, errorDetail: detail}
}

// BadRequest - errors where input criteria is invalid
var BadRequest = "BadRequest"

// NewBadRequestError from detail
func NewBadRequestError(detail string) Error {
	return &databaseError{errorType: BadRequest, errorDetail: detail}
}

// GenericError - uncategorized errors
var GenericError = "Other"

// NewGenericError from detail
func NewGenericError(detail string) Error {
	return &databaseError{errorType: "Other", errorDetail: detail}
}

// Forbidden - errors where operation is forbidden
var Forbidden = "Forbidden"

// NewForbiddenError from detail
func NewForbiddenError(detail string) Error {
	return &databaseError{errorType: Forbidden, errorDetail: detail}
}

// WrapError (raw nullable errors) into db.Error
func WrapError(wrapped error) Error {
	if wrapped == nil {
		return nil
	}
	return &databaseError{errorType: "Other", errorDetail: wrapped.Error()}
}
