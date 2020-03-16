package app

import (
	"os"

	"github.com/saharsh-samples/go-mux-sql-starter/http"
)

// ContextIn describes dependecies needed by this package
type ContextIn struct {
	StartupTimeoutInSeconds int
	HTTPServer              http.Server
	ShutdownHooks           []ShutdownHook
}

// ContextOut describes dependencies exported by this package
type ContextOut struct {
	App    App
	Signal chan<- os.Signal
	Status <-chan Status
}

// Bootstrap initializes this module with ContextIn and exports
// resulting ContextOut
func Bootstrap(in *ContextIn) *ContextOut {

	// channels
	signal := make(chan os.Signal, 1)
	status := make(chan Status, 3)

	// context out
	out := &ContextOut{}
	out.App = &app{
		startupTimeoutInSeconds: in.StartupTimeoutInSeconds,
		server:                  in.HTTPServer,
		shutdownHooks:           in.ShutdownHooks,
		sigs:                    signal,
		status:                  status,
	}
	out.Signal = signal
	out.Status = status

	return out
}
