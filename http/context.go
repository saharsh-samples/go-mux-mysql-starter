package http

// ContextIn describes dependecies needed by this package
type ContextIn struct {
	Port             int
	RoutesToRegister []Routes
}

// ContextOut describes dependencies exported by this package
type ContextOut struct {
	Server Server
}

// Bootstrap initializes this module with ContextIn and exports
// resulting ContextOut
func Bootstrap(in *ContextIn) *ContextOut {

	out := &ContextOut{}
	out.Server = &server{
		port:   in.Port,
		routes: in.RoutesToRegister,
	}

	return out
}
