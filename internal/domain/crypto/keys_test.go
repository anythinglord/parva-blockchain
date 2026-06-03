package crypto

import (
	"sync"
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

func TestSignAndVerify(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() failed: %v", err)
	}

	msg := []byte("hello, parva blockchain")
	sig, err := Sign(kp.PrivateKey, msg)
	if err != nil {
		t.Fatalf("Sign() returned error: %v", err)
	}
	if len(sig) == 0 {
		t.Fatal("Sign() returned empty signature")
	}

	valid := Verify(kp.PublicKey, msg, sig)
	if !valid {
		t.Fatal("Verify() returned false for valid signature")
	}
}

func TestVerifyFailsWrongKey(t *testing.T) {
	kp1, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() 1 failed: %v", err)
	}
	kp2, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() 2 failed: %v", err)
	}

	msg := []byte("test message")
	sig, err := Sign(kp1.PrivateKey, msg)
	if err != nil {
		t.Fatalf("Sign() failed: %v", err)
	}

	// Verify with wrong key
	valid := Verify(kp2.PublicKey, msg, sig)
	if valid {
		t.Fatal("Verify() returned true with wrong public key")
	}
}

func TestVerifyFailsTamperedMessage(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() failed: %v", err)
	}

	original := []byte("original message")
	sig, err := Sign(kp.PrivateKey, original)
	if err != nil {
		t.Fatalf("Sign() failed: %v", err)
	}

	// Verify with tampered message
	tampered := []byte("tampered message")
	valid := Verify(kp.PublicKey, tampered, sig)
	if valid {
		t.Fatal("Verify() returned true with tampered message")
	}
}

func TestVerifyFailsTamperedSignature(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() failed: %v", err)
	}

	msg := []byte("test message")
	sig, err := Sign(kp.PrivateKey, msg)
	if err != nil {
		t.Fatalf("Sign() failed: %v", err)
	}

	// Flip a bit in the signature
	if len(sig) > 0 {
		tamperedSig := make([]byte, len(sig))
		copy(tamperedSig, sig)
		tamperedSig[0] ^= 0x01

		valid := Verify(kp.PublicKey, msg, tamperedSig)
		if valid {
			t.Fatal("Verify() returned true with tampered signature")
		}
	}
}

func TestVerifyNoPanicNilInput(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() failed: %v", err)
	}

	msg := []byte("test message")
	sig, err := Sign(kp.PrivateKey, msg)
	if err != nil {
		t.Fatalf("Sign() failed: %v", err)
	}

	tests := []struct {
		name string
		pub  []byte
		m    []byte
		s    []byte
	}{
		{name: "nil public key", pub: nil, m: msg, s: sig},
		{name: "nil message", pub: kp.PublicKey, m: nil, s: sig},
		{name: "nil signature", pub: kp.PublicKey, m: msg, s: nil},
		{name: "empty public key", pub: []byte{}, m: msg, s: sig},
		{name: "empty signature", pub: kp.PublicKey, m: msg, s: []byte{}},
		{name: "all nil", pub: nil, m: nil, s: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Must not panic
			result := Verify(tt.pub, tt.m, tt.s)
			if result {
				t.Fatal("Verify() returned true for nil/empty input")
			}
		})
	}
}

func TestSignInvalidKey(t *testing.T) {
	msg := []byte("test message")

	tests := []struct {
		name string
		key  []byte
	}{
		{name: "nil private key", key: nil},
		{name: "empty private key", key: []byte{}},
		{name: "too short (8 bytes)", key: []byte("deadbeef")},
		{name: "too long (64 bytes)", key: make([]byte, 64)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sig, err := Sign(tt.key, msg)
			if err == nil {
				t.Fatal("Sign() expected error for invalid key, got nil")
			}
			if sig != nil {
				t.Fatal("Sign() expected nil signature on error")
			}
		})
	}
}

func TestSignInvalidMessage(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() failed: %v", err)
	}

	sig, err := Sign(kp.PrivateKey, nil)
	if err == nil {
		t.Fatal("Sign() expected error for nil message, got nil")
	}
	if sig != nil {
		t.Fatal("Sign() expected nil signature on error")
	}
}

func TestSignVerifyRoundTrip(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() failed: %v", err)
	}

	messages := [][]byte{
		[]byte("short"),
		[]byte("a longer message with more bytes to sign"),
		[]byte{0x00, 0x01, 0x02, 0xFF},
		[]byte{},
		make([]byte, 1024), // 1 KB message
	}

	for _, msg := range messages {
		sig, err := Sign(kp.PrivateKey, msg)
		if err != nil {
			t.Fatalf("Sign(%q) failed: %v", string(msg), err)
		}

		valid := Verify(kp.PublicKey, msg, sig)
		if !valid {
			t.Fatalf("Verify() returned false for valid signature on message %q", string(msg))
		}
	}
}

func TestVerifyConcurrent(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() failed: %v", err)
	}

	msg := []byte("concurrent test message")
	sig, err := Sign(kp.PrivateKey, msg)
	if err != nil {
		t.Fatalf("Sign() failed: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			valid := Verify(kp.PublicKey, msg, sig)
			if !valid {
				t.Error("Verify() returned false in concurrent call")
			}
		}()
	}
	wg.Wait()
}
