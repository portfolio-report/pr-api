package argon2

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Inspired by https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go

// HashPasswordDefault returns hash using argon2id and default parameters
func HashPasswordDefault(password string) (string, error) {
	return HashPassword(password, 32, 16, 1, 65536, 2)
}

// HashPassword returns hash using argon2id and the provided parameters
func HashPassword(password string, keyLength, saltLength, iterations, memory uint32, parallelism uint8) (string, error) {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("could not create salt: %w", err)
	}

	var key []byte
	key = argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLength)

	hash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		memory,
		iterations,
		parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key))

	return hash, nil
}

// VerifyPassword checks if password matches the hash, supports argon2id and argon2i
func VerifyPassword(password, hash string) (match bool, err error) {
	hashParts := strings.Split(hash, "$")
	if len(hashParts) != 6 {
		return false, errors.New("invalid hash format")
	}

	variant := hashParts[1]

	var version int
	_, err = fmt.Sscanf(hashParts[2], "v=%d", &version)
	if err != nil {
		return false, fmt.Errorf("could not parse version: %w", err)
	}
	if version != argon2.Version {
		return false, fmt.Errorf("incompatible version: %d", version)
	}

	var memory uint32
	var iterations uint32
	var parallelism uint8

	_, err = fmt.Sscanf(hashParts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return false, fmt.Errorf("could not parse parameters: %w", err)
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(hashParts[4])
	if err != nil {
		return false, fmt.Errorf("could not decode salt: %w", err)
	}

	key, err := base64.RawStdEncoding.Strict().DecodeString(hashParts[5])
	if err != nil {
		return false, fmt.Errorf("could not decode key: %w", err)
	}

	var derivedKey []byte
	switch variant {
	case "argon2id":
		derivedKey = argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(key)))
	case "argon2i":
		derivedKey = argon2.Key([]byte(password), salt, iterations, memory, parallelism, uint32(len(key)))
	default:
		return false, fmt.Errorf("unknown variant: %s", variant)
	}

	if subtle.ConstantTimeEq(int32(len(key)), int32(len(derivedKey))) == 0 {
		return false, nil
	}
	return subtle.ConstantTimeCompare(key, derivedKey) == 1, nil
}
