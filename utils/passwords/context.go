package passwords

// ContextIn describes dependecies needed by this package
type ContextIn struct {
	Argon2Config Argon2Config
}

// ContextOut describes dependencies exported by this package
type ContextOut struct {
	PasswordHasher PasswordHasher
}

// Bootstrap initializes this module with ContextIn and exports
// resulting ContextOut
func Bootstrap(in *ContextIn) *ContextOut {

	out := &ContextOut{}
	out.PasswordHasher = &hasher{config: in.Argon2Config}

	return out
}
