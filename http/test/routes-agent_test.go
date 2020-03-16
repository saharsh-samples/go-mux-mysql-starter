package test

import (
	"net/http"
	"testing"

	base "github.com/saharsh-samples/go-mux-sql-starter/http"
)

type route struct{}

// Register endpoint+method handlers
func (resource *route) Register(agent base.RoutesAgent) {
	agent.RegisterGet("/test", resource.dummyGet)
	agent.RegisterPost("/test", resource.dummyPost)
	agent.RegisterPut("/test", resource.dummyPut)
	agent.RegisterDelete("/test", resource.dummyDelete)
}

func (resource *route) dummyGet(w http.ResponseWriter, r *http.Request)    {}
func (resource *route) dummyPost(w http.ResponseWriter, r *http.Request)   {}
func (resource *route) dummyPut(w http.ResponseWriter, r *http.Request)    {}
func (resource *route) dummyDelete(w http.ResponseWriter, r *http.Request) {}

func TestMockRoutesAgent(t *testing.T) {

	// Arrange
	agent := NewMockRoutesAgent()
	route := route{}

	// Act
	route.Register(agent)

	// Assert
	verifyThatTestRoute := agent.VerifyThatRoute(t, "/test")

	// GET
	verifyThatTestRoute.ForHTTPMethod(http.MethodGet).UsesHandler(route.dummyGet)
	verifyThatTestRoute.ForHTTPMethod(http.MethodPost).UsesHandler(route.dummyPost)
	verifyThatTestRoute.ForHTTPMethod(http.MethodPut).UsesHandler(route.dummyPut)
	verifyThatTestRoute.ForHTTPMethod(http.MethodDelete).UsesHandler(route.dummyDelete)
}

func TestMockRoutesAgent_with_handlerFuncRegistrationOverride(t *testing.T) {

	// Arrange
	agent := NewMockRoutesAgent()
	route := route{}

	i := 0
	agent.OverrideHandlerFuncRegistration(func(method string, url string, registered func(http.ResponseWriter, *http.Request)) interface{} {
		if i < 3 {
			i++
			return route.dummyGet
		}
		return nil
	})

	// Act
	route.Register(agent)

	// Assert
	verifyThatTestRoute := agent.VerifyThatRoute(t, "/test")

	verifyThatTestRoute.ForHTTPMethod(http.MethodGet).UsesHandler(route.dummyGet)
	verifyThatTestRoute.ForHTTPMethod(http.MethodPost).UsesHandler(route.dummyGet)
	verifyThatTestRoute.ForHTTPMethod(http.MethodPut).UsesHandler(route.dummyGet)

	verifyThatTestRoute.ForHTTPMethod(http.MethodDelete).UsesHandler(route.dummyDelete)
}
