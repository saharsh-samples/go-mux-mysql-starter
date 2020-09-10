package test

import (
	"fmt"
	"strconv"
	"syscall"

	"github.com/saharsh-samples/go-mux-sql-starter/app"
	base "github.com/saharsh-samples/go-mux-sql-starter/http"
	"github.com/saharsh-samples/go-mux-sql-starter/test"
)

// StartTestServer with specified routes. Returns chosen network port
func StartTestServer(routes []base.Routes, t test.T) (int, *app.ContextOut) {
	return StartTestServerWithMiddleWare(nil, routes, t)
}

// StartTestServerWithMiddleWare with specified middlewares and routes. Returns chosen network port
func StartTestServerWithMiddleWare(middlewares base.Middlewares, routes []base.Routes, t test.T) (int, *app.ContextOut) {

	// bootstrap
	server := base.Bootstrap(&base.ContextIn{
		Port:                  0,
		MiddlewaresToRegister: middlewares,
		RoutesToRegister:      routes,
	}).Server

	appCtx := app.Bootstrap(&app.ContextIn{
		StartupTimeoutInSeconds: 5,
		HTTPServer:              server,
		ShutdownHooks:           []app.ShutdownHook{},
	})

	// Run test app in another goroutine
	go appCtx.App.Run()
	appStatus := <-appCtx.Status
	appStatus = <-appCtx.Status
	test.AssertEquals("", app.ReadyStatus, appStatus.Status, t)

	port, portErr := strconv.Atoi(appStatus.Detail)
	test.AssertTrue(fmt.Sprintf("Error getting server port: %v", portErr), portErr == nil, t)
	return port, appCtx
}

// StopTestServer gracefully
func StopTestServer(appCtx *app.ContextOut, t test.T) {
	appCtx.Signal <- syscall.SIGTERM
	appStatus := <-appCtx.Status
	test.AssertEquals("", app.TerminatedStatus, appStatus.Status, t)
}
