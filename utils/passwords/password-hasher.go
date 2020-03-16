package passwords

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// ---
// Password Hashing
//
// Shamelessly copied from
// https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go
// ---

// PasswordHasher generates hashes for cleartext passwords
type PasswordHasher interface {
	GeneratePasswordHash(password string) (encodedHash string, err error)
}

// DefaultArgon2Memory value
const DefaultArgon2Memory = 64 * 1024

// DefaultArgon2Iterations value
const DefaultArgon2Iterations = 3

// DefaultArgon2Parallelism value
const DefaultArgon2Parallelism = 2

// DefaultArgon2SaltLength value
const DefaultArgon2SaltLength = 16

// DefaultArgon2KeyLength value
const DefaultArgon2KeyLength = 32

// Argon2Config values for configuring the Argon2 hashing algorithm
type Argon2Config struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

type hasher struct {
	config Argon2Config
}

// GeneratePasswordHash generates an Argon2 hash for specified string
func (hasher *hasher) GeneratePasswordHash(password string) (encodedHash string, err error) {

	p := hasher.config

	// Generate a cryptographically secure random salt.
	salt, err := generateRandomBytes(p.SaltLength)
	if err != nil {
		return "", err
	}

	// Pass the plaintext password, salt and parameters to the argon2.IDKey
	// function. This will generate a hash of the password using the Argon2id
	// variant.
	hash := argon2.IDKey(
		[]byte(password),
		salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength,
	)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation.
	encodedHash = fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, p.Memory, p.Iterations, p.Parallelism, b64Salt, b64Hash,
	)

	return encodedHash, nil
}

// ComparePasswordAndHash verifies a password is same as password used to generate
// the given hash
func ComparePasswordAndHash(password, encodedHash string) (match bool, err error) {

	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func decodeHash(encodedHash string) (p *Argon2Config, salt, hash []byte, err error) {

	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, errors.New("the encoded hash is not in the correct format")
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, errors.New("incompatible version of argon2")
	}

	p = &Argon2Config{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.SaltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.KeyLength = uint32(len(hash))

	return p, salt, hash, nil

}
