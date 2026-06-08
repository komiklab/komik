package utils

import (
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

const (
	saltBytes  = 16     // 128-bit salt
	iterations = 210000 // OWASP recommendation for SHA-512
	keyBytes   = 64     // 512-bit output key size
)

func HashPassword(password string) (string, error) {
	// 1. Generate a cryptographically secure random salt
	salt := make([]byte, saltBytes)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// 2. Derive the key using PBKDF2-HMAC-SHA512
	hash := pbkdf2.Key([]byte(password), salt, iterations, keyBytes, sha512.New)

	// 3. Encode to base64 for easy storage in a database column
	b64Salt := base64.StdEncoding.EncodeToString(salt)
	b64Hash := base64.StdEncoding.EncodeToString(hash)

	// Format matching standard formats: $pbkdf2-sha512$iterations$salt$hash
	return fmt.Sprintf("$pbkdf2-sha512$%d$%s$%s", iterations, b64Salt, b64Hash), nil
}

// VerifyPassword securely validates an incoming plaintext password against a stored hash string
func VerifyPassword(password, storedHash string) (bool, error) {
	// 1. Parse the stored hash string back into components
	var parsedIterations int
	var b64Salt, b64Hash string

	_, err := fmt.Sscanf(storedHash, "$pbkdf2-sha512$%d$%s$%s", &parsedIterations, &b64Salt, &b64Hash)
	if err != nil {
		return false, fmt.Errorf("invalid hash format: %w", err)
	}

	// 2. Decode the salt and the target hash from Base64
	salt, err := base64.StdEncoding.DecodeString(b64Salt)
	if err != nil {
		return false, err
	}
	targetHash, err := base64.StdEncoding.DecodeString(b64Hash)
	if err != nil {
		return false, err
	}

	// 3. Compute the PBKDF2 hash of the incoming password using the parsed parameters
	computedHash := pbkdf2.Key([]byte(password), salt, parsedIterations, len(targetHash), sha512.New)

	// 4. Use constant-time comparison to prevent side-channel timing attacks
	if subtle.ConstantTimeCompare(targetHash, computedHash) == 1 {
		return true, nil
	}
	return false, nil
}
