package test

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/saharsh-samples/go-mux-sql-starter/db"
)

// NewDatabaseWithMockConnection creates a real db.Database instance
// with mocked driver using github.com/DATA-DOG/go-sqlmock
func NewDatabaseWithMockConnection(t *testing.T) (db.Database, sqlmock.Sqlmock, func()) {

	handle, mock, openErr := sqlmock.New()
	if openErr != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", openErr)
	}

	database := db.Bootstrap(&db.ContextIn{DatabaseHandle: handle}).Database
	return database, mock, func() {
		mock.ExpectClose()
		database.Close()
	}
}

// MockDatabase is a test friendly implementation of db.Database
// Use this struct to mock out the basic CRUD methods bound to
// db.Database. To test transactions or mock out the underlying
// db.Connection, use NewDatabaseWithMockConnection and
// github.com/DATA-DOG/go-sqlmock methods
type MockDatabase struct {

	// GetConnection
	GetConnectionReturn db.Connection

	// Close
	CloseCalled bool

	// CreateOne
	CreateOneCommand    string
	CreateOneArgs       []interface{}
	CreateOneQuery      string
	CreateOneDestWriter func([]interface{})
	CreateOneError      db.Error

	// LookupOne
	LookupOneQuery      string
	LookupOneArgs       []interface{}
	LookupOneDestWriter func([]interface{})
	LookupOneError      db.Error

	// UpdateOne
	UpdateOneID         string
	UpdateOneCommand    string
	UpdateOneArgs       []interface{}
	UpdateOneQuery      string
	UpdateOneDestWriter func([]interface{})
	UpdateOneError      db.Error

	// DeleteOne
	DeleteOneID         string
	DeleteOneCommand    string
	DeleteOneQuery      string
	DeleteOneDestWriter func([]interface{})
	DeleteOneError      db.Error
}

// GetConnection to run database commands directly
func (database *MockDatabase) GetConnection() db.Connection {
	return database.GetConnectionReturn
}

// Close the database handle gracefully
func (database *MockDatabase) Close() {
	database.CloseCalled = true
}

// WithTransaction creates a new transaction and handles rollback/commit
// based on the error object returned by the wrapped code
func (database *MockDatabase) WithTransaction(wrapped func(db.Connection) db.Error) (txExecError db.Error) {
	// TODO
	return nil
}

// CreateOne row in DB
func (database *MockDatabase) CreateOne(command string, args []interface{}, query string, dest []interface{}) db.Error {
	database.CreateOneCommand = command
	database.CreateOneArgs = args
	database.CreateOneQuery = query
	if database.CreateOneDestWriter != nil {
		database.CreateOneDestWriter(dest)
	}
	return database.CreateOneError
}

// LookupOne row in DB
func (database *MockDatabase) LookupOne(query string, args []interface{}, dest []interface{}) db.Error {
	database.LookupOneQuery = query
	database.LookupOneArgs = args
	if database.LookupOneDestWriter != nil {
		database.LookupOneDestWriter(dest)
	}
	return database.LookupOneError
}

// UpdateOne row in DB
func (database *MockDatabase) UpdateOne(id string, updateCommand string, updateArgs []interface{}, query string, dest []interface{}) db.Error {
	database.UpdateOneID = id
	database.UpdateOneCommand = updateCommand
	database.UpdateOneArgs = updateArgs
	database.UpdateOneQuery = query
	if database.UpdateOneDestWriter != nil {
		database.UpdateOneDestWriter(dest)
	}
	return database.UpdateOneError
}

// DeleteOne row in DB
func (database *MockDatabase) DeleteOne(id string, deleteCommand string, query string, dest []interface{}) db.Error {
	database.DeleteOneID = id
	database.DeleteOneCommand = deleteCommand
	database.DeleteOneQuery = query
	if database.DeleteOneDestWriter != nil {
		database.DeleteOneDestWriter(dest)
	}
	return database.DeleteOneError
}
