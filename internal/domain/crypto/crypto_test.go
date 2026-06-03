package crypto

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"testing"
)

// Known test vectors from crypto/hash.md spec.
var knownVectors = []struct {
	name     string
	input    []byte
	expected string // SHA-256 hex
}{
	{
		name:     "hello",
		input:    []byte("hello"),
		expected: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
	},
	{
		name:     "empty",
		input:    []byte{},
		expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	},
	{
		name:     "nil",
		input:    nil,
		expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	},
	{
		name:     "parva-blockchain",
		input:    []byte("parva-blockchain"),
		expected: "", // computed at test time via stdlib
	},
}

func TestHashSHA256(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{name: "known vector 'hello'", input: []byte("hello")},
		{name: "empty slice", input: []byte{}},
		{name: "nil input", input: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HashSHA256(tt.input)

			// Must always return 32 bytes
			if len(result) != 32 {
				t.Fatalf("HashSHA256() returned %d bytes, want 32", len(result))
			}

			// Verify against stdlib
			h := sha256.Sum256(tt.input)
			if !bytes.Equal(result, h[:]) {
				t.Fatalf("HashSHA256() = %x, want %x", result, h[:])
			}
		})
	}
}

func TestHashSHA256KnownVector(t *testing.T) {
	// Known vector: SHA-256("hello")
	input := []byte("hello")
	expected := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"

	result := HashSHA256(input)
	resultHex := hex.EncodeToString(result)

	if resultHex != expected {
		t.Fatalf("HashSHA256(%q) = %s, want %s", string(input), resultHex, expected)
	}

	if len(result) != 32 {
		t.Fatalf("HashSHA256() returned %d bytes, want 32", len(result))
	}
}

func TestHashSHA256Empty(t *testing.T) {
	// Empty slice
	result := HashSHA256([]byte{})
	if len(result) != 32 {
		t.Fatalf("HashSHA256([]byte{}) returned %d bytes, want 32", len(result))
	}

	expected := sha256.Sum256([]byte{})
	if !bytes.Equal(result, expected[:]) {
		t.Fatalf("HashSHA256([]byte{}) = %x, want %x", result, expected[:])
	}
}

func TestHashSHA256Nil(t *testing.T) {
	// Nil input — must NOT panic and must equal SHA-256("")
	result := HashSHA256(nil)
	if len(result) != 32 {
		t.Fatalf("HashSHA256(nil) returned %d bytes, want 32", len(result))
	}

	expected := sha256.Sum256([]byte{})
	if !bytes.Equal(result, expected[:]) {
		t.Fatalf("HashSHA256(nil) = %x, want %x (SHA-256 of empty)", result, expected[:])
	}
}

func TestHashSHA256Consistency(t *testing.T) {
	input := []byte("parva-blockchain")

	// Call 3 times, all must be identical
	r1 := HashSHA256(input)
	r2 := HashSHA256(input)
	r3 := HashSHA256(input)

	if !bytes.Equal(r1, r2) {
		t.Fatal("HashSHA256() is not deterministic: call 1 and 2 differ")
	}
	if !bytes.Equal(r2, r3) {
		t.Fatal("HashSHA256() is not deterministic: call 2 and 3 differ")
	}
}

func TestHashSHA256LargeInput(t *testing.T) {
	// 1 MB input — must return 32 bytes without panic or timeout
	size := 1024 * 1024 // 1 MB
	data := make([]byte, size)
	_, err := rand.Read(data)
	if err != nil {
		t.Fatalf("failed to generate random data: %v", err)
	}

	result := HashSHA256(data)
	if len(result) != 32 {
		t.Fatalf("HashSHA256(1MB) returned %d bytes, want 32", len(result))
	}

	// Verify against stdlib
	expected := sha256.Sum256(data)
	if !bytes.Equal(result, expected[:]) {
		t.Fatal("HashSHA256(1MB) does not match stdlib")
	}
}

func TestHashDoubleSHA256(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{name: "known vector 'hello'", input: []byte("hello")},
		{name: "empty slice", input: []byte{}},
		{name: "nil input", input: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HashDoubleSHA256(tt.input)

			// Must always return 32 bytes
			if len(result) != 32 {
				t.Fatalf("HashDoubleSHA256() returned %d bytes, want 32", len(result))
			}

			// Verify composition: DoubleSHA256(X) == SHA256(SHA256(X))
			first := sha256.Sum256(tt.input)
			second := sha256.Sum256(first[:])
			if !bytes.Equal(result, second[:]) {
				t.Fatalf("HashDoubleSHA256() = %x, want %x (SHA256∘SHA256)", result, second[:])
			}
		})
	}
}

func TestHashDoubleSHA256Consistency(t *testing.T) {
	input := []byte("parva-blockchain")

	r1 := HashDoubleSHA256(input)
	r2 := HashDoubleSHA256(input)
	r3 := HashDoubleSHA256(input)

	if !bytes.Equal(r1, r2) {
		t.Fatal("HashDoubleSHA256() is not deterministic: call 1 and 2 differ")
	}
	if !bytes.Equal(r2, r3) {
		t.Fatal("HashDoubleSHA256() is not deterministic: call 2 and 3 differ")
	}
}

func TestHashDoubleSHA256Composition(t *testing.T) {
	// DoubleSHA256(X) MUST equal SHA256(SHA256(X)) for any input
	inputs := [][]byte{
		[]byte("hello"),
		[]byte("parva-blockchain"),
		[]byte{},
		nil,
		[]byte("x"),
		make([]byte, 1024), // 1 KB zeros
	}

	for _, input := range inputs {
		doubleResult := HashDoubleSHA256(input)
		singleResult := HashSHA256(HashSHA256(input))

		if !bytes.Equal(doubleResult, singleResult) {
			t.Fatalf("HashDoubleSHA256(%v) ≠ HashSHA256(HashSHA256(%v)):\n  double = %x\n  single = %x",
				input, input, doubleResult, singleResult)
		}
	}
}
