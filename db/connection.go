package db

import "database/sql"

// Connection abstracts out the relevant operations common to sql.DB and
// sql.Tx into a base interface
type Connection interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// ----------
// BASIC CRUD
// ----------

// CreateOne new row in DB
func CreateOne(conn Connection, insertCommand string, insertArgs []interface{}, query string, dest []interface{}) Error {

	// insert
	insert, insertError := conn.Exec(insertCommand, insertArgs...)
	if insertError != nil {
		return WrapError(insertError)
	}

	// grab last inserted id
	id, idRetrievalError := insert.LastInsertId()
	if idRetrievalError != nil {
		return WrapError(idRetrievalError)
	}

	// return inserted data
	return WrapError(LookupOne(conn, query, []interface{}{id}, dest))
}

// LookupOne row in DB
func LookupOne(conn Connection, query string, args []interface{}, dest []interface{}) Error {
	row := conn.QueryRow(query, args...)
	scanError := row.Scan(dest...)
	switch scanError {
	case sql.ErrNoRows:
		return NewNotFoundError("")
	case nil:
		return nil
	default:
		return NewGenericError(scanError.Error())
	}
}

// UpdateOne row in DB
func UpdateOne(conn Connection, id string, updateCommand string, updateArgs []interface{}, query string, dest []interface{}) Error {

	// update
	_, updateError := conn.Exec(updateCommand, updateArgs...)
	if updateError != nil {
		return WrapError(updateError)
	}

	// look up updated data using id
	return LookupOne(conn, query, []interface{}{id}, dest)
}

// DeleteOne row in DB
func DeleteOne(conn Connection, id string, deleteCommand string, query string, dest []interface{}) Error {

	// look up data being deleted
	lookupError := LookupOne(conn, query, []interface{}{id}, dest)
	if lookupError != nil {
		return lookupError
	}

	// update
	_, deleteError := conn.Exec(deleteCommand, id)
	if deleteError != nil {
		return WrapError(deleteError)
	}

	return nil
}
