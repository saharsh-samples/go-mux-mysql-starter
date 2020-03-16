package passwords

import (
	"testing"

	"github.com/saharsh-samples/go-mux-sql-starter/test"
)

func TestPasswordHashing(t *testing.T) {

	// create hasher
	var hasher PasswordHasher = &hasher{
		config: Argon2Config{
			Memory:      DefaultArgon2Memory,
			Iterations:  DefaultArgon2Iterations,
			Parallelism: DefaultArgon2Parallelism,
			SaltLength:  DefaultArgon2SaltLength,
			KeyLength:   DefaultArgon2KeyLength,
		},
	}

	// Generate hash
	hash, err := hasher.GeneratePasswordHash("P@ssw0rd")
	test.AssertFalse("Expected hash to be nonempty", hash == "", t)
	test.AssertTrue("Expected hashing to succeed", err == nil, t)
	// fmt.Printf("Hashed Password: '%v'\n", hash)

	// Compare valid password against hash
	match, err2 := ComparePasswordAndHash("P@ssw0rd", hash)
	test.AssertTrue("Expected hash to match password", match, t)
	test.AssertTrue("Expected matching to succeed", err2 == nil, t)

	// Compare invalid password against hash
	match, err2 = ComparePasswordAndHash("Passw0rd", hash)
	test.AssertFalse("Expected hash to NOT match password", match, t)
	test.AssertTrue("Expected matching to succeed", err2 == nil, t)
}

func TestComparePasswordAndHash_where_hash_has_incorrect_number_of_fields(t *testing.T) {

	_, err := ComparePasswordAndHash(
		"P@ssw0rd",
		"$v=19"+
			"$m=65536,t=3,p=2"+
			"$gGid0u4NBFEUYt5jxH8a6g"+
			"$rJkjgQFtGZO6srvudP7Ayl6aLXAbpW2GJcN7R7GVHZQ",
	)

	test.AssertFalse("Expected matching to error out", err == nil, t)
	test.AssertEquals("", "the encoded hash is not in the correct format", err.Error(), t)
}

func TestComparePasswordAndHash_where_hash_has_corrupt_argon2_version(t *testing.T) {

	_, err := ComparePasswordAndHash(
		"P@ssw0rd",
		"$argon2id"+
			"$v=abc"+
			"$m=65536,t=3,p=2"+
			"$gGid0u4NBFEUYt5jxH8a6g"+
			"$rJkjgQFtGZO6srvudP7Ayl6aLXAbpW2GJcN7R7GVHZQ",
	)

	test.AssertFalse("Expected matching to error out", err == nil, t)
	test.AssertEquals("", "expected integer", err.Error(), t)
}

func TestComparePasswordAndHash_where_hash_has_unknown_argon2_version(t *testing.T) {

	_, err := ComparePasswordAndHash(
		"P@ssw0rd",
		"$argon2id"+
			"$v=18"+
			"$m=65536,t=3,p=2"+
			"$gGid0u4NBFEUYt5jxH8a6g"+
			"$rJkjgQFtGZO6srvudP7Ayl6aLXAbpW2GJcN7R7GVHZQ",
	)

	test.AssertFalse("Expected matching to error out", err == nil, t)
	test.AssertEquals("", "incompatible version of argon2", err.Error(), t)
}

func TestComparePasswordAndHash_where_hash_has_corrupt_argon2_memory(t *testing.T) {

	_, err := ComparePasswordAndHash(
		"P@ssw0rd",
		"$argon2id"+
			"$v=19"+
			"$m=6553i,t=3,p=2"+
			"$gGid0u4NBFEUYt5jxH8a6g"+
			"$rJkjgQFtGZO6srvudP7Ayl6aLXAbpW2GJcN7R7GVHZQ",
	)

	test.AssertFalse("Expected matching to error out", err == nil, t)
	test.AssertEquals("", "input does not match format", err.Error(), t)
}
