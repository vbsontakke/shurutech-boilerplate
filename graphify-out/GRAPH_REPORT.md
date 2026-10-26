# Graph Report - shurutech-boilerplate  (2026-10-26)

## Corpus Check
- Corpus is ~8,458 words - fits in a single context window. You may not need a graph.

## Summary
- 344 nodes · 679 edges · 15 communities (13 shown, 2 thin omitted)
- Extraction: 97% EXTRACTED · 3% INFERRED · 0% AMBIGUOUS · INFERRED: 17 edges (avg confidence: 0.89)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- Server Bootstrap
- Project Docs
- HTTP Request Handling
- Database Layer & Tests
- Frontend Dependencies
- Docker Compose Services
- Frontend API Client
- Config Loading
- App Error Types
- TypeScript App Config
- TypeScript Node Config
- Request Validation
- Frontend Lint Config
- TypeScript Project Refs
- Go Module Root

## God Nodes (most connected - your core abstractions)
1. `compilerOptions` - 18 edges
2. `compilerOptions` - 15 edges
3. `Design Template` - 12 edges
4. `Do()` - 10 edges
5. `Validation()` - 9 edges
6. `NewTestDB()` - 9 edges
7. `Backend Requirements` - 9 edges
8. `New()` - 8 edges
9. `NewRouter()` - 8 edges
10. `newTestRouter()` - 8 edges

## Surprising Connections (you probably didn't know these)
- `Project Conventions` --conceptually_related_to--> `Entities and Relationships`  [AMBIGUOUS]
  CLAUDE.md → DESIGN.md
- `Status Mapping` --semantically_similar_to--> `Status Code Mapping`  [INFERRED] [semantically similar]
  README.md → CLAUDE.md
- `Boilerplate Page Title` --conceptually_related_to--> `React TypeScript Vite Frontend`  [INFERRED]
  frontend/index.html → README.md
- `DecodeAndValidate()` --calls--> `Validate()`  [INFERRED]
  backend/internal/httpx/decode.go → backend/internal/httpx/validate.go
- `Adding a Feature` --conceptually_related_to--> `Adding a Feature`  [INFERRED]
  README.md → CLAUDE.md

## Import Cycles
- None detected.

## Communities (15 total, 2 thin omitted)

### Community 0 - "Server Bootstrap"
Cohesion: 0.10
Nodes (9): main(), newLogger(), run(), hasCode(), IsCheckViolation(), IsForeignKeyViolation(), IsUniqueViolation(), createReq (+1 more)

### Community 1 - "Project Docs"
Cohesion: 0.07
Nodes (46): Adding a Feature, apperror Package, database Package, When You Finish a Step, httpx Package, Layout, Money as BIGINT Minor Units, Project Conventions (+38 more)

### Community 2 - "HTTP Request Handling"
Cohesion: 0.08
Nodes (30): Validation(), NewHandler(), DecodeAndValidate(), DecodeJSON(), MethodNotAllowedHandler(), NotFoundHandler(), decode(), TestCursorRoundTrip() (+22 more)

### Community 3 - "Database Layer & Tests"
Cohesion: 0.11
Nodes (27): TestConnectAndMigrate(), TestIsUniqueViolation(), TestWithTx_RollsBackOnError(), Migrate(), WithTx(), TestError_HidesInternalDetails(), TestError_WrappedAppErrorKeepsKind(), TestParsePagination() (+19 more)

### Community 4 - "Frontend Dependencies"
Cohesion: 0.10
Nodes (24): devDependencies, oxlint@^1.88.0, @types/node@^24.19.2, @types/react@^19.3.0, @types/react-dom@^19.3.0, typescript@~6.0.3, vite@^8.3.4, @vitejs/plugin-react@^6.1.2 (+16 more)

### Community 5 - "Docker Compose Services"
Cohesion: 0.10
Nodes (24): POSTGRES_DB, POSTGRES_PASSWORD, POSTGRES_USER, interval, retries, test, timeout, container_name (+16 more)

### Community 6 - "Frontend API Client"
Cohesion: 0.15
Nodes (17): dependencies, react@^19.3.1, react-dom@^19.3.1, api(), ApiError, FieldError, getDbHealth(), getHealth() (+9 more)

### Community 7 - "Config Loading"
Cohesion: 0.12
Nodes (7): getEnv(), Config, DBConfig, Load(), splitCSV(), NewPool(), DBTX

### Community 8 - "App Error Types"
Cohesion: 0.19
Nodes (14): Error, FieldError, Conflict(), Forbidden(), From(), Kind, Internal(), New() (+6 more)

### Community 9 - "TypeScript App Config"
Cohesion: 0.10
Nodes (19): compilerOptions, allowArbitraryExtensions, allowImportingTsExtensions, erasableSyntaxOnly, jsx, lib, module, moduleDetection (+11 more)

### Community 10 - "TypeScript Node Config"
Cohesion: 0.12
Nodes (16): compilerOptions, allowImportingTsExtensions, erasableSyntaxOnly, lib, module, moduleDetection, noEmit, noFallthroughCasesInSwitch (+8 more)

### Community 11 - "Request Validation"
Cohesion: 0.31
Nodes (4): fieldPath(), issueFor(), newValidator(), Validate()

### Community 12 - "Frontend Lint Config"
Cohesion: 0.33
Nodes (5): plugins, rules, react/only-export-components, react/rules-of-hooks, $schema

## Ambiguous Edges - Review These
- `Project Conventions` → `Entities and Relationships`  [AMBIGUOUS]
  CLAUDE.md · relation: conceptually_related_to

## Knowledge Gaps
- **16 isolated node(s):** `github.com/vbsontakke/boilerplate`, `When You Finish a Step`, `Money as BIGINT Minor Units`, `set_updated_at Trigger`, `Assumptions and Clarifications` (+11 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 110 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **2 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Error` connect `App Error Types` to `HTTP Request Handling`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **What connects `github.com/vbsontakke/boilerplate`, `When You Finish a Step`, `Money as BIGINT Minor Units` to the rest of the system?**
  _16 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Server Bootstrap` be split into smaller, more focused modules?**
  _Cohesion score 0.10407239819004525 - nodes in this community are weakly interconnected._
- **Why does `DBTX` connect `Config Loading` to `Server Bootstrap`?**
  _High betweenness centrality (0.008) - this node is a cross-community bridge._
- **Should `Project Docs` be split into smaller, more focused modules?**
  _Cohesion score 0.06763285024154589 - nodes in this community are weakly interconnected._
- **Why does `Handler` connect `HTTP Request Handling` to `Server Bootstrap`?**
  _High betweenness centrality (0.008) - this node is a cross-community bridge._

### Low-confidence Hints
_AMBIGUOUS edges — the extractor was unsure. Verify before acting on these._

- **What is the exact relationship between `Project Conventions` and `Entities and Relationships`?**
  _Edge tagged AMBIGUOUS (relation: conceptually_related_to) - confidence is low._