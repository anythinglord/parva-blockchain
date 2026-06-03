package crypto

import "crypto/sha256"

// HashSHA256 computes the SHA-256 hash of data.
// Returns 32 bytes. Nil and empty input return SHA-256("").
func HashSHA256(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}

// HashDoubleSHA256 computes SHA-256(SHA-256(data)).
// Returns 32 bytes.
func HashDoubleSHA256(data []byte) []byte {
	first := HashSHA256(data)
	return HashSHA256(first)
}
