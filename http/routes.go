package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

// ------
// Routes
// ------

// Routes is base type for types that implement routes
type Routes interface {
	Register(agent RoutesAgent)
}

// ------------
// Routes Agent
// ------------

// RoutesAgent is used to expose HTTP endpoints
type RoutesAgent interface {
	RegisterGet(path string, f func(w http.ResponseWriter, r *http.Request))
	RegisterPost(path string, f func(w http.ResponseWriter, r *http.Request))
	RegisterPut(path string, f func(w http.ResponseWriter, r *http.Request))
	RegisterDelete(path string, f func(w http.ResponseWriter, r *http.Request))
}

// --------
// Internal
// --------

type routesAgent struct {
	router *mux.Router
}

func (agent *routesAgent) RegisterGet(path string, f func(w http.ResponseWriter, r *http.Request)) {
	agent.router.HandleFunc(path, f).Methods("GET")
}

func (agent *routesAgent) RegisterPost(path string, f func(w http.ResponseWriter, r *http.Request)) {
	agent.router.HandleFunc(path, f).Methods("POST")
}

func (agent *routesAgent) RegisterPut(path string, f func(w http.ResponseWriter, r *http.Request)) {
	agent.router.HandleFunc(path, f).Methods("PUT")
}

func (agent *routesAgent) RegisterDelete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	agent.router.HandleFunc(path, f).Methods("DELETE")
}
