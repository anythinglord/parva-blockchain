# Mini Blockchain — Proposal

> Documento vivo. Este proposal se refina a medida que avanzamos en el proyecto.
> Alineado con las directrices de `AGENTS.md`: TDD estricto, Clean Architecture, GitFlow.

---

## 1. Resumen Ejecutivo

Construir una mini-blockchain en Go desde cero, con un enfoque educativo y modular. El objetivo no es solo tener algo funcionando, sino entender cada capa: criptografía, transacciones, mempool, consenso PoS, y comunicación P2P.

Arrancamos con una base mínima e iteramos. Cada etapa agrega una capacidad concreta y verificable.

---

## 2. Stack Tecnológico

| Capa            | Tecnología                                    |
| --------------- | --------------------------------------------- |
| Lenguaje        | Go 1.22+                                      |
| Criptografía    | ECDSA (`secp256k1`), SHA-256, estándar `crypto` |
| Red P2P         | `libp2p` (conexión entre nodos, descubrimiento) |
| Persistencia    | A definir — arrancamos en memoria, luego disco |
| Testing         | `testing` estándar + `go test -bench`         |
| Perfilado       | `pprof`, `go test -bench`, escape analysis     |

> **Decisión abierta:** para storage podemos evaluar BadgerDB, LevelDB (vía `go-ethereum`), o incluso un formato plano. Se define en la etapa de persistencia.

---

## 3. Arquitectura General

Vamos con **Clean Architecture / Hexagonal** desde el día uno.

```
cmd/
  node/            → entry point del nodo
internal/
  domain/          → entidades del negocio (Transaction, Block, Account, Validator)
    crypto/        → firmas, hashes, claves
    transaction/   → modelo y validación de transacciones
    block/         → estructura de bloque y header
    consensus/     → lógica de PoS (validadores, slashing, finalidad)
  application/     → casos de uso (orquestación)
    mempool/       → gestión de transacciones pendientes
    blockchain/    → cadena de bloques, fork choice
    node/          → ciclo de vida del nodo
  infrastructure/  → implementaciones concretas
    p2p/           → transporte libp2p, peers, mensajes
    storage/       → persistencia (archivos, DB)
    api/           → HTTP/gRPC para debugging o RPC liviano
```

### Principios

- **Domain** no depende de nada externo. Zero imports de infraestructura.
- **Application** orquesta, no implementa. Usa interfaces definidas en domain.
- **Infrastructure** implementa esas interfaces. Es reemplazable.
- Todo archivo `.go` de producción tiene su contraparte `_test.go`.

---

## 4. Componentes Principales

### 4.1 Capa Criptográfica

- Generación y manejo de claves ECDSA (`secp256k1`).
- Firma y verificación de transacciones.
- Hashing SHA-256 para bloques y estructuras.
- (Futuro) Merkle Tree para resumen de transacciones en bloque.

### 4.2 Transacciones

- Modelo: `Sender`, `Receiver`, `Value`, `Nonce`, `Signature`, `Timestamp`.
- Validación: firma, nonce, saldo suficiente.
- Serialización/deserialización para red y almacenamiento.

### 4.3 Mempool

- Cola de transacciones pendientes.
- Priorización por fee (cuando tengamos comisiones).
- Expulsión por capacidad o expiración.
- Thread-safe para concurrencia.

### 4.4 Bloques y Blockchain

- Estructura: `Header` (Version, PrevHash, MerkleRoot, Timestamp, Height, Validator) + `Transactions`.
- Validación de bloques: hash, firma del validador, transacciones.
- Fork choice rule: simplest (cadena más larga) al inicio, luego veremos.

### 4.5 Consenso PoS (Simple)

- Conjunto de validadores.
- Selección de proponente por turno o peso.
- Validadores firman bloques propuestos.
- (Futuro) Slashing por doble firma o ausencia.
- (Futuro) Finalidad tipo Casper FFG o similar.

### 4.6 Red P2P

- Identidad de nodo con par de claves.
- Conexión entre peers vía `libp2p`.
- Protocolo de mensajes: `TxBroadcast`, `BlockBroadcast`, `PeerDiscovery`.
- (Futuro) Sync de estado entre nodos.

### 4.7 Almacenamiento

- **Fase 1:** Todo en memoria (maps, slices).
- **Fase 2:** Persistencia a disco (archivo JSON, gob, o similar).
- **Fase 3:** Base de datos embebida (BadgerDB / LevelDB).

---

## 5. Plan de Implementación por Fases

Cada fase produce código funcional y testeado. No se avanza a la siguiente sin pasar toda la suite de tests.

| Fase | Nombre                        | Entrega                                                                 |
| ---- | ----------------------------- | ----------------------------------------------------------------------- |
| 1    | Estructura y Criptografía     | Carpetas Clean Architecture + generación de claves + firma/verificación |
| 2    | Modelo de Transacciones       | TX model, validación, serialización, tests                             |
| 3    | Mempool                       | Mempool thread-safe con inserción/eliminación                           |
| 4    | Blockchain y Almacenamiento   | Block struct, Blockchain chain, storage en memoria                      |
| 5    | Consenso PoS                  | Validators set, proposer selection, block finalization                  |
| 6    | Red P2P                       | Conexión libp2p, broadcast de TX y bloques                              |
| 7    | Integración                   | Nodo completo: minar, propagar, sync básico                             |
| 8    | Persistencia                  | Disco / DB                                                             |
| 9+   | Mejoras                       | Ver sección 7                                                           |

---

## 6. Convenciones y Workflow

### 6.1 TDD Estricto (de AGENTS.md)

Ciclo obligatorio por cada cambio:

1. **RED** — Escribir test que falle (compilación o assertion).
2. **GREEN** — Mínimo código necesario para que pase.
3. **REFACTOR** — Limpiar, optimizar, mantener tests verdes.

> No se escribe ni una línea de producción sin su test correspondiente.

### 6.2 GitFlow (de AGENTS.md)

- `main` — solo releases estables.
- `develop` — integración continua.
- `feature/*` — ramas desde `develop`.
- `bugfix/*` — ramas desde `develop`.
- No se commitea directo a `main` ni `develop`.

### 6.3 Commits

- Convencionales: `feat:`, `fix:`, `refactor:`, `test:`, `docs:`.
- Cada commit deja los tests pasando.
- Un commit por unidad lógica de trabajo.

### 6.4 Seguridad

- **NUNCA** tocar `.env` ni credenciales en el repo (AGENTS.md 2.1).
- Usar `.env.example` con placeholders.
- `.env` en `.gitignore` desde el día uno.

---

## 7. Mejoras Futuras (para definir más adelante)

| Mejora                      | Detalle breve                                                  |
| --------------------------- | -------------------------------------------------------------- |
| Merkle Patricia Trie        | Para estado de cuentas eficiente                               |
| Slashing                    | Penalizar validadores maliciosos o inactivos                   |
| Contratos Inteligentes      | VM mínima tipo stack o WASM                                    |
| API RPC                     | JSON-RPC para interactuar con el nodo                          |
| Explorador de Bloques       | Web UI básica                                                   |
| Wallet CLI                  | Comando para crear TX, ver saldo, etc.                         |
| Sharding / Sidechains       | Escalabilidad horizontal                                       |
| Governance on-chain         | Votación de parámetros por validadores                         |
| Benchmarks y Perfilado      | `go test -bench`, `pprof`, optimización de allocaciones         |

---

## 8. Hitos Clave

- **[Fase 2]** Poder crear y firmar una transacción desde código.
- **[Fase 4]** Tener una blockchain funcional en memoria con validación.
- **[Fase 5]** Minería PoS funcionando con al menos 3 validadores.
- **[Fase 6]** Dos nodos conectados intercambiando transacciones.
- **[Fase 7]** Nodo autónomo que mina, propaga y persiste.

---

## 9. Glosario

| Término       | Significado                                                  |
| ------------- | ------------------------------------------------------------ |
| **PoS**       | Proof of Stake — consenso donde validadores "apuestan" tokens |
| **Mempool**   | Conjunto de transacciones pendientes de ser incluidas en un bloque |
| **libp2p**    | Librería modular para redes P2P                              |
| **ECDSA**     | Algoritmo de firma digital basado en curvas elípticas        |
| **secp256k1** | Curva elíptica usada en Bitcoin y Ethereum                    |
| **Fork**      | División de la cadena por desacuerdo en el estado            |
| **Validator** | Nodo que participa en la producción de bloques vía PoS       |
| **Finality**  | Garantía de que un bloque no será revertido                  |

---

> **Próximo paso:** Arrancar con Fase 1 — estructura de carpetas Clean Architecture + capa criptográfica con TDD.
