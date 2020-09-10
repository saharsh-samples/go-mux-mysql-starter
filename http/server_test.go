package http

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/saharsh-samples/go-mux-sql-starter/test"
	"github.com/saharsh-samples/go-mux-sql-starter/utils"
)

func TestPreInitState(t *testing.T) {

	// arrange
	server := Bootstrap(&ContextIn{Port: 0}).Server

	// act
	isReady := server.IsReady()
	_, portError := server.Port()
	shutdownError := server.Shutdown()

	// assert
	test.AssertFalse("Expected IsReady to return false", isReady, t)
	test.AssertEquals("", "Server is not running", portError.Error(), t)
	test.AssertEquals("", "Server is not running", shutdownError.Error(), t)

}

func TestRun_with_bad_port(t *testing.T) {

	// arrange
	server := Bootstrap(&ContextIn{Port: -1}).Server

	// assert via defer
	defer test.AssertPanic("Expected panic during Run", t)

	// act
	server.Run()

}

func TestRun_with_no_routes(t *testing.T) {

	// arrange
	server := Bootstrap(&ContextIn{Port: 0}).Server

	// act
	go server.Run()
	startupError := utils.WaitTill(func() bool { return server.IsReady() }, 10)

	// assert
	test.AssertTrue("Expected no errors starting up", startupError == nil, t)
	test.AssertTrue("Expected IsReady to return false", server.IsReady(), t)

	port, portErr := server.Port()
	test.AssertTrue("Expected no errors getting port", portErr == nil, t)
	test.AssertTrue("Expected server to pick a random port", port > 0, t)

	test.AssertTrue("Expected no errors shutting down", server.Shutdown() == nil, t)

}

func TestRun_with_some_routes_and_middlewares(t *testing.T) {

	// arrange
	requestCount := 0
	middlewares := make(Middlewares, 1)
	middlewares[0] = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestCount++
			h.ServeHTTP(w, r)
		})
	}

	server := Bootstrap(&ContextIn{
		Port:                  0,
		RoutesToRegister:      []Routes{&route1{}, &route2{}},
		MiddlewaresToRegister: middlewares,
	}).Server

	// act
	go server.Run()
	startupError := utils.WaitTill(func() bool { return server.IsReady() }, 10)

	// assert

	// ---
	// Assert healthy startup
	// ---

	test.AssertTrue("Expected no errors starting up", startupError == nil, t)
	test.AssertTrue("Expected IsReady to return false", server.IsReady(), t)

	port, portErr := server.Port()
	test.AssertTrue("Expected no errors getting port", portErr == nil, t)
	test.AssertTrue("Expected server to pick a random port", port > 0, t)

	// ---
	// Test routes
	// ---

	url := fmt.Sprintf("http://localhost:%d/route1", port)

	// Get
	resp, respErr := http.Get(url)
	test.AssertTrue("Expected no errors doing GET /route1", respErr == nil, t)
	test.AssertEquals("", 200, resp.StatusCode, t)

	// Post (expect 405 since post is defined for route2, not route1, below)
	resp, respErr = http.Post(url, "text/plain", nil)
	test.AssertTrue("Expected no errors doing POST /route1", respErr == nil, t)
	test.AssertEquals("", 405, resp.StatusCode, t)

	// Post
	url = fmt.Sprintf("http://localhost:%d/route2", port)
	resp, respErr = http.Post(url, "text/plain", nil)
	test.AssertTrue("Expected no errors doing POST /route1", respErr == nil, t)
	test.AssertEquals("", 200, resp.StatusCode, t)

	// Middleware
	test.AssertEquals("", 2, requestCount, t)

	// ---
	// Test shutdown
	// ---

	test.AssertTrue("Expected no errors shutting down", server.Shutdown() == nil, t)

}

type route1 struct{}

func (resource *route1) Register(agent RoutesAgent) {
	agent.RegisterGet("/route1", SuccessHandler)
	agent.RegisterPut("/route1", SuccessHandler)
}

type route2 struct{}

func (resource *route2) Register(agent RoutesAgent) {
	agent.RegisterPost("/route2", SuccessHandler)
	agent.RegisterDelete("/route2", SuccessHandler)
}

func SuccessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
