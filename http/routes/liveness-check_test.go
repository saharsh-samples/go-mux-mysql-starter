package routes

import (
	"fmt"
	"net/http"
	"testing"

	myHttp "github.com/saharsh-samples/go-mux-sql-starter/http"
	myHttpTest "github.com/saharsh-samples/go-mux-sql-starter/http/test"
	"github.com/saharsh-samples/go-mux-sql-starter/test"
)

func TestGet(t *testing.T) {

	serverPort, appCtx := myHttpTest.StartTestServer([]myHttp.Routes{&LivenessCheck{}}, t)
	defer myHttpTest.StopTestServer(appCtx, t)
	url := fmt.Sprintf("http://localhost:%d/healthz", serverPort)

	resp, err := http.Get(url)
	test.AssertTrue("expected no errors", err == nil, t)
	defer resp.Body.Close()
	test.AssertEquals("", 200, resp.StatusCode, t)

}
