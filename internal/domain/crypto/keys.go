package crypto

import (
	"fmt"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

// KeyPair represents an ECDSA secp256k1 key pair.
// PrivateKey is 32 bytes, PublicKey is 33 bytes (compressed).
type KeyPair struct {
	PrivateKey []byte
	PublicKey  []byte
}

// GenerateKeyPair creates a new ECDSA secp256k1 key pair.
func GenerateKeyPair() (*KeyPair, error) {
	privateKey, err := secp256k1.GeneratePrivateKey()
	if err != nil {
		return nil, fmt.Errorf("crypto: failed to generate private key: %w", err)
	}

	privBytes := privateKey.Serialize()
	pubBytes := privateKey.PubKey().SerializeCompressed()

	return &KeyPair{
		PrivateKey: privBytes,
		PublicKey:  pubBytes,
	}, nil
}

// PublicKeyFromPrivateKey derives the compressed public key from a private key.
// The private key MUST be 32 bytes. Returns error for nil or wrong-length input.
func PublicKeyFromPrivateKey(privateKey []byte) ([]byte, error) {
	if privateKey == nil {
		return nil, fmt.Errorf("crypto: private key is nil")
	}
	if len(privateKey) != 32 {
		return nil, fmt.Errorf("crypto: private key must be 32 bytes, got %d", len(privateKey))
	}

	privKey := secp256k1.PrivKeyFromBytes(privateKey)
	return privKey.PubKey().SerializeCompressed(), nil
}
