package db

import (
	"database/sql"
)

// Database should be used as the highest level abstraction of the database
type Database interface {
	Close()
	GetConnection() Connection
	WithTransaction(wrapped func(Connection) Error) (txExecError Error)
	CreateOne(insertCommand string, insertArgs []interface{}, query string, dest []interface{}) Error
	LookupOne(query string, args []interface{}, dest []interface{}) Error
	UpdateOne(id interface{}, updateCommand string, updateArgs []interface{}, query string, dest []interface{}) Error
	DeleteOne(id interface{}, deleteCommand string, query string, dest []interface{}) Error
}

type database struct {
	dbHandle *sql.DB
}

// Close closes the underlying database handle.
// Should only be called before app termination
func (database *database) Close() {
	err := database.dbHandle.Close()
	if err != nil {
		panic(err.Error())
	}
}

// GetConnection to run database commands directly
func (database *database) GetConnection() Connection {
	return database.dbHandle
}

// ---
// TRANSACTION SUPPORT
// ---

// WithTransaction creates a new transaction and handles rollback/commit
// based on the error object returned by the wrapped code
func (database *database) WithTransaction(wrapped func(Connection) Error) (txExecError Error) {

	tx, txBeginError := database.dbHandle.Begin()
	if txBeginError != nil {
		return WrapError(txBeginError)
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			tx.Rollback()
			panic(p)
		} else if txExecError != nil {
			// something went wrong, rollback
			tx.Rollback()
		} else {
			// all good, commit
			txExecError = WrapError(tx.Commit())
		}
	}()

	return wrapped(tx)
}

// ---
// WRAPPERS TO BIND BASIC CRUD TO database TYPE
//
// NOTE: This is done as a convenience for when these operations are done
// standalone and not as part of a transaction.
//
// DO NOT USE THESE INSIDE A TRANSACTION
// ---

// CreateOne row in DB
func (database *database) CreateOne(insertCommand string, insertArgs []interface{}, query string, dest []interface{}) Error {
	return CreateOne(database.dbHandle, insertCommand, insertArgs, query, dest)
}

// LookupOne row in DB
func (database *database) LookupOne(query string, args []interface{}, dest []interface{}) Error {
	return LookupOne(database.dbHandle, query, args, dest)
}

// UpdateOne row in DB
func (database *database) UpdateOne(id interface{}, updateCommand string, updateArgs []interface{}, query string, dest []interface{}) Error {
	return UpdateOne(database.dbHandle, id, updateCommand, updateArgs, query, dest)
}

// DeleteOne row in DB
func (database *database) DeleteOne(id interface{}, deleteCommand string, query string, dest []interface{}) Error {
	return DeleteOne(database.dbHandle, id, deleteCommand, query, dest)
}
