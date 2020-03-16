package test

import (
	"fmt"
	"net/http"
	"testing"

	base "github.com/saharsh-samples/go-mux-sql-starter/http"
	"github.com/saharsh-samples/go-mux-sql-starter/test"
)

type testRoute struct{}

// Register endpoint+method handlers
func (resource *testRoute) Register(agent base.RoutesAgent) {
	agent.RegisterGet("/test", resource.Test)
}

// Create returns a 200
func (resource *testRoute) Test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestStartAndStopTestServer(t *testing.T) {

	serverPort, appCtx := StartTestServer([]base.Routes{&testRoute{}}, t)
	defer StopTestServer(appCtx, t)
	url := fmt.Sprintf("http://localhost:%d/test", serverPort)

	resp, err := http.Get(url)
	test.AssertTrue("expected no errors", err == nil, t)
	defer resp.Body.Close()
	test.AssertEquals("", 200, resp.StatusCode, t)
}
