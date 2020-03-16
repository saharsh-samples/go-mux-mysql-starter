package routes

import (
	"github.com/saharsh-samples/go-mux-sql-starter/http"
)

// ContextIn describes dependecies needed by this package
type ContextIn struct {
	// Add external dependencies here
}

// ContextOut describes dependencies exported by this package
type ContextOut struct {
	RoutesToRegister []http.Routes
}

// Bootstrap initializes this module with ContextIn and exports
// resulting ContextOut
func Bootstrap(in *ContextIn) *ContextOut {

	out := &ContextOut{}
	out.RoutesToRegister = []http.Routes{
		&LivenessCheck{},
		// Add exported routes here
	}

	return out
}
