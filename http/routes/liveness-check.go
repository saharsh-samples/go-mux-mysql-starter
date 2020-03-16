package routes

import (
	"net/http"

	base "github.com/saharsh-samples/go-mux-sql-starter/http"
)

// LivenessCheck can be used as a simple liveness health check
type LivenessCheck struct{}

// Register endpoint+method handlers
func (resource *LivenessCheck) Register(agent base.RoutesAgent) {
	agent.RegisterGet("/", resource.Get)
	agent.RegisterGet("/healthz", resource.Get)
}

// Get returns a 200
func (resource *LivenessCheck) Get(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
