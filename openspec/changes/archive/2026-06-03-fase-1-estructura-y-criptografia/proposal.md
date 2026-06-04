# Proposal: Fase 1 — Estructura y Criptografía

## Intent

Establecer la base del proyecto: inicializar el módulo Go, crear el esqueleto Clean Architecture completo, e implementar la capa criptográfica con ECDSA (secp256k1) y SHA-256. Sin esta base no podemos construir nada más.

## Scope

### In Scope
- `go mod init github.com/anythinglord/parva-blockchain`
- `.gitignore` con entradas estándar Go
- Árbol Clean Architecture completo con placeholders (`doc.go`)
- `internal/domain/crypto/`: keygen, sign, verify, SHA-256 hash
- Tests TDD estrictos (RED → GREEN → REFACTOR)
- Dependencia externa: `github.com/decred/dcrd/dcrec/secp256k1/v4`

### Out of Scope
- Transacciones (Fase 2), mempool (Fase 3), bloques (Fase 4)
- Consenso PoS (Fase 5), red P2P (Fase 6), persistencia (Fase 8)

## Capabilities

### New Capabilities
- `crypto-keys`: Generación ECDSA secp256k1, firma y verificación
- `crypto-hash`: Hashing SHA-256 utility

### Modified Capabilities
None — proyecto nuevo, sin specs previas.

## Approach

1. **Module init**: `go mod init github.com/anythinglord/parva-blockchain`
2. **Scaffold**: `cmd/` + `internal/{domain,application,infrastructure}/` con `doc.go`
3. **Dependencia**: `decred/dcrd/dcrec/secp256k1/v4` sobre `btcd/btcec` (menos modular) y `go-ethereum/secp256k1` (CGo). Go stdlib no soporta secp256k1.
4. **TDD estricto** (config.yaml: `strict_tdd: true`, runner: `go test`). Domain zero external deps excepto secp256k1.
5. **Verificación**: `go test ./...` + `staticcheck ./...` + `go vet ./...`

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `go.mod` | New | Module + dependency declaration |
| `.gitignore` | Modified | Go std ignores, IDE, OS |
| `cmd/node/main.go` | New | Entry point |
| `internal/domain/crypto/` | New | Keygen, sign, verify, hash |
| `internal/*/*/doc.go` | New | 10 placeholder packages |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| secp256k1 ausente en stdlib | High | decred/dcrec — Pure Go, sin CGo |
| Serialización de claves | Low | Hex raw bytes en Fase 1 |

## Rollback Plan

`git revert <commit>` sobre `develop`. Si hay commits posteriores, revertir solo archivos tocados con `git checkout develop~1 -- <paths>`.

## Dependencies

- `github.com/decred/dcrd/dcrec/secp256k1/v4` (ISC license, Pure Go)

## Success Criteria

- [ ] `go test ./...` — 0 fallos
- [ ] `staticcheck ./...` — 0 warnings
- [ ] `go vet ./...` — 0 errores
- [ ] `go build ./...` — compila
- [ ] 9+ tests para crypto (keygen, sign, verify, hash, edge cases)
- [ ] `.gitignore` cubre `*.exe`, `*.test`, `*.out`, `.env`, `.idea/`, `.vscode/`
- [ ] Clean Architecture: domain no importa infrastructure
- [ ] Cada `.go` tiene su `_test.go`
