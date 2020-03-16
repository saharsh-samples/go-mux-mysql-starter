package db

import "database/sql"

// ContextIn describes dependecies needed by this package
type ContextIn struct {
	DatabaseHandle *sql.DB
}

// ContextOut describes dependencies exported by this package
type ContextOut struct {
	Database Database
}

// Bootstrap initializes this module with ContextIn and exports
// resulting ContextOut
func Bootstrap(in *ContextIn) *ContextOut {

	// create and export out context
	out := &ContextOut{}
	out.Database = &database{dbHandle: in.DatabaseHandle}

	return out
}
