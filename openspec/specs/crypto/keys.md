# crypto/keys — ECDSA secp256k1 Key Management

## Purpose

Provide cryptographic key generation, signing, and verification using the ECDSA algorithm over the secp256k1 elliptic curve. This is the identity primitive for all actors (validators, nodes, wallets) in the Parva blockchain.

## Package

`github.com/anythinglord/parva-blockchain/internal/domain/crypto`

This is a **domain** package. It MUST NOT import anything from `internal/application` or `internal/infrastructure`. The only external dependency allowed is `github.com/decred/dcrd/dcrec/secp256k1/v4`.

## Type Definitions

```go
// KeyPair represents an ECDSA secp256k1 key pair
type KeyPair struct {
    PrivateKey []byte // raw private key bytes (32 bytes)
    PublicKey  []byte // raw public key bytes, compressed format (33 bytes)
}
```

## Requirements

### Requirement: GenerateKeyPair

The system MUST generate a new ECDSA key pair over the secp256k1 curve.

- Private key MUST be exactly 32 bytes.
- Public key MUST use compressed format (33 bytes) for storage efficiency.
- The generator MUST use `crypto/rand` as the entropy source.
- The function MUST return an error if the underlying entropy source fails.

#### Scenario: Happy path — generates valid key pair

- GIVEN no existing key material
- WHEN `GenerateKeyPair()` is called
- THEN a `*KeyPair` is returned with no error
- AND `PrivateKey` is 32 bytes long
- AND `PublicKey` is 33 bytes long (compressed)
- AND the public key is mathematically derivable from the private key

#### Scenario: Edge case — entropy source failure (unreachable in practice)

- GIVEN an entropy source that cannot provide random bytes
- WHEN `GenerateKeyPair()` is called
- THEN an error MUST be returned
- AND the returned `*KeyPair` MUST be nil

### Requirement: Sign

The system MUST produce an ECDSA signature over a message hash using a secp256k1 private key.

- The function MUST accept a raw private key (32 bytes) and a message hash (any length).
- The function MUST return an error if the private key is invalid (nil, wrong length, or not on curve).
- ECDSA signatures are non-deterministic by nature. Multiple calls with identical inputs MAY produce different signatures, but all MUST be valid under `Verify`.

#### Scenario: Happy path — signs a message hash

- GIVEN a valid 32-byte private key and a 32-byte message hash
- WHEN `Sign(privateKey, messageHash)` is called
- THEN a signature byte slice is returned with no error
- AND the signature is non-empty
- AND `Verify(publicKey, messageHash, signature)` returns true

#### Scenario: Edge case — nil private key

- GIVEN a nil `privateKey`
- WHEN `Sign(nil, messageHash)` is called
- THEN an error MUST be returned
- AND the signature MUST be nil

#### Scenario: Edge case — empty message hash

- GIVEN a valid private key and an empty `messageHash`
- WHEN `Sign(privateKey, []byte{})` is called
- THEN an error MAY be returned (empty hash is rejected)

#### Scenario: Adversarial — wrong private key length

- GIVEN a private key that is not 32 bytes
- WHEN `Sign(invalidKey, messageHash)` is called
- THEN an error MUST be returned
- AND the signature MUST be nil

### Requirement: Verify

The system MUST verify that an ECDSA signature is valid for a given message hash and public key over secp256k1.

- The function MUST NOT panic under any input condition.
- The function MUST return `false` (not panic) for nil or malformed inputs.

#### Scenario: Happy path — valid signature passes

- GIVEN a message hash, a valid signature, and the correct public key
- WHEN `Verify(publicKey, messageHash, signature)` is called
- THEN `true` is returned

#### Scenario: Adversarial — wrong public key

- GIVEN a message hash signed with key A
- WHEN `Verify(publicKeyB, messageHash, signature)` is called (where publicKeyB belongs to a different key pair)
- THEN `false` is returned

#### Scenario: Adversarial — tampered message

- GIVEN a valid signature for message hash H1
- WHEN `Verify(publicKey, H2, signature)` is called (where H2 != H1)
- THEN `false` is returned

#### Scenario: Adversarial — nil public key

- GIVEN a nil public key
- WHEN `Verify(nil, messageHash, signature)` is called
- THEN `false` is returned (no panic)

#### Scenario: Adversarial — nil signature

- GIVEN a nil signature
- WHEN `Verify(publicKey, messageHash, nil)` is called
- THEN `false` is returned (no panic)

#### Scenario: Edge case — empty signature

- GIVEN a valid public key and an empty signature slice
- WHEN `Verify(publicKey, messageHash, []byte{})` is called
- THEN `false` is returned (no panic)

### Requirement: PublicKeyFromPrivateKey

The system MUST derive the compressed secp256k1 public key from a private key deterministically.

- Given the same private key, the function MUST always return the same public key.
- The function MUST return an error for invalid private keys.

#### Scenario: Happy path — deterministic derivation

- GIVEN a valid 32-byte private key
- WHEN `PublicKeyFromPrivateKey(privateKey)` is called
- THEN a 33-byte compressed public key is returned with no error
- AND calling it again with the same private key returns identical bytes

#### Scenario: Edge case — nil private key

- GIVEN a nil private key
- WHEN `PublicKeyFromPrivateKey(nil)` is called
- THEN an error MUST be returned

### Requirement: Thread Safety

All exported functions in this package MUST be safe for concurrent read access.

- `Verify` MAY be called concurrently from multiple goroutines with no synchronization required.
- `GenerateKeyPair` MUST NOT share mutable global state between calls.

### Requirement: Round-Trip Integrity

A full sign/verify round-trip MUST always succeed for any valid key pair.

#### Scenario: Generate → Sign → Verify

- GIVEN a freshly generated key pair
- WHEN `Sign(kp.PrivateKey, hash)` is called, then `Verify(kp.PublicKey, hash, sig)` is called
- THEN `Verify` returns `true`

#### Scenario: Generate → Derive → Sign → Verify

- GIVEN a freshly generated key pair
- WHEN `PublicKeyFromPrivateKey(kp.PrivateKey)` is called and returns derivedPubKey
- AND `Sign(kp.PrivateKey, hash)` is called returning `sig`
- THEN `Verify(derivedPubKey, hash, sig)` returns `true`

## Error Handling

| Condition | Function | Behavior |
|-----------|----------|----------|
| Nil private key | Sign | Return error |
| Nil private key | PublicKeyFromPrivateKey | Return error |
| Invalid private key length | Sign | Return error |
| Invalid private key length | PublicKeyFromPrivateKey | Return error |
| Nil/empty public key | Verify | Return false |
| Nil/empty signature | Verify | Return false |
| Entropy failure | GenerateKeyPair | Return error |

Errors MUST use `fmt.Errorf` with descriptive messages (e.g., "crypto: invalid private key length: got N bytes, want 32").

## Security Considerations

1. **secp256k1 vs stdlib curves**: secp256k1 is NOT in Go stdlib. The `decred/dcrd/dcrec/secp256k1/v4` dependency is a Pure Go implementation with no CGo, reducing cross-compilation risk.
2. **Compressed keys**: 33-byte compressed public keys are chosen over 65-byte uncompressed for storage efficiency. This is the standard for Bitcoin-derived protocols.
3. **Non-deterministic signatures**: Standard ECDSA uses random `k` per signature. This is acceptable for Phase 1. Deterministic signatures (RFC 6979) MAY be added later if replay protection requires it.
4. **No key serialization**: Phase 1 uses raw bytes only. PEM/DER encoding is out of scope until key storage is needed (Phase 8).
5. **No constant-time guarantees**: Phase 1 does not require constant-time operations. A production audit SHOULD verify no timing side-channels exist.

## Test Assertions Checklist

| # | Test | Assertions |
|---|------|------------|
| 1 | `TestGenerateKeyPair` | KeyPair non-nil, PrivateKey len=32, PublicKey len=33, no error |
| 2 | `TestSignAndVerify` | Sign returns sig, Verify returns true for same key |
| 3 | `TestVerifyFailsWrongKey` | Verify returns false for different public key |
| 4 | `TestVerifyFailsTamperedMessage` | Verify returns false for different message hash |
| 5 | `TestVerifyNoPanicNilInput` | Verify(nil, hash, sig) = false, Verify(pub, nil, sig) = false, Verify(pub, hash, nil) = false — no panics |
| 6 | `TestPublicKeyFromPrivateKey` | Returns 33 bytes, matches original KeyPair.PublicKey |
| 7 | `TestPublicKeyFromPrivateKeyNil` | Returns error for nil input |
| 8 | `TestSignInvalidKey` | Returns error for nil or wrong-length key |
| 9 | `TestGenerateSignVerifyRoundTrip` | Generate → Sign → Verify returns true |
| 10 | `TestPublicKeyDerivationDeterministic` | Same private key → same public key on repeated calls |
