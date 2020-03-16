package app

import (
	"fmt"
	"os"

	"github.com/saharsh-samples/go-mux-sql-starter/http"
	"github.com/saharsh-samples/go-mux-sql-starter/utils"
)

// App interface
type App interface {
	Run() Status
}

// Status of application
type Status struct {
	Status string
	Detail string
}

const (

	// InitializingStatus means app is still intializing
	InitializingStatus = "Initializing"

	// ReadyStatus means app is fully initialized and running
	ReadyStatus = "Ready"

	// TerminatedStatus means app gracefully terminated
	TerminatedStatus = "Terminated"

	// ErrorStatus indicates app encountered an error
	ErrorStatus = "Error"
)

// ShutdownHook function type
type ShutdownHook func()

type app struct {
	startupTimeoutInSeconds int
	server                  http.Server
	shutdownHooks           []ShutdownHook
	sigs                    <-chan os.Signal
	status                  chan<- Status
}

// Run the application
func (app *app) Run() Status {

	defer close(app.status)

	app.status <- Status{Status: InitializingStatus}

	// register shutdown hooks
	for _, shutdownHook := range app.shutdownHooks {
		defer shutdownHook()
	}

	// run http server
	go app.server.Run()

	// make sure server starts
	utils.WaitTill(func() bool {
		return app.server.IsReady()
	}, app.startupTimeoutInSeconds)

	if app.server.IsReady() {
		fmt.Println("Startup successful. App running...")
		port, _ := app.server.Port()
		app.status <- Status{Status: ReadyStatus, Detail: fmt.Sprintf("%d", port)}
	} else {
		status := Status{Status: ErrorStatus, Detail: "Timed out waiting for server to start"}
		app.status <- status
		return status
	}

	// wait for termination signal
	sig := <-app.sigs
	fmt.Println("Recieved ", sig.String(), " signal. Terminating...")

	// shutdown http server
	err := app.server.Shutdown()

	// exit communicating error (if any)
	if err != nil {
		status := Status{Status: ErrorStatus, Detail: err.Error()}
		app.status <- status
		return status
	}
	status := Status{Status: TerminatedStatus}
	app.status <- status
	return status
}
