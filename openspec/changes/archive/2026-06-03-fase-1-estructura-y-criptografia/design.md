# Design: Fase 1 — Estructura y Criptografía

## Executive Summary

Initialize the `parva-blockchain` Go module, scaffold the full Clean Architecture tree with placeholder packages, and implement the domain crypto layer (ECDSA secp256k1 key generation, signing, verification, and SHA-256 hashing) using strict TDD. The crypto package is the foundation — no other phase can proceed without it.

## Package Architecture

```
cmd/node/main.go  ──→  entry point (placeholder, moves from root)
                            │
                    ┌───────┴────────┐
                    │  internal/      │
                    │  domain/        │  ← zero infra deps
                    │    crypto/      │  ← Phase 1 implementation
                    │    transaction/  │  ← doc.go placeholder
                    │    block/        │  ← doc.go placeholder
                    │    consensus/    │  ← doc.go placeholder
                    ├────────────────┤
                    │  application/   │  ← all doc.go placeholders
                    │    mempool/     │
                    │    blockchain/  │
                    │    node/        │
                    ├────────────────┤
                    │  infrastructure/│  ← all doc.go placeholders
                    │    p2p/         │
                    │    storage/     │
                    │    api/         │
                    └────────────────┘

Dependency graph for Phase 1:

  internal/domain/crypto
      ├── crypto/sha256 (stdlib)
      ├── crypto/rand  (stdlib)
      └── github.com/decred/dcrd/dcrec/secp256k1/v4      (external)
          └── github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa  (sub-pkg)
```

## Key Design Decisions (ADRs)

| # | Decision | Options | Trade-off | Chosen |
|---|----------|---------|-----------|--------|
| 1 | secp256k1 lib | stdlib (no), btcd/btcec (coupled), go-ethereum (CGo), **decred/dcrec** | decred: pure Go, crypto/elliptic.Curve compat, ISC license, optimized ecdsa sub-pkg | **decred/dcrec/secp256k1/v4** |
| 2 | Public key format | Uncompressed (65B), **Compressed (33B)** | 33B saves ~50% in blocks/network; slight CPU cost for decompression | **Compressed** (Bitcoin standard) |
| 3 | Hash return type | `[32]byte` vs **`[]byte`** | Array gives compile-time length guarantee; slice is idiomatic Go, easier serialization | **`[]byte`** (stdlib pattern) |
| 4 | Deterministic sigs | RFC 6979 vs **non-deterministic ECDSA** | Deterministic: reproducible, testable; Non-det: standard ECDSA, acceptable for Phase 1 | **Non-deterministic** (defer RFC 6979) |
| 5 | Constant-time | **No in Phase 1** vs full CT | CT prevents timing attacks; Phase 1 is educational, not production | **Not in Phase 1** (noted tech debt) |
| 6 | Sign/Verify sig type | `*ecdsa.Signature` vs **raw `[]byte`** | Typed struct is safer; raw bytes are serialization-ready, simpler for Phase 1 | **Raw `[]byte`** (DER-encoded) |

## Implementation Plan (TDD Batches)

### Batch 1: Scaffold + Module
- **RED**: Create `go.mod`, folder tree with `doc.go`, update `.gitignore`
- **GREEN**: `go mod init`, `mkdir -p`, write files, `go mod tidy`
- **VERIFY**: `go build ./...`, `go vet ./...`

### Batch 2: Hash Functions
- **RED**: `TestHashSHA256`, `TestHashDoubleSHA256`, `TestHashEmpty`, `TestHashConsistency`
- **GREEN**: Implement `HashSHA256(data []byte) []byte`, `HashDoubleSHA256(data []byte) []byte`
- **REFACTOR**: Check naming, edge cases
- **VERIFY**: `go test ./internal/domain/crypto/...`

### Batch 3: Key Generation
- **RED**: `TestGenerateKeyPair`, `TestPublicKeyFromPrivateKey`, `TestPublicKeyFromPrivateKeyNil`
- **GREEN**: Implement `GenerateKeyPair() (*KeyPair, error)`, `PublicKeyFromPrivateKey(priv []byte) ([]byte, error)`
- **REFACTOR**: Validate error messages, ensure entropy failure path
- **VERIFY**: `go test ./internal/domain/crypto/...`

### Batch 4: Sign + Verify
- **RED**: `TestSignAndVerify`, `TestVerifyFailsWrongKey`, `TestVerifyFailsTamperedMessage`, `TestVerifyNoPanicNilInput`, `TestSignInvalidKey`, `TestGenerateSignVerifyRoundTrip`, `TestPublicKeyDerivationDeterministic`
- **GREEN**: Implement `Sign(privateKey, msg []byte) ([]byte, error)`, `Verify(publicKey, msg, sig []byte) bool`
- **REFACTOR**: Ensure nil/empty inputs handled, no panics
- **VERIFY**: `go test ./...` + `staticcheck ./...` + `go vet ./...`

## Test Architecture

- **Framework**: Go `testing` only (no testify, no gomega)
- **Pattern**: Table-driven tests with named sub-tests (`t.Run`)
- **Coverage scope**: All exported functions, all spec scenarios, adversarial inputs
- **No caching**: `go test -count=1` during development
- **Concurrent safety test**: `Verify` called from `t.Parallel()` sub-tests

### Edge Cases Covered

| Input | Expected Behavior |
|-------|------------------|
| `nil` input to HashSHA256 | Returns SHA-256("") — valid 32 bytes |
| `[]byte{}` to HashSHA256 | Same as nil — valid 32 bytes |
| `nil` private key to Sign | Error |
| Wrong-length private key to Sign | Error |
| `nil` public key to Verify | `false` (no panic) |
| `nil` signature to Verify | `false` (no panic) |
| Empty signature to Verify | `false` (no panic) |
| Tampered message to Verify | `false` |
| Wrong key to Verify | `false` |
| 1 MB input to HashSHA256 | Valid 32 bytes, no timeout |
| Determinism (3x same input) | Identical bytes every call |

### Tests Required

| Package | Tests | Count |
|---------|-------|-------|
| `domain/crypto` | hash tests | 8 |
| `domain/crypto` | key tests | 10 |
| **Total** | | **18** |

## Sign/Verify Sequence

```
Sign(privateKey [32]byte, msg []byte):
  1. Validate privateKey length == 32
  2. Parse privateKey → secp256k1.PrivateKey via secp256k1.PrivKeyFromBytes
  3. Compute signature: ecdsa.Sign(privateKey, msg)
     → uses random k, returns *Signature
  4. Serialize signature: signature.Serialize() → DER-encoded []byte
  5. Return ([]byte, nil)

Verify(publicKey [33]byte, msg []byte, sig []byte):
  1. If any input is nil/empty → return false
  2. Parse publicKey → secp256k1.PublicKey via secp256k1.ParsePubKey
     (if parse fails → return false)
  3. Parse signature via ecdsa.ParseDERSignature
     (if parse fails → return false)
  4. Verify: ecdsa.Verify(publicKey, msg, signature)
  5. Return result (bool)
```

## File-by-File Breakdown

| File | Action | Content |
|------|--------|---------|
| `go.mod` | Create | Module `github.com/anythinglord/parva-blockchain`, Go 1.25, dep `decred/dcrec/secp256k1/v4` |
| `.gitignore` | Modify | Add Go binaries (`*.exe`, `*.test`, `*.out`), IDE (`.idea/`, `.vscode/`), OS (`.DS_Store`) |
| `cmd/node/main.go` | Create | `package main` with empty `func main(){}` — placeholder |
| `internal/domain/crypto/crypto.go` | Create | `HashSHA256`, `HashDoubleSHA256` — stdlib only |
| `internal/domain/crypto/keys.go` | Create | `KeyPair` struct, `GenerateKeyPair`, `PublicKeyFromPrivateKey`, `Sign`, `Verify` |
| `internal/domain/crypto/crypto_test.go` | Create | 8 hash tests (table-driven) |
| `internal/domain/crypto/keys_test.go` | Create | 10 key tests (table-driven) |
| `internal/*/*/doc.go` (×9) | Create | `// Package ... is a placeholder for Phase N` |
| `main.go` (root) | Delete | Moved to `cmd/node/main.go` |

## Risks and Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| `decred/dcrec` API changes in v4 | Low | Medium | Pin to v4, use well-documented public API surfaces |
| Non-deterministic sigs flaky in CI | Medium | Low | Tests use `Verify()` not byte comparison; no flakiness |
| Module path mismatch | Low | High | Confirm with `git remote -v`: `github.com/anythinglord/parva-blockchain` |
| Placeholder packages fail `staticcheck` | Low | Low | `doc.go` only — `package X` compiles fine, `staticcheck` skips empty packages |

## Open Questions

- [ ] Final module path: `github.com/anythinglord/parva-blockchain` (from git remote) or `github.com/gentleman-programming/mini-blockchain` (from exploration)? — Resolved: use git remote.
- [ ] Should `main.go` in root be deleted or kept as redirect? — Design says delete, move to `cmd/node/`.