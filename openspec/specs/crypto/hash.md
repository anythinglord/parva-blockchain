# crypto/hash — SHA-256 Hashing Utilities

## Purpose

Provide SHA-256 and Double SHA-256 hashing convenience functions. These are the fundamental building blocks for block hashing, transaction IDs, and Merkle trees in the Parva blockchain.

## Package

`github.com/anythinglord/parva-blockchain/internal/domain/crypto`

This is a **domain** package. It MUST NOT import anything from `internal/application` or `internal/infrastructure`. Only stdlib `crypto/sha256` is allowed.

## Requirements

### Requirement: HashSHA256

The system MUST compute the SHA-256 digest of arbitrary input data.

- The function MUST always return exactly 32 bytes.
- The function MUST be deterministic: identical input always produces identical output.
- The function MUST handle nil and empty input gracefully, returning the valid SHA-256 of the empty string.
- The function MUST NOT panic for any input value.

#### Scenario: Happy path — hashes known input

- GIVEN the byte sequence `[]byte("hello")`
- WHEN `HashSHA256(input)` is called
- THEN a 32-byte slice is returned
- AND the result matches the well-known SHA-256 hash of "hello" (`2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824`)

#### Scenario: Determinism — same input → same output

- GIVEN a fixed input `[]byte("parva-blockchain")`
- WHEN `HashSHA256(input)` is called twice
- THEN both calls return byte-identical results

#### Scenario: Edge case — nil input

- GIVEN a nil input
- WHEN `HashSHA256(nil)` is called
- THEN a 32-byte slice is returned (no panic)
- AND the result equals `SHA-256("")` = `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`

#### Scenario: Edge case — empty slice

- GIVEN an empty byte slice `[]byte{}`
- WHEN `HashSHA256(input)` is called
- THEN a 32-byte slice is returned
- AND the result is identical to the nil-input case

#### Scenario: Adversarial — large input

- GIVEN a 1 MB input of random bytes
- WHEN `HashSHA256(input)` is called
- THEN a 32-byte slice is returned with no error or panic
- AND the operation completes in reasonable time

### Requirement: HashDoubleSHA256

The system MUST compute SHA-256(SHA-256(data)) — double hashing, as used in Bitcoin-derived protocols for transaction and block identifiers.

- The function MUST always return exactly 32 bytes.
- The function MUST be deterministic.
- The function MUST handle nil and empty input gracefully.

#### Scenario: Happy path — double hashes correctly

- GIVEN the input `[]byte("hello")`
- WHEN `HashDoubleSHA256(input)` is called
- THEN a 32-byte slice is returned
- AND the result equals `SHA-256(SHA-256("hello"))` — verifiable by hashing the single-hash output again

#### Scenario: Determinism — same input → same output

- GIVEN a fixed input `[]byte("parva-blockchain")`
- WHEN `HashDoubleSHA256(input)` is called twice
- THEN both calls return byte-identical results

#### Scenario: Composition — DoubleSHA256 is SHA256∘SHA256

- GIVEN any input X
- WHEN `HashDoubleSHA256(X)` is called
- THEN the result MUST equal `HashSHA256(HashSHA256(X))`

#### Scenario: Edge case — nil input

- GIVEN a nil input
- WHEN `HashDoubleSHA256(nil)` is called
- THEN a 32-byte slice is returned (no panic)

#### Scenario: Edge case — empty slice

- GIVEN an empty byte slice
- WHEN `HashDoubleSHA256([]byte{})` is called
- THEN a 32-byte slice is returned
- AND the result equals `HashSHA256(HashSHA256([]byte{}))`

### Requirement: No External Dependencies

Both functions MUST use only `crypto/sha256` from the Go standard library.

- No external dependency is needed or allowed for hashing.
- This ensures the hash functions remain zero-cost dependencies for downstream packages.

## Error Handling

Neither function returns an error. SHA-256 from `crypto/sha256` does not fail for any input (nil, empty, any length). Both functions always succeed.

| Condition | Behavior |
|-----------|----------|
| Nil input | Returns valid SHA-256 of empty string |
| Empty slice | Returns valid SHA-256 of empty string |
| Any other input | Returns valid SHA-256 digest (32 bytes) |

## Security Considerations

1. **SHA-256 is collision-resistant** for the foreseeable future. No blockchain-specific weaknesses apply.
2. **Double SHA-256** mitigates length-extension attacks (though SHA-256 is partially vulnerable). This is the Bitcoin convention for transaction and block IDs.
3. **No constant-time requirement**: Hashing is data-independent by nature in `crypto/sha256` — the stdlib implementation is already constant-time.
4. **No allocation guarantees**: These functions allocate a new `[]byte` on each call. Callers that need zero-alloc should use `crypto/sha256.Sum256` directly and reuse buffers.

## Test Assertions Checklist

| # | Test | Assertions |
|---|------|------------|
| 1 | `TestHashSHA256` | Known input produces expected hex output, result is 32 bytes |
| 2 | `TestHashDoubleSHA256` | Known input double hash is correct, result is 32 bytes |
| 3 | `TestHashConsistency` | Same input called 3x returns identical bytes |
| 4 | `TestHashDoubleSHA256Composition` | `DoubleSHA256(X) == SHA256(SHA256(X))` for known and random inputs |
| 5 | `TestHashNilInput` | HashSHA256(nil) returns 32 bytes = SHA-256(""), no panic |
| 6 | `TestHashEmptyInput` | HashSHA256([]byte{}) returns 32 bytes = SHA-256(""), same as nil |
| 7 | `TestDoubleHashNilInput` | HashDoubleSHA256(nil) returns 32 bytes, no panic |
| 8 | `TestDoubleHashEmptyInput` | HashDoubleSHA256([]byte{}) returns 32 bytes, matches composition property |
