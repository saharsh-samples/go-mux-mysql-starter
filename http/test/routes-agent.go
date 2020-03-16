package test

import (
	"net/http"
	"reflect"
	"runtime"

	base "github.com/saharsh-samples/go-mux-sql-starter/http"
	"github.com/saharsh-samples/go-mux-sql-starter/test"
)

// ---
// Exposed types/functions
// ---

// NewMockRoutesAgent with test friendly features
func NewMockRoutesAgent() MockRoutesAgent {
	return &mockRoutesAgent{
		httpHandlers: make(map[string]string),
	}
}

// HandlerFuncRegistrationOverride function type
type HandlerFuncRegistrationOverride func(method string, url string, registered func(http.ResponseWriter, *http.Request)) interface{}

// MockRoutesAgent mocks base.RoutesAgent and adds test friendly features
type MockRoutesAgent interface {
	base.RoutesAgent
	OverrideHandlerFuncRegistration(HandlerFuncRegistrationOverride)
	VerifyThatRoute(t test.T, url string) RouteVerifierFactory
}

// RouteVerifierFactory creates RouteVerifier instances
type RouteVerifierFactory interface {
	ForHTTPMethod(string) RouteVerifier
}

// RouteVerifier is used to verify proper configuration of HTTP routes
type RouteVerifier interface {
	UsesHandler(interface{}) RouteVerifier
}

// StringifyHandlerFunc for comparisons in testing
func StringifyHandlerFunc(handlerFunc interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(handlerFunc).Pointer()).Name()
}

// ---
// MockRoutesAgent impl
// ---

type mockRoutesAgent struct {
	httpHandlers                    map[string]string
	overrideHandlerFuncRegistration HandlerFuncRegistrationOverride
}

func (agent *mockRoutesAgent) register(method string, path string, f func(w http.ResponseWriter, r *http.Request)) {
	var handlerFunc interface{} = f
	if agent.overrideHandlerFuncRegistration != nil {
		override := agent.overrideHandlerFuncRegistration(method, path, f)
		if override != nil {
			handlerFunc = override
		}
	}
	agent.httpHandlers[method+":"+path] = StringifyHandlerFunc(handlerFunc)
}

func (agent *mockRoutesAgent) RegisterGet(path string, f func(w http.ResponseWriter, r *http.Request)) {
	agent.register(http.MethodGet, path, f)
}

func (agent *mockRoutesAgent) RegisterPost(path string, f func(w http.ResponseWriter, r *http.Request)) {
	agent.register(http.MethodPost, path, f)
}

func (agent *mockRoutesAgent) RegisterPut(path string, f func(w http.ResponseWriter, r *http.Request)) {
	agent.register(http.MethodPut, path, f)
}

func (agent *mockRoutesAgent) RegisterDelete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	agent.register(http.MethodDelete, path, f)
}

func (agent *mockRoutesAgent) OverrideHandlerFuncRegistration(override HandlerFuncRegistrationOverride) {
	agent.overrideHandlerFuncRegistration = override
}

func (agent *mockRoutesAgent) VerifyThatRoute(t test.T, url string) RouteVerifierFactory {
	return &routeVerifierFactory{t: t, agent: agent, url: url}
}

func (agent *mockRoutesAgent) getHandler(method string, path string) (string, bool) {
	handler, found := agent.httpHandlers[method+":"+path]
	return handler, found
}

// ---
// RouteVerifier impl
// ---

type routeVerifierFactory struct {
	t     test.T
	agent *mockRoutesAgent
	url   string
}

func (factory *routeVerifierFactory) ForHTTPMethod(method string) RouteVerifier {
	return &routeVerifier{t: factory.t, agent: factory.agent, url: factory.url, method: method}
}

type routeVerifier struct {
	t      test.T
	agent  *mockRoutesAgent
	url    string
	method string
}

func (v *routeVerifier) UsesHandler(handler interface{}) RouteVerifier {
	expectedHandler := StringifyHandlerFunc(handler)
	actualHandler, _ := v.agent.getHandler(v.method, v.url)
	test.AssertEquals("Expected "+expectedHandler+" to be handler function for "+v.method+" "+v.url, expectedHandler, actualHandler, v.t)
	return v
}
