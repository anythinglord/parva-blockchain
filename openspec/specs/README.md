# Parva Blockchain — OpenSpec (Source of Truth)

This directory contains the authoritative specification documents for the Parva blockchain.

## Domains

### `crypto/`

Cryptographic primitives for the Parva blockchain:

| File | Description | Status |
|------|-------------|--------|
| `hash.md` | SHA-256 and Double SHA-256 hashing utilities | ✅ Implemented (Fase 1) |
| `keys.md` | ECDSA secp256k1 key generation, signing, and verification | ✅ Implemented (Fase 1) |
| `design.md` | Architecture design and ADRs for Fase 1 | ✅ Finalized |
| `tasks.md` | Implementation tasks for Fase 1 | ✅ All 13 tasks complete |

## Future Domains

| Domain | Planned Phase | Status |
|--------|--------------|--------|
| `transaction/` | Fase 2 | 🔲 Planned |
| `mempool/` | Fase 3 | 🔲 Planned |
| `block/` | Fase 4 | 🔲 Planned |
| `consensus/` | Fase 5 | 🔲 Planned |
| `p2p/` | Fase 6 | 🔲 Planned |
| `node/` | Fase 7 | 🔲 Planned |
| `storage/` | Fase 8 | 🔲 Planned |
| `api/` | Fase 9 | 🔲 Planned |

## Conventions

- Specs use Given/When/Then for scenarios
- Requirements use RFC 2119 keywords (MUST, SHALL, SHOULD, MAY)
- Architecture decisions are documented as ADRs with options and trade-offs
