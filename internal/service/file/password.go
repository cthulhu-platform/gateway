package file

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argon2Time    = 1
	argon2Memory  = 64 * 1024 // 64 MB
	argon2Threads = 4
	argon2KeyLen  = 32
	saltLen       = 16
)

// HashPassword hashes a password using Argon2ID
// Returns hash in format: "argon2id$<base64-salt>$<base64-hash>"
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	// Generate random salt
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Hash password with Argon2ID
	hash := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)

	// Encode salt and hash as base64
	saltB64 := base64.RawStdEncoding.EncodeToString(salt)
	hashB64 := base64.RawStdEncoding.EncodeToString(hash)

	// Return format: "argon2id$<salt>$<hash>"
	return "argon2id$" + saltB64 + "$" + hashB64, nil
}

// VerifyPassword verifies a password against a stored hash
// Hash format: "argon2id$<base64-salt>$<base64-hash>"
func VerifyPassword(password, hash string) bool {
	if password == "" || hash == "" {
		return false
	}

	// Parse hash format: "argon2id$<salt>$<hash>"
	parts := strings.Split(hash, "$")
	if len(parts) != 3 || parts[0] != "argon2id" {
		return false
	}

	// Decode salt
	salt, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	// Decode stored hash
	storedHash, err := base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil {
		return false
	}

	// Compute hash with same parameters
	computedHash := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)

	// Constant-time comparison to prevent timing attacks
	return subtle.ConstantTimeCompare(storedHash, computedHash) == 1
}
