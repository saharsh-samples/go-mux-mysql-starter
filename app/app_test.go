package app

import (
	"errors"
	"syscall"
	"testing"

	"github.com/saharsh-samples/go-mux-sql-starter/test"
)

// happy server
type happyServer struct{ ready bool }

func (server *happyServer) Run()               { server.ready = true }
func (server *happyServer) IsReady() bool      { return server.ready }
func (server *happyServer) Port() (int, error) { return 8080, nil }
func (server *happyServer) Shutdown() error {
	server.ready = false
	return nil
}

// lazy server
type lazyServer struct{ happyServer }

func (server *lazyServer) Run() { /* no-op */ }

// dead server
type badShutdownServer struct{ happyServer }

func (server *badShutdownServer) Shutdown() error {
	server.ready = false
	return errors.New("Simulated error")
}

// ---
// tests
// ---

func TestHappyPath(t *testing.T) {

	hooksRan := false

	ctx := Bootstrap(&ContextIn{
		StartupTimeoutInSeconds: 1,
		HTTPServer:              &happyServer{},
		ShutdownHooks:           []ShutdownHook{func() { hooksRan = true }},
	})

	// start server
	go func() {
		test.AssertEquals("", TerminatedStatus, ctx.App.Run().Status, t)
	}()

	// first status is initializing with no detail
	appStatus := <-ctx.Status
	test.AssertEquals("", InitializingStatus, appStatus.Status, t)
	test.AssertEquals("", "", appStatus.Detail, t)

	// second status is ready with server port as detail
	appStatus = <-ctx.Status
	test.AssertEquals("", ReadyStatus, appStatus.Status, t)
	test.AssertEquals("", "8080", appStatus.Detail, t)

	// send SIGTERM to app
	ctx.Signal <- syscall.SIGTERM

	// last status is terminated with no detail
	appStatus = <-ctx.Status
	test.AssertEquals("", TerminatedStatus, appStatus.Status, t)
	test.AssertEquals("", "", appStatus.Detail, t)

	// make sure status channel is now closed
	_, open := <-ctx.Status
	test.AssertFalse("Expected app status channel to close", open, t)

	// make sure hooks ran
	test.AssertTrue("Expected shutdown hooks to run", hooksRan, t)
}

func TestLazyPath(t *testing.T) {

	hooksRan := false

	ctx := Bootstrap(&ContextIn{
		StartupTimeoutInSeconds: 1,
		HTTPServer:              &lazyServer{},
		ShutdownHooks:           []ShutdownHook{func() { hooksRan = true }},
	})

	// start server
	go func() {
		status := ctx.App.Run()
		test.AssertEquals("", ErrorStatus, status.Status, t)
		test.AssertEquals("", "Timed out waiting for server to start", status.Detail, t)
	}()

	// first status is initializing with no detail
	appStatus := <-ctx.Status
	test.AssertEquals("", InitializingStatus, appStatus.Status, t)
	test.AssertEquals("", "", appStatus.Detail, t)

	// second status is error with timeout error as detail
	appStatus = <-ctx.Status
	test.AssertEquals("", ErrorStatus, appStatus.Status, t)
	test.AssertEquals("", "Timed out waiting for server to start", appStatus.Detail, t)

	// make sure status channel is now closed
	_, open := <-ctx.Status
	test.AssertFalse("Expected app status channel to close", open, t)

	// make sure hooks ran
	test.AssertTrue("Expected shutdown hooks to run", hooksRan, t)

}

func TestBadShutdownPath(t *testing.T) {

	hooksRan := false

	ctx := Bootstrap(&ContextIn{
		StartupTimeoutInSeconds: 1,
		HTTPServer:              &badShutdownServer{},
		ShutdownHooks:           []ShutdownHook{func() { hooksRan = true }},
	})

	// start server
	go func() {
		status := ctx.App.Run()
		test.AssertEquals("", ErrorStatus, status.Status, t)
		test.AssertEquals("", "Simulated error", status.Detail, t)
	}()

	// first status is initializing with no detail
	appStatus := <-ctx.Status
	test.AssertEquals("", InitializingStatus, appStatus.Status, t)
	test.AssertEquals("", "", appStatus.Detail, t)

	// second status is ready with server port as detail
	appStatus = <-ctx.Status
	test.AssertEquals("", ReadyStatus, appStatus.Status, t)
	test.AssertEquals("", "8080", appStatus.Detail, t)

	// send SIGTERM to app
	ctx.Signal <- syscall.SIGTERM

	// last status is error with reported error as detail
	appStatus = <-ctx.Status
	test.AssertEquals("", ErrorStatus, appStatus.Status, t)
	test.AssertEquals("", "Simulated error", appStatus.Detail, t)

	// make sure status channel is now closed
	_, open := <-ctx.Status
	test.AssertFalse("Expected app status channel to close", open, t)

	// make sure hooks ran
	test.AssertTrue("Expected shutdown hooks to run", hooksRan, t)

}
