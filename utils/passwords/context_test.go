package passwords

import (
	"testing"

	"github.com/saharsh-samples/go-mux-sql-starter/test"
)

func TestBootstrap(t *testing.T) {

	out := Bootstrap(&ContextIn{
		Argon2Config: Argon2Config{},
	})

	// verify hasher
	hasher := out.PasswordHasher
	test.AssertFalse("", hasher == nil, t)
}
