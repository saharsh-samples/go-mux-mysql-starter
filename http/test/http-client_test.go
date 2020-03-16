package test

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	base "github.com/saharsh-samples/go-mux-sql-starter/http"
	httpUtils "github.com/saharsh-samples/go-mux-sql-starter/http/utils"
	"github.com/saharsh-samples/go-mux-sql-starter/test"
)

type pingPong struct {
	Value string
}

func (body *pingPong) Validate() error {
	if body.Value != "Ping" && body.Value != "Pong" {
		return errors.New("Value must be either 'Ping' or 'Pong'")
	}
	return nil
}

type pingPongRoute struct {
	jsonUtils httpUtils.JSONUtils
}

// Register endpoint+method handlers
func (resource *pingPongRoute) Register(agent base.RoutesAgent) {
	agent.RegisterGet("/ping", resource.Ping)
	agent.RegisterPost("/ping", resource.Pong)
	agent.RegisterPut("/ping", resource.Pong)
	agent.RegisterDelete("/ping", resource.Ping)
	agent.RegisterGet("/secure-ping", resource.SecurePing)
}

func (resource *pingPongRoute) Ping(w http.ResponseWriter, r *http.Request) {
	resource.jsonUtils.SetJSONResponse(w, http.StatusOK, &pingPong{Value: "Ping"})
}

func (resource *pingPongRoute) Pong(w http.ResponseWriter, r *http.Request) {

	ping := &pingPong{}
	err := resource.jsonUtils.ParseJSONRequest(r, ping, w)
	if err != nil {
		return
	}

	if ping.Value != "Ping" {
		resource.jsonUtils.BadRequest(w, "Did not get Ping")
		return
	}

	resource.jsonUtils.SetJSONResponse(w, http.StatusOK, &pingPong{Value: "Pong"})
}

func (resource *pingPongRoute) SecurePing(w http.ResponseWriter, r *http.Request) {

	authValues := r.Header["Authorization"]
	authErr := ""
	if len(authValues) != 1 {
		authErr = "Expected ONE and only ONE Authroization Header"
	} else if authValues[0] != "Bearer api-token" {
		authErr = "Unauthorized"
	}

	if authErr != "" {
		resource.jsonUtils.Forbidden(w, authErr)
	}

	resource.jsonUtils.SetJSONResponse(w, http.StatusOK, &pingPong{Value: "Ping"})
}

func TestGet(t *testing.T) {

	// Arrange
	serverPort, appCtx := StartTestServer(
		[]base.Routes{
			&pingPongRoute{
				jsonUtils: httpUtils.Bootstrap(nil).JSONUtils,
			},
		},
		t,
	)
	defer StopTestServer(appCtx, t)
	url := fmt.Sprintf("http://localhost:%d/ping", serverPort)

	// Act and Assert
	ping := &pingPong{}
	test.AssertEquals("", http.StatusOK, Get(t, url, ping, ""), t)
	test.AssertEquals("", "Ping", ping.Value, t)
}

func TestGetWithoutBody(t *testing.T) {

	// Arrange
	serverPort, appCtx := StartTestServer(
		[]base.Routes{
			&pingPongRoute{
				jsonUtils: httpUtils.Bootstrap(nil).JSONUtils,
			},
		},
		t,
	)
	defer StopTestServer(appCtx, t)
	url := fmt.Sprintf("http://localhost:%d/secure-ping", serverPort)

	// Act and Assert
	test.AssertEquals("", http.StatusOK, GetWithoutBody(t, url, "api-token"), t)
}

func TestPost(t *testing.T) {

	// Arrange
	serverPort, appCtx := StartTestServer(
		[]base.Routes{
			&pingPongRoute{
				jsonUtils: httpUtils.Bootstrap(nil).JSONUtils,
			},
		},
		t,
	)
	defer StopTestServer(appCtx, t)
	url := fmt.Sprintf("http://localhost:%d/ping", serverPort)

	// Act and Assert
	ping := &pingPong{Value: "Ping"}
	pong := &pingPong{}
	test.AssertEquals("", http.StatusOK, Post(t, url, ping, pong, ""), t)
	test.AssertEquals("", "Pong", pong.Value, t)
}

func TestPut(t *testing.T) {

	// Arrange
	serverPort, appCtx := StartTestServer(
		[]base.Routes{
			&pingPongRoute{
				jsonUtils: httpUtils.Bootstrap(nil).JSONUtils,
			},
		},
		t,
	)
	defer StopTestServer(appCtx, t)
	url := fmt.Sprintf("http://localhost:%d/ping", serverPort)

	// Act and Assert
	ping := &pingPong{Value: "Ping"}
	pong := &pingPong{}
	test.AssertEquals("", http.StatusOK, Put(t, url, ping, pong, ""), t)
	test.AssertEquals("", "Pong", pong.Value, t)
}

func TestDelete(t *testing.T) {

	// Arrange
	serverPort, appCtx := StartTestServer(
		[]base.Routes{
			&pingPongRoute{
				jsonUtils: httpUtils.Bootstrap(nil).JSONUtils,
			},
		},
		t,
	)
	defer StopTestServer(appCtx, t)
	url := fmt.Sprintf("http://localhost:%d/ping", serverPort)

	// Act and Assert
	ping := &pingPong{}
	test.AssertEquals("", http.StatusOK, Delete(t, url, ping, ""), t)
	test.AssertEquals("", "Ping", ping.Value, t)
}

func TestDeleteWithoutBody(t *testing.T) {

	// Arrange
	serverPort, appCtx := StartTestServer(
		[]base.Routes{
			&pingPongRoute{
				jsonUtils: httpUtils.Bootstrap(nil).JSONUtils,
			},
		},
		t,
	)
	defer StopTestServer(appCtx, t)
	url := fmt.Sprintf("http://localhost:%d/ping", serverPort)

	// Act and Assert
	test.AssertEquals("", http.StatusOK, DeleteWithoutBody(t, url, ""), t)
}

// Error testing

type dummyT struct {
	errorMsg string
}

func (d *dummyT) Error(args ...interface{}) {
	d.errorMsg = fmt.Sprint(args...)
	panic(d.errorMsg)
}

func (d *dummyT) Errorf(format string, args ...interface{}) {
	d.errorMsg = fmt.Sprintf(format, args...)
	panic(d.errorMsg)
}

func Test_doRequest_for_bad_request_body(t *testing.T) {
	d := &dummyT{}
	defer func() {
		test.AssertTrue("Expected error related to json marshalling", strings.Contains(d.errorMsg, "unsupported type: chan string"), t)
	}()
	defer test.AssertPanic("Expected panic", t)
	doRequest(d, "should-not-matter", http.MethodGet, make(chan string), nil, "")
}

func Test_doRequest_for_malformed_url(t *testing.T) {
	d := &dummyT{}
	defer func() {
		test.AssertTrue("Expected error related to malformed url", strings.Contains(d.errorMsg, "missing protocol scheme"), t)
	}()
	defer test.AssertPanic("Expected panic", t)
	doRequest(d, ":malformed-url", http.MethodGet, nil, nil, "")
}

func Test_doRequest_for_non_http_url(t *testing.T) {
	d := &dummyT{}
	defer func() {
		test.AssertTrue("Expected error related to non http url", strings.Contains(d.errorMsg, "unsupported protocol scheme"), t)
	}()
	defer test.AssertPanic("Expected panic", t)
	doRequest(d, "non-htp-url", http.MethodGet, nil, nil, "")
}

func Test_doRequest_for_bad_json_unmarshalling(t *testing.T) {

	// Arrange
	serverPort, appCtx := StartTestServer(
		[]base.Routes{
			&pingPongRoute{
				jsonUtils: httpUtils.Bootstrap(nil).JSONUtils,
			},
		},
		t,
	)
	defer StopTestServer(appCtx, t)
	url := fmt.Sprintf("http://localhost:%d/ping", serverPort)

	// Act and Assert
	d := &dummyT{}
	defer func() {
		test.AssertTrue("Expected error related to non http url", strings.Contains(d.errorMsg, "non-pointer string"), t)
	}()
	defer test.AssertPanic("Expected panic", t)
	Get(d, url, "", "")
}
