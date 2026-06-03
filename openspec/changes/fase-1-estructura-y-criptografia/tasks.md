# Tasks: Fase 1 — Estructura y Criptografía

## Workload Forecast

| Campo | Valor |
|-------|-------|
| Líneas estimadas cambiadas | ~470-520 |
| Riesgo presupuesto 400 líneas | Medio |
| Chained PRs recomendado | No |
| Split sugerido | PR único |
| Delivery strategy | ask-on-risk |
| Chain strategy | pending |

```
Decision needed before apply: Yes
Chained PRs recommended: No
Chain strategy: pending
400-line budget risk: Medium
```

### Work Units Sugeridas

| Unidad | Meta | PR | Notas |
|--------|------|-----|-------|
| 1 | Scaffold + Hash + Keys + Sign/Verify | PR único | Todo en `internal/domain/crypto/`, un solo concern |

---

## Grupo 1: Project Scaffold

```yaml
id: "1.1"
name: go mod init + .gitignore
group: Scaffold
description: "Init module github.com/anythinglord/parva-blockchain, update .gitignore with Go/IDE/OS entries"
tdd: {phase: null}
verification:
  commands: ["cat go.mod", "cat .gitignore"]
  expected: "go.mod válido, .gitignore cubre *.exe *.test *.out .env .idea/ .vscode/ .DS_Store"
acceptance:
  - "go.mod exists with module path matching git remote"
  - ".gitignore excludes binaries, IDE, OS files"
dependencies: []
```

```yaml
id: "1.2"
name: Clean Architecture folder tree + doc.go + cmd/node/main.go
group: Scaffold
description: "Create 9 package dirs, 9 doc.go placeholders, move root main.go → cmd/node/main.go"
tdd: {phase: null}
verification:
  commands: ["go build ./...", "go vet ./..."]
  expected: "Build and vet pass, no compile errors"
acceptance:
  - "go build ./... exits with 0"
  - "cmd/node/main.go exists, root main.go deleted"
  - "All 9 internal packages have doc.go"
dependencies: ["1.1"]
```

```yaml
id: "1.3"
name: Add secp256k1 dependency
group: Scaffold
description: "go get github.com/decred/dcrd/dcrec/secp256k1/v4, go mod tidy, verify build"
tdd: {phase: null}
verification:
  commands: ["go build ./...", "go vet ./..."]
  expected: "secp256k1 resolves and compiles, go.sum present"
acceptance:
  - "go.mod contains decred/dcrec/secp256k1/v4"
  - "go build ./... succeeds"
dependencies: ["1.2"]
```

---

## Grupo 2: Hash Functions

```yaml
id: "2.1"
name: "RED — Write hash tests"
group: Hash Functions
description: "Create crypto_test.go with TestHashSHA256, TestHashDoubleSHA256, TestHashEmpty, TestHashConsistency, TestHashNilInput"
tdd:
  phase: RED
  test_files: ["internal/domain/crypto/crypto_test.go"]
  source_files: []
  require_test_first: true
verification:
  commands: ["go test -count=1 ./internal/domain/crypto/... 2>&1 || true"]
  expected: "Tests fail to compile (no crypto.go yet) or fail on missing symbols"
acceptance:
  - "Test file compiles (even if tests fail due to missing impl)"
  - "All spec scenarios from crypto/hash.md covered"
dependencies: ["1.3"]
```

```yaml
id: "2.2"
name: "GREEN — Implement HashSHA256 and HashDoubleSHA256"
group: Hash Functions
description: "Add HashSHA256(data []byte) []byte and HashDoubleSHA256(data []byte) []byte to crypto.go"
tdd:
  phase: GREEN
  test_files: ["internal/domain/crypto/crypto_test.go"]
  source_files: ["internal/domain/crypto/crypto.go"]
  require_test_first: false
verification:
  commands: ["go test -count=1 ./internal/domain/crypto/...", "go vet ./internal/domain/crypto/..."]
  expected: "All hash tests pass, no vet warnings"
acceptance:
  - "go test passes all hash tests"
  - "HashSHA256(nil) returns SHA-256(\"\"), no panic"
dependencies: ["2.1"]
```

```yaml
id: "2.3"
name: "REFACTOR — Clean up hash code"
group: Hash Functions
description: "Check edge cases, rename anything unclear, ensure composition property holds"
tdd:
  phase: REFACTOR
  test_files: ["internal/domain/crypto/crypto_test.go"]
  source_files: ["internal/domain/crypto/crypto.go"]
  require_test_first: false
verification:
  commands: ["go test -count=1 ./internal/domain/crypto/...", "go vet ./internal/domain/crypto/..."]
  expected: "Tests still green, vet clean"
acceptance:
  - "No unused or redundant code"
  - "HashDoubleSHA256(X) === HashSHA256(HashSHA256(X))"
dependencies: ["2.2"]
```

---

## Grupo 3: Key Generation

```yaml
id: "3.1"
name: "RED — Write keygen tests"
group: Key Generation
description: "Add tests to keys_test.go: TestGenerateKeyPair, TestPublicKeyFromPrivateKey, TestPublicKeyFromPrivateKeyNil, TestKeyPairLengths"
tdd:
  phase: RED
  test_files: ["internal/domain/crypto/keys_test.go"]
  source_files: []
  require_test_first: true
verification:
  commands: ["go test -count=1 ./internal/domain/crypto/... 2>&1 || true"]
  expected: "Tests fail to compile (keys.go missing)"
acceptance:
  - "Tests cover GenerateKeyPair, PublicKeyFromPrivateKey, nil/invalid inputs"
  - "All spec scenarios from crypto/keys.md covered"
dependencies: ["1.3", "2.3"]
```

```yaml
id: "3.2"
name: "GREEN — Implement GenerateKeyPair and PublicKeyFromPrivateKey"
group: Key Generation
description: "Implement KeyPair struct, GenerateKeyPair() (*KeyPair, error), PublicKeyFromPrivateKey(priv []byte) ([]byte, error) in keys.go"
tdd:
  phase: GREEN
  test_files: ["internal/domain/crypto/keys_test.go"]
  source_files: ["internal/domain/crypto/keys.go"]
  require_test_first: false
verification:
  commands: ["go test -count=1 ./internal/domain/crypto/...", "go vet ./internal/domain/crypto/..."]
  expected: "Keygen tests pass, vet clean"
acceptance:
  - "GenerateKeyPair returns 32B priv + 33B pub, no error"
  - "PublicKeyFromPrivateKey(nil) returns error"
dependencies: ["3.1"]
```

```yaml
id: "3.3"
name: "REFACTOR — Validate keygen error paths"
group: Key Generation
description: "Ensure error messages use crypto: prefix, entropy failure path exists, no dead code"
tdd:
  phase: REFACTOR
  test_files: ["internal/domain/crypto/keys_test.go"]
  source_files: ["internal/domain/crypto/keys.go"]
  require_test_first: false
verification:
  commands: ["go test -count=1 ./internal/domain/crypto/...", "go vet ./internal/domain/crypto/..."]
  expected: "Tests green, vet clean"
acceptance:
  - "All errors use fmt.Errorf with crypto: prefix"
  - "No unused parameters or variables"
dependencies: ["3.2"]
```

---

## Grupo 4: Sign & Verify

```yaml
id: "4.1"
name: "RED — Write sign/verify tests"
group: Sign & Verify
description: "Add to keys_test.go: TestSignAndVerify, TestVerifyFailsWrongKey, TestVerifyFailsTamperedMessage, TestVerifyNoPanicNilInput, TestSignInvalidKey, TestGenerateSignVerifyRoundTrip, TestPublicKeyDerivationDeterministic"
tdd:
  phase: RED
  test_files: ["internal/domain/crypto/keys_test.go"]
  source_files: []
  require_test_first: true
verification:
  commands: ["go test -count=1 ./internal/domain/crypto/... 2>&1 || true"]
  expected: "Tests fail (Sign/Verify not yet implemented)"
acceptance:
  - "All 7 sign/verify tests written, covering happy path + adversarial"
  - "Table-driven with t.Run subtests"
dependencies: ["3.3"]
```

```yaml
id: "4.2"
name: "GREEN — Implement Sign and Verify"
group: Sign & Verify
description: "Implement Sign(privateKey, msg []byte) ([]byte, error) and Verify(publicKey, msg, sig []byte) bool in keys.go"
tdd:
  phase: GREEN
  test_files: ["internal/domain/crypto/keys_test.go"]
  source_files: ["internal/domain/crypto/keys.go"]
  require_test_first: false
verification:
  commands: ["go test -count=1 ./internal/domain/crypto/...", "go vet ./internal/domain/crypto/..."]
  expected: "All sign/verify tests pass"
acceptance:
  - "Sign returns DER-encoded signature for valid key"
  - "Verify returns false (no panic) for nil/empty inputs"
dependencies: ["4.1"]
```

```yaml
id: "4.3"
name: "REFACTOR — Nil/empty safety and concurrency test"
group: Sign & Verify
description: "Ensure no panics for any nil/empty input, add t.Parallel() sub-test for concurrent Verify calls"
tdd:
  phase: REFACTOR
  test_files: ["internal/domain/crypto/keys_test.go"]
  source_files: ["internal/domain/crypto/keys.go"]
  require_test_first: false
verification:
  commands: ["go test -count=1 -race ./internal/domain/crypto/...", "go vet ./internal/domain/crypto/..."]
  expected: "Race detector clean, all tests green"
acceptance:
  - "Verify can be called concurrently with no data races"
  - "All edge cases from specs handled"
dependencies: ["4.2"]
```

```yaml
id: "4.4"
name: "FINAL VERIFY — Full suite"
group: Sign & Verify
description: "Run go test ./..., go vet ./..., staticcheck ./... on entire project"
tdd: {phase: null}
verification:
  commands: ["go test -count=1 ./...", "go vet ./...", "staticcheck ./..."]
  expected: "All tests pass, vet clean, staticcheck clean (0 warnings)"
acceptance:
  - "go test ./... — 0 failures"
  - "go vet ./... — 0 errors"
  - "staticcheck ./... — 0 warnings"
  - "go build ./... — compiles"
dependencies: ["4.3"]
```

## Resumen

| Grupo | Tasks | Resultado |
|-------|-------|-----------|
| Scaffold | 3 | Módulo + árbol + dependencia lista |
| Hash | 3 | SHA-256 + Double SHA-256 implementado |
| Keys (keygen) | 3 | KeyPair, GenerateKeyPair, PublicKeyFromPrivateKey |
| Sign/Verify | 4 | Sign, Verify + final full-suite verify |
| **Total** | **13** | |

### Orden de Implementación

1. **Scaffold** (1.1→1.2→1.3) — todo lo demás depende de esto
2. **Hash** (2.1→2.2→2.3) — independiente, puede ir en paralelo con Keys si se quisiera
3. **Keygen** (3.1→3.2→3.3) — debe ir antes que Sign/Verify
4. **Sign/Verify** (4.1→4.2→4.3→4.4) — depende de keys para generar pares

### Review Workload Forecast

- **Líneas estimadas**: ~470-520 (tests ~300, prod ~130, scaffold ~60, movimientos/borrados ~7)
- **Riesgo 400 líneas**: Medio — supera por poco, pero todo es un solo package (`crypto`)
- **Chained PRs**: No recomendado — todo el código cae en `internal/domain/crypto/`, separar en PRs encadenados agregaría más fricción que beneficio
- **Delivery strategy**: `ask-on-risk` → se necesita decisión del usuario antes de apply
- **Decisión necesaria**: Sí — el estimado supera 400 líneas, usuario debe aprobar PR único o elegir split

### Siguiente Paso

Si se aprueba como PR único: listo para `sdd-apply`. Si se prefiere dividir, se puede dividir en 2 PRs: (1) Scaffold + Hash, (2) Keygen + Sign/Verify — pero no se recomienda por la cohesión del package.
