package crypto

import (
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() returned error: %v", err)
	}
	if kp == nil {
		t.Fatal("GenerateKeyPair() returned nil KeyPair")
	}
	if len(kp.PrivateKey) != 32 {
		t.Fatalf("PrivateKey length = %d, want 32", len(kp.PrivateKey))
	}
	if len(kp.PublicKey) != 33 {
		t.Fatalf("PublicKey length = %d, want 33", len(kp.PublicKey))
	}
}

func TestGenerateKeyPairUnique(t *testing.T) {
	kp1, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() 1 failed: %v", err)
	}
	kp2, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() 2 failed: %v", err)
	}

	// Two consecutive calls must produce different keys
	if string(kp1.PrivateKey) == string(kp2.PrivateKey) {
		t.Fatal("Two GenerateKeyPair calls produced identical private keys")
	}
	if string(kp1.PublicKey) == string(kp2.PublicKey) {
		t.Fatal("Two GenerateKeyPair calls produced identical public keys")
	}
}

func TestPublicKeyFromPrivateKey(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() failed: %v", err)
	}

	pub, err := PublicKeyFromPrivateKey(kp.PrivateKey)
	if err != nil {
		t.Fatalf("PublicKeyFromPrivateKey() returned error: %v", err)
	}
	if len(pub) != 33 {
		t.Fatalf("PublicKeyFromPrivateKey() returned %d bytes, want 33", len(pub))
	}
	if string(pub) != string(kp.PublicKey) {
		t.Fatal("PublicKeyFromPrivateKey() returned different public key than KeyPair.PublicKey")
	}
}

func TestPublicKeyFromPrivateKeyDeterministic(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() failed: %v", err)
	}

	// Call 3 times with same private key — all must return same public key
	pub1, err := PublicKeyFromPrivateKey(kp.PrivateKey)
	if err != nil {
		t.Fatalf("PublicKeyFromPrivateKey() call 1 failed: %v", err)
	}
	pub2, err := PublicKeyFromPrivateKey(kp.PrivateKey)
	if err != nil {
		t.Fatalf("PublicKeyFromPrivateKey() call 2 failed: %v", err)
	}
	pub3, err := PublicKeyFromPrivateKey(kp.PrivateKey)
	if err != nil {
		t.Fatalf("PublicKeyFromPrivateKey() call 3 failed: %v", err)
	}

	if string(pub1) != string(pub2) {
		t.Fatal("PublicKeyFromPrivateKey() is not deterministic: call 1 and 2 differ")
	}
	if string(pub2) != string(pub3) {
		t.Fatal("PublicKeyFromPrivateKey() is not deterministic: call 2 and 3 differ")
	}
}

func TestPublicKeyFromPrivateKeyInvalid(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{name: "nil private key", input: nil},
		{name: "empty slice", input: []byte{}},
		{name: "too short (8 bytes)", input: []byte("deadbeef")},
		{name: "too long (64 bytes)", input: make([]byte, 64)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pub, err := PublicKeyFromPrivateKey(tt.input)
			if err == nil {
				t.Fatal("PublicKeyFromPrivateKey() expected error, got nil")
			}
			if pub != nil {
				t.Fatal("PublicKeyFromPrivateKey() expected nil public key on error")
			}
		})
	}
}

func TestKeyPairLengths(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() failed: %v", err)
	}

	if len(kp.PrivateKey) != 32 {
		t.Errorf("PrivateKey length = %d, want 32", len(kp.PrivateKey))
	}
	if len(kp.PublicKey) != 33 {
		t.Errorf("PublicKey length = %d, want 33", len(kp.PublicKey))
	}
}
