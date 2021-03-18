package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Middlewares runs before every route
type Middlewares []func(http.Handler) http.Handler

// TLSConfiguration for server
type TLSConfiguration struct {
	// CertFile to use for serving TLS traffic
	CertFile string
	// KeyFile to use for serving TLS traffic
	KeyFile string
}

// ContextIn describes dependecies needed by this package
type ContextIn struct {
	Port                  int
	RoutesToRegister      []Routes
	MiddlewaresToRegister Middlewares
	TLSConfiguration      *TLSConfiguration
}

// ContextOut describes dependencies exported by this package
type ContextOut struct {
	Server Server
}

// Bootstrap initializes this module with ContextIn and exports
// resulting ContextOut
func Bootstrap(in *ContextIn) *ContextOut {

	// transform middleware slice
	middlewares := make([]mux.MiddlewareFunc, len(in.MiddlewaresToRegister))
	for i, middleware := range in.MiddlewaresToRegister {
		middlewares[i] = middleware
	}

	out := &ContextOut{}
	out.Server = &server{
		port:        in.Port,
		routes:      in.RoutesToRegister,
		middlewares: middlewares,
		tlsConfig:   in.TLSConfiguration,
	}

	return out
}
