package crypto

import (
	"fmt"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	secp_ecdsa "github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
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

// Sign signs the message with the given private key.
// Private key MUST be 32 bytes. Returns DER-encoded signature.
func Sign(privateKey []byte, msg []byte) ([]byte, error) {
	if privateKey == nil {
		return nil, fmt.Errorf("crypto: private key is nil")
	}
	if len(privateKey) != 32 {
		return nil, fmt.Errorf("crypto: private key must be 32 bytes, got %d", len(privateKey))
	}
	if msg == nil {
		return nil, fmt.Errorf("crypto: message is nil")
	}

	priv := secp256k1.PrivKeyFromBytes(privateKey)
	signature := secp_ecdsa.Sign(priv, msg)
	return signature.Serialize(), nil
}

// Verify checks that the signature is valid for the given message and public key.
// Public key MUST be 33 bytes (compressed). Returns false for invalid inputs (no panic).
func Verify(publicKey []byte, msg []byte, sig []byte) bool {
	if publicKey == nil || msg == nil || sig == nil {
		return false
	}

	pubKey, err := secp256k1.ParsePubKey(publicKey)
	if err != nil {
		return false
	}

	signature, err := secp_ecdsa.ParseDERSignature(sig)
	if err != nil {
		return false
	}

	return signature.Verify(msg, pubKey)
}
