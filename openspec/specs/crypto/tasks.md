# Tasks: Fase 1 — Estructura y Criptografía

## Completed (All 13 tasks)

### Grupo 1: Project Scaffold
- [x] **1.1** go mod init + .gitignore
- [x] **1.2** Clean Architecture folder tree + doc.go + cmd/node/main.go
- [x] **1.3** Add secp256k1 dependency

### Grupo 2: Hash Functions
- [x] **2.1** RED — Write hash tests
- [x] **2.2** GREEN — Implement HashSHA256 and HashDoubleSHA256
- [x] **2.3** REFACTOR — Clean up hash code

### Grupo 3: Key Generation
- [x] **3.1** RED — Write keygen tests
- [x] **3.2** GREEN — Implement GenerateKeyPair and PublicKeyFromPrivateKey
- [x] **3.3** REFACTOR — Validate keygen error paths

### Grupo 4: Sign & Verify
- [x] **4.1** RED — Write sign/verify tests
- [x] **4.2** GREEN — Implement Sign and Verify
- [x] **4.3** REFACTOR — Nil/empty safety and concurrency test
- [x] **4.4** FINAL VERIFY — Full suite

### Verification
- All 24 tests pass (9 hash + 6 keygen + 9 sign/verify)
- go vet ./... — clean
- staticcheck ./... — clean
- go build ./... — compiles
- Race detector — clean (Verify concurrent test passes)
