## Exploration: Fase 1 — Estructura y Criptografía

### Current State

The project is a fresh scaffold with a single commit (`d58b153`). It has:
- `AGENTS.md` — defines TDD lifecycle (RED-GREEN-REFACTOR), GitFlow, Clean Architecture, .env isolation
- `proposal.md` — full project proposal with Clean Architecture layout and phased plan
- `main.go` — "Hello World" in root package
- `openspec/config.yaml` — SDD config with strict TDD, staticcheck, go vet, gofmt
- `openspec/changes/` and `openspec/specs/` — both empty
- `.gitignore` — only `.env` (missing Go-specific entries)
- `.atl/` — skill registry and cache
- No `go.mod` — module not initialized
- No `_test.go` files
- No folder structure beyond root
- Git on `develop` branch (correct for GitFlow), no uncommitted production code

### Affected Areas

- **Root** (`/`) — needs `go.mod` initialization, `.gitignore` update
- **`cmd/node/main.go`** — new entry point (move from root `main.go`)
- **`internal/domain/crypto/`** — new package for Phase 1
- **`internal/domain/transaction/`** — empty placeholder
- **`internal/domain/block/`** — empty placeholder
- **`internal/domain/consensus/`** — empty placeholder
- **`internal/application/mempool/`** — empty placeholder
- **`internal/application/blockchain/`** — empty placeholder
- **`internal/application/node/`** — empty placeholder
- **`internal/infrastructure/p2p/`** — empty placeholder
- **`internal/infrastructure/storage/`** — empty placeholder
- **`internal/infrastructure/api/`** — empty placeholder

### Environment & Tools

| Item | Value |
|------|-------|
| Go version | `go1.25.5 X:nodwarf5 linux/amd64` |
| Module | Not initialized — need `go mod init github.com/gentleman-programming/mini-blockchain` |
| Linter | `staticcheck` available at `/home/jesus/go/bin/staticcheck` |
| Type checker | `go vet` built-in |
| Formatter | `gofmt` (standard) |
| golangci-lint | **Not installed** |
| gofumpt | **Not installed** |

### Key Discovery: secp256k1 Is NOT in Go Standard Library

**Critical finding.** Go's `crypto/elliptic` only supports NIST P-224, P-256, P-384, P-521. Go's `crypto/ecdh` adds X25519 but still no secp256k1.

**Recommended external dependency:** `github.com/decred/dcrd/dcrec/secp256k1/v4`

| Library | Type | Pros | Cons |
|---------|------|------|------|
| **decred/dcrec/secp256k1/v4** ✅ | Pure Go | No CGo, implements `crypto/elliptic.Curve`, has dedicated `ecdsa` sub-package, production-proven (Decred), ISC license, Go 1.17+ | Additional dependency |
| btcd/btcec/v2 | Pure Go | Also pure Go, Bitcoin-proven | Less modular, coupled to btcd types |
| go-ethereum/secp256k1 | CGo wrapper | Fastest (C libsecp256k1) | CGo required (cross-compile pain), heavier dep |

**Recommendation:** Use `decred/dcrd/dcrec/secp256k1/v4` — its `ecdsa` sub-package is optimized specifically for secp256k1 and is the idiomatic choice.

### Dependencies Required (Phase 1)

| Package | Purpose | Type |
|---------|---------|------|
| `crypto/ecdsa` | ECDSA interface | Standard library |
| `crypto/sha256` | SHA-256 hashing | Standard library |
| `crypto/rand` | Cryptographic random reader | Standard library |
| `github.com/decred/dcrd/dcrec/secp256k1/v4` | secp256k1 curve + optimized ECDSA | External |
| `github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa` | secp256k1-specific signing/verification | External (sub-package) |
| `github.com/decred/dcrd/dcrec/secp256k1/v4/schnorr` | (Future) Schnorr signatures | External (sub-package) |

### Folder Structure to Create

```
cmd/
  node/
    main.go                   ← Entry point (move from root)
internal/
  domain/
    crypto/
      keygen.go               ← Key generation
      keygen_test.go          ← TDD: key generation tests
      sign.go                 ← Signing + verification
      sign_test.go            ← TDD: sign/verify tests
      hash.go                 ← SHA-256 convenience wrappers
      hash_test.go            ← TDD: hash tests
    transaction/              ← EMPTY (Phase 2)
    block/                    ← EMPTY (Phase 4)
    consensus/                ← EMPTY (Phase 5)
  application/
    mempool/                  ← EMPTY (Phase 3)
    blockchain/               ← EMPTY (Phase 4)
    node/                     ← EMPTY (Phase 7)
  infrastructure/
    p2p/                      ← EMPTY (Phase 6)
    storage/                  ← EMPTY (Phase 8)
    api/                      ← EMPTY (Phase 9+)
```

### .gitignore Update Needed

Current: only `.env`. Missing Go entries:
```
# Binaries
*.exe
*.test
*.out

# Go workspace
go.work
go.work.sum

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Environment
.env
```

### Test Strategy for Phase 1

Following strict TDD (RED-GREEN-REFACTOR) per AGENTS.md:

| Test | What It Covers |
|------|---------------|
| `TestGenerateKeyPair` | Generates key pair, verifies private/public key are non-nil, curve is secp256k1 |
| `TestSignAndVerify` | Signs a known message, verifies with correct public key |
| `TestVerifyFailsWrongKey` | Verification with wrong public key returns false |
| `TestVerifyFailsTamperedMessage` | Verification with tampered message returns false |
| `TestPublicKeyFromPrivateKey` | Derives public key from private key deterministically |
| `TestHashSHA256` | SHA-256 produces correct hex output for known input |
| `TestHashConsistency` | Same input always produces same hash |
| `TestHashEmpty` | Empty input hash is valid |
| `TestKeyPairRoundTrip` | Serialize/deserialize keys (bytes) round-trips correctly |

All tests use `go test -v -count=1` (no caching) during development.

### Risks and Edge Cases

| Risk | Mitigation |
|------|-----------|
| **secp256k1 not in stdlib** | ✅ Addressed — using decred/dcrec/secp256k1/v4 |
| **Go 1.25 compatibility** | Low risk — decred module requires Go 1.17+, well below 1.25 |
| **Private key serialization** | Must decide encoding format (PEM, DER, raw bytes hex). Recommendation: hex-encoded raw bytes for simplicity in Phase 1 |
| **Non-deterministic ECDSA** | `crypto/ecdsa` signatures are non-deterministic — same message + same key produces different signatures each time. This is OK for Phase 1; verification must accept any valid signature |
| **No go.sum yet** | Will be created on first `go mod tidy` |
| **Static analysis on empty packages** | Package placeholders that only contain a doc comment should compile; empty Go packages with `_test.go` but no `.go` will fail. Solution: each package gets at least a doc+export file even if empty |

### Ready for Proposal

Yes. The exploration is complete. Clear path forward:
1. Initialize Go module
2. Update `.gitignore`
3. Create folder structure
4. Select and add `decred/dcrd/dcrec/secp256k1/v4` dependency
5. Implement crypto layer TDD: RED → GREEN → REFACTOR per test
6. Verify with `go test ./...` and `staticcheck ./...`
