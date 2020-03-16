package utils

// ContextIn describes dependecies needed by this package
type ContextIn struct {
	// Nothing
}

// ContextOut describes dependencies exported by this package
type ContextOut struct {
	JSONUtils JSONUtils
	URLUtils  URLUtils
}

// Bootstrap initializes this module with ContextIn and exports
// resulting ContextOut
func Bootstrap(in *ContextIn) *ContextOut {

	out := &ContextOut{}
	out.JSONUtils = &jsonUtils{}
	out.URLUtils = &urlUtils{}

	return out
}
